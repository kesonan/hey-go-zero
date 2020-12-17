package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-xorm/builder"
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	courseFieldNames          = builderx.FieldNames(&Course{})
	courseRows                = strings.Join(courseFieldNames, ",")
	courseRowsExpectAutoSet   = strings.Join(stringx.Remove(courseFieldNames, "id", "create_time", "update_time"), ",")
	courseRowsWithPlaceHolder = strings.Join(stringx.Remove(courseFieldNames, "id", "create_time", "update_time"), "=?,") + "=?"

	cacheCourseIdPrefix   = "cache#Course#id#"
	cacheCourseNamePrefix = "cache#Course#name#"
)

type (
	CourseModel interface {
		Insert(data Course) (sql.Result, error)
		FindOne(id int64) (*Course, error)
		FindOneByName(name string) (*Course, error)
		Update(data Course) error
		Delete(id int64) error
		FindAllCount() (int, error)
		FindLimit(page, size int) ([]*Course, error)
		FindByIds(ids []int64) ([]*Course, error)
	}

	defaultCourseModel struct {
		sqlc.CachedConn
		table string
	}

	Course struct {
		Id          int64     `db:"id"`
		Name        string    `db:"name"`         // 书籍名称
		Description string    `db:"description"`  // 书籍描述
		Classify    string    `db:"classify"`     // 书籍分类，目前仅支持 【天文|地理|数学|物理|机械|航天|医学|信息|互联网|计算机】
		GenderLimit int64     `db:"gender_limit"` // 性别限制 0-不限，1-男，2-女
		MemberLimit int64     `db:"member_limit"` // 限制人数 0-不限
		Credit      int64     `db:"credit"`       // 学分
		CreateTime  time.Time `db:"create_time"`
		UpdateTime  time.Time `db:"update_time"`
	}
)

func NewCourseModel(conn sqlx.SqlConn, c cache.CacheConf) CourseModel {
	return &defaultCourseModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "course",
	}
}

func (m *defaultCourseModel) Insert(data Course) (sql.Result, error) {
	courseNameKey := fmt.Sprintf("%s%v", cacheCourseNamePrefix, data.Name)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, courseRowsExpectAutoSet)
		return conn.Exec(query, data.Name, data.Description, data.Classify, data.GenderLimit, data.MemberLimit, data.Credit)
	}, courseNameKey)
	return ret, err
}

func (m *defaultCourseModel) FindOne(id int64) (*Course, error) {
	courseIdKey := fmt.Sprintf("%s%v", cacheCourseIdPrefix, id)
	var resp Course
	err := m.QueryRow(&resp, courseIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where id = ? limit 1", courseRows, m.table)
		return conn.QueryRow(v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultCourseModel) FindOneByName(name string) (*Course, error) {
	courseNameKey := fmt.Sprintf("%s%v", cacheCourseNamePrefix, name)
	var resp Course
	err := m.QueryRowIndex(&resp, courseNameKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where name = ? limit 1", courseRows, m.table)
		if err := conn.QueryRow(&resp, query, name); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultCourseModel) Update(data Course) error {
	courseIdKey := fmt.Sprintf("%s%v", cacheCourseIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where id = ?", m.table, courseRowsWithPlaceHolder)
		return conn.Exec(query, data.Name, data.Description, data.Classify, data.GenderLimit, data.MemberLimit, data.Credit, data.Id)
	}, courseIdKey)
	return err
}

func (m *defaultCourseModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	courseIdKey := fmt.Sprintf("%s%v", cacheCourseIdPrefix, id)
	courseNameKey := fmt.Sprintf("%s%v", cacheCourseNamePrefix, data.Name)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where id = ?", m.table)
		return conn.Exec(query, id)
	}, courseIdKey, courseNameKey)
	return err
}

func (m *defaultCourseModel) FindAllCount() (int, error) {
	query := fmt.Sprintf("select count(id) from %s", m.table)
	var count int
	err := m.CachedConn.QueryRowNoCache(&count, query)
	return count, err
}

func (m *defaultCourseModel) FindLimit(page, size int) ([]*Course, error) {
	query := fmt.Sprintf("select %s from %s order by id limit ?,?", courseRows, m.table)
	var resp []*Course
	err := m.CachedConn.QueryRowsNoCache(&resp, query, (page-1)*size, size)
	return resp, err
}

func (m *defaultCourseModel) FindByIds(ids []int64) ([]*Course, error) {
	query, args, err := builder.Select(courseRows).From(m.table).Where(builder.Eq{"id": ids}).ToSQL()
	if err != nil {
		return nil, err
	}

	var resp []*Course
	err = m.CachedConn.QueryRowsNoCache(&resp, query, args...)

	return resp, err
}

func (m *defaultCourseModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheCourseIdPrefix, primary)
}

func (m *defaultCourseModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where id = ? limit 1", courseRows, m.table)
	return conn.QueryRow(v, query, primary)
}
