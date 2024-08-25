package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	CPool *pgxpool.Pool
}

type AfterConnectFn func(context.Context, *pgx.Conn) error
