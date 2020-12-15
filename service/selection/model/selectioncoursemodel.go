package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	selectionCourseFieldNames          = builderx.FieldNames(&SelectionCourse{})
	selectionCourseRows                = strings.Join(selectionCourseFieldNames, ",")
	selectionCourseRowsExpectAutoSet   = strings.Join(stringx.Remove(selectionCourseFieldNames, "id", "create_time", "update_time"), ",")
	selectionCourseRowsWithPlaceHolder = strings.Join(stringx.Remove(selectionCourseFieldNames, "id", "create_time", "update_time"), "=?,") + "=?"

	cacheSelectionCourseIdPrefix = "cache#SelectionCourse#id#"
)

type (
	SelectionCourseModel interface {
		Insert(data SelectionCourse) (sql.Result, error)
		FindOne(id int64) (*SelectionCourse, error)
		Update(data SelectionCourse) error
		Delete(id int64) error
	}

	defaultSelectionCourseModel struct {
		sqlc.CachedConn
		table string
	}

	SelectionCourse struct {
		SelectionId int64     `db:"selection_id"` // 选课任务id
		CourseId    int64     `db:"course_id"`    // 课程id
		TeacherId   int64     `db:"teacher_id"`   // 任教教师
		CreateTime  time.Time `db:"create_time"`
		UpdateTime  time.Time `db:"update_time"`
		Id          int64     `db:"id"`
	}
)

func NewSelectionCourseModel(conn sqlx.SqlConn, c cache.CacheConf) SelectionCourseModel {
	return &defaultSelectionCourseModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "selection_course",
	}
}

func (m *defaultSelectionCourseModel) Insert(data SelectionCourse) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?)", m.table, selectionCourseRowsExpectAutoSet)
	ret, err := m.ExecNoCache(query, data.SelectionId, data.CourseId, data.TeacherId)

	return ret, err
}

func (m *defaultSelectionCourseModel) FindOne(id int64) (*SelectionCourse, error) {
	selectionCourseIdKey := fmt.Sprintf("%s%v", cacheSelectionCourseIdPrefix, id)
	var resp SelectionCourse
	err := m.QueryRow(&resp, selectionCourseIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where id = ? limit 1", selectionCourseRows, m.table)
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

func (m *defaultSelectionCourseModel) Update(data SelectionCourse) error {
	selectionCourseIdKey := fmt.Sprintf("%s%v", cacheSelectionCourseIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where id = ?", m.table, selectionCourseRowsWithPlaceHolder)
		return conn.Exec(query, data.SelectionId, data.CourseId, data.TeacherId, data.Id)
	}, selectionCourseIdKey)
	return err
}

func (m *defaultSelectionCourseModel) Delete(id int64) error {

	selectionCourseIdKey := fmt.Sprintf("%s%v", cacheSelectionCourseIdPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where id = ?", m.table)
		return conn.Exec(query, id)
	}, selectionCourseIdKey)
	return err
}

func (m *defaultSelectionCourseModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheSelectionCourseIdPrefix, primary)
}

func (m *defaultSelectionCourseModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where id = ? limit 1", selectionCourseRows, m.table)
	return conn.QueryRow(v, query, primary)
}
