package isql

import (
	"context"
	"database/sql"
)

func NewDB(d *sql.DB) DB {
	return &db{
		DB: d,
	}
}

var _ DB = (*db)(nil)

type db struct {
	*sql.DB
}

func (db *db) Query(query string, args ...interface{}) (Rows, error) {
	return db.DB.Query(query, args...)
}

func (db *db) QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	return db.DB.QueryContext(ctx, query, args...)
}

func (db *db) QueryRow(query string, args ...interface{}) Row {
	return db.DB.QueryRow(query, args...)
}

func (db *db) QueryRowContext(ctx context.Context, query string, args ...interface{}) Row {
	return db.DB.QueryRowContext(ctx, query, args...)
}

// type

type Bytes []byte

type BytesArray [][]byte

type OfflineBool uint32

const (
	DontPush    OfflineBool = 0 // 推送失败过
)
