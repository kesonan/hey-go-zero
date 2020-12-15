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
	selectionFieldNames          = builderx.FieldNames(&Selection{})
	selectionRows                = strings.Join(selectionFieldNames, ",")
	selectionRowsExpectAutoSet   = strings.Join(stringx.Remove(selectionFieldNames, "id", "create_time", "update_time"), ",")
	selectionRowsWithPlaceHolder = strings.Join(stringx.Remove(selectionFieldNames, "id", "create_time", "update_time"), "=?,") + "=?"

	cacheSelectionNamePrefix = "cache#Selection#name#"
	cacheSelectionIdPrefix   = "cache#Selection#id#"
)

type (
	SelectionModel interface {
		Insert(data Selection) (sql.Result, error)
		FindOne(id int64) (*Selection, error)
		FindOneByName(name string) (*Selection, error)
		Update(data Selection) error
		Delete(id int64) error
	}

	defaultSelectionModel struct {
		sqlc.CachedConn
		table string
	}

	Selection struct {
		MaxCredit    int64     `db:"max_credit"`   // 最大可修学分
		StartTime    int64     `db:"start_time"`   // 选课开始时间
		EndTime      int64     `db:"end_time"`     // 选课结束时间
		Notification string    `db:"notification"` // 选课通知内容，500字以内
		CreateTime   time.Time `db:"create_time"`
		UpdateTime   time.Time `db:"update_time"`
		Id           int64     `db:"id"`
		Name         string    `db:"name"` // 选课名称
	}
)

func NewSelectionModel(conn sqlx.SqlConn, c cache.CacheConf) SelectionModel {
	return &defaultSelectionModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "selection",
	}
}

func (m *defaultSelectionModel) Insert(data Selection) (sql.Result, error) {
	selectionNameKey := fmt.Sprintf("%s%v", cacheSelectionNamePrefix, data.Name)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, selectionRowsExpectAutoSet)
		return conn.Exec(query, data.MaxCredit, data.StartTime, data.EndTime, data.Notification, data.Name)
	}, selectionNameKey)
	return ret, err
}

func (m *defaultSelectionModel) FindOne(id int64) (*Selection, error) {
	selectionIdKey := fmt.Sprintf("%s%v", cacheSelectionIdPrefix, id)
	var resp Selection
	err := m.QueryRow(&resp, selectionIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where id = ? limit 1", selectionRows, m.table)
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

func (m *defaultSelectionModel) FindOneByName(name string) (*Selection, error) {
	selectionNameKey := fmt.Sprintf("%s%v", cacheSelectionNamePrefix, name)
	var resp Selection
	err := m.QueryRowIndex(&resp, selectionNameKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where name = ? limit 1", selectionRows, m.table)
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

func (m *defaultSelectionModel) Update(data Selection) error {
	selectionIdKey := fmt.Sprintf("%s%v", cacheSelectionIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where id = ?", m.table, selectionRowsWithPlaceHolder)
		return conn.Exec(query, data.MaxCredit, data.StartTime, data.EndTime, data.Notification, data.Name, data.Id)
	}, selectionIdKey)
	return err
}

func (m *defaultSelectionModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	selectionNameKey := fmt.Sprintf("%s%v", cacheSelectionNamePrefix, data.Name)
	selectionIdKey := fmt.Sprintf("%s%v", cacheSelectionIdPrefix, id)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where id = ?", m.table)
		return conn.Exec(query, id)
	}, selectionNameKey, selectionIdKey)
	return err
}

func (m *defaultSelectionModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheSelectionIdPrefix, primary)
}

func (m *defaultSelectionModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where id = ? limit 1", selectionRows, m.table)
	return conn.QueryRow(v, query, primary)
}
