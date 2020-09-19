package isql

import (
	"context"
	"database/sql"
)

type DB interface {
	Query(query string, args ...interface{}) (Rows, error)

	QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error)

	QueryRow(query string, args ...interface{}) Row

	QueryRowContext(ctx context.Context, query string, args ...interface{}) Row

	Exec(query string, args ...interface{}) (sql.Result, error)

	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	Close() error
}

type Rows interface {
	Columns() ([]string, error)
	Next() bool
	Close() error
	Err() error
	Scan(dest ...interface{}) error
}

type Row interface {
	Scan(dest ...interface{}) error
}
