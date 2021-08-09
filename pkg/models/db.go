package models

import (
	"context"

	"github.com/jackc/pgx"
)

// Queryer is one of the pgx's queryable interfaces
type Queryer interface {
	ExecEx(context.Context, string, *pgx.QueryExOptions, ...interface{}) (pgx.CommandTag, error)
	QueryRowEx(context.Context, string, *pgx.QueryExOptions, ...interface{}) *pgx.Row
	QueryEx(context.Context, string, *pgx.QueryExOptions, ...interface{}) (*pgx.Rows, error)
}

// DB is just an interface that both pgx.Conn and pgx.ConnPool satisfy
type DB interface {
	Queryer
	Begin() (*pgx.Tx, error)
}
