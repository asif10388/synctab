package database

import "context"

const (
	defaultPort     = 5432
	defaultUser     = "postgres"
	defaultPassword = "password"
	defaultHost     = "localhost"
	defaultDBName   = "synctabdb"
)

func (db *Database) Init(ctx context.Context)
