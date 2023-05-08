package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Option func(db *sql.DB)

func WithMaxIdleConns(num int) Option {
	return func(db *sql.DB) {
		db.SetMaxIdleConns(num)
	}
}

func WithMaxOpenConns(num int) Option {
	return func(db *sql.DB) {
		db.SetMaxOpenConns(num)
	}
}

func WithConnMaxLifetime(d time.Duration) Option {
	return func(db *sql.DB) {
		db.SetConnMaxLifetime(d)
	}
}

func OpenDB(ds Datasource, opt ...Option) (*sql.DB, error) {
	db, err := sql.Open(ds.Driver(), ds.DSN())
	if err != nil {
		return nil, fmt.Errorf("database connection: %w", err)
	}

	for _, o := range opt {
		o(db)
	}
	return db, nil
}
