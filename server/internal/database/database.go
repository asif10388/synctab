package database

import (
	"context"
	"fmt"

	env "github.com/asif10388/synctab/internal/environment"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func (db *Database) Init(ctx context.Context, afterConnectCb AfterConnectFn) error {
	e := env.GetEnvironment()
	log.Info().Msg("initializing connection to configdb server")

	host := e.GetStrEnv("SYNCTAB_CONFIGDB_HOST")
	port := e.GetStrEnv("SYNCTAB_CONFIGDB_PORT")
	user := e.GetStrEnv("SYNCTAB_CONFIGDB_USER")
	secure := e.GetStrEnv("SYNCTAB_CONFIGDB_SECURE")
	dbName := e.GetStrEnv("SYNCTAB_CONFIGDB_DBNAME")
	password := e.GetStrEnv("SYNCTAB_CONFIGDB_PASSWORD")
	protocol := e.GetStrEnv("SYNCTAB_CONFIGDB_PROTOCOL")
	maxConns := e.GetStrEnv("SYNCTAB_CONFIGDB_MAX_CONNS")

	dbUrl := protocol + "://" + user + ":" + password + "@" + host + ":" + port + "/" + dbName + "?sslmode=" + secure + "&pool_max_conns=" + maxConns

	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse connection string")
		return fmt.Errorf("failed to parse connection string: %v", err)
	}

	if afterConnectCb != nil {
		config.AfterConnect = afterConnectCb
	}

	cpool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to configuration database server")
		return fmt.Errorf("failed to connect to Postgres server: %v", err)
	}

	db.CPool = cpool

	var greeting string
	err = cpool.QueryRow(ctx, "select 'Hello World!'").Scan(&greeting)
	if err != nil {
		log.Error().Err(err).Msg("connection test to configuration database server failed")
	}

	log.Info().Msg("successfully connected to configuration database")

	return nil
}

func (db *Database) Close(ctx context.Context) {
	if db != nil && db.CPool != nil {
		db.CPool.Close()
	}
}
