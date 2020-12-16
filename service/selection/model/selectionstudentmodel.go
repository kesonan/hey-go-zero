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
	selectionStudentFieldNames          = builderx.FieldNames(&SelectionStudent{})
	selectionStudentRows                = strings.Join(selectionStudentFieldNames, ",")
	selectionStudentRowsExpectAutoSet   = strings.Join(stringx.Remove(selectionStudentFieldNames, "id", "create_time", "update_time"), ",")
	selectionStudentRowsWithPlaceHolder = strings.Join(stringx.Remove(selectionStudentFieldNames, "id", "create_time", "update_time"), "=?,") + "=?"

	cacheSelectionStudentIdPrefix = "cache#SelectionStudent#id#"
)

type (
	SelectionStudentModel interface {
		Insert(data SelectionStudent) (sql.Result, error)
		FindOne(id int64) (*SelectionStudent, error)
		Update(data SelectionStudent) error
		Delete(id int64) error
		FindByStudentId(studentId int64) ([]*SelectionStudent, error)
		FindBySelectionCourseId(selectionCourseId int64) ([]*SelectionStudent, error)
		FindByStudentIdAndSelectionCourseId(studentId, selectionCourseId int64) (*SelectionStudent, error)
	}

	defaultSelectionStudentModel struct {
		sqlc.CachedConn
		table string
	}

	SelectionStudent struct {
		CreateTime        time.Time `db:"create_time"`
		UpdateTime        time.Time `db:"update_time"`
		Id                int64     `db:"id"`
		SelectionCourseId int64     `db:"selection_course_id"` // 选课任务课程id
		StudentId         int64     `db:"student_id"`          // 学生id
	}
)

func NewSelectionStudentModel(conn sqlx.SqlConn, c cache.CacheConf) SelectionStudentModel {
	return &defaultSelectionStudentModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "selection_student",
	}
}

func (m *defaultSelectionStudentModel) Insert(data SelectionStudent) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, selectionStudentRowsExpectAutoSet)
	ret, err := m.ExecNoCache(query, data.SelectionCourseId, data.StudentId)

	return ret, err
}

func (m *defaultSelectionStudentModel) FindOne(id int64) (*SelectionStudent, error) {
	selectionStudentIdKey := fmt.Sprintf("%s%v", cacheSelectionStudentIdPrefix, id)
	var resp SelectionStudent
	err := m.QueryRow(&resp, selectionStudentIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where id = ? limit 1", selectionStudentRows, m.table)
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

func (m *defaultSelectionStudentModel) Update(data SelectionStudent) error {
	selectionStudentIdKey := fmt.Sprintf("%s%v", cacheSelectionStudentIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where id = ?", m.table, selectionStudentRowsWithPlaceHolder)
		return conn.Exec(query, data.SelectionCourseId, data.StudentId, data.Id)
	}, selectionStudentIdKey)
	return err
}

func (m *defaultSelectionStudentModel) Delete(id int64) error {

	selectionStudentIdKey := fmt.Sprintf("%s%v", cacheSelectionStudentIdPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where id = ?", m.table)
		return conn.Exec(query, id)
	}, selectionStudentIdKey)
	return err
}

func (m *defaultSelectionStudentModel) FindByStudentId(studentId int64) ([]*SelectionStudent, error) {
	query := fmt.Sprintf("select %s from %s where student_id = ?", selectionStudentRows, m.table)
	var resp []*SelectionStudent
	err := m.QueryRowsNoCache(&resp, query, studentId)
	return resp, err
}

func (m *defaultSelectionStudentModel) FindBySelectionCourseId(selectionCourseId int64) ([]*SelectionStudent, error) {
	query := fmt.Sprintf("select %s from %s where selection_course_id = ?", selectionStudentRows, m.table)
	var resp []*SelectionStudent
	err := m.QueryRowNoCache(&resp, query, selectionCourseId)
	return resp, err
}

func (m *defaultSelectionStudentModel) FindByStudentIdAndSelectionCourseId(studentId, selectionCourseId int64) (*SelectionStudent, error) {
	query := fmt.Sprintf("select %s from %s where student_id = ? and selection_course_id = ? limit 1", selectionStudentRows, m.table)
	var resp SelectionStudent
	err := m.QueryRowNoCache(&resp, query, studentId, selectionCourseId)
	return &resp, err
}

func (m *defaultSelectionStudentModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheSelectionStudentIdPrefix, primary)
}

func (m *defaultSelectionStudentModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where id = ? limit 1", selectionStudentRows, m.table)
	return conn.QueryRow(v, query, primary)
}
