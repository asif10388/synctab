package database

import (
	"context"
	"fmt"
	"time"

	database "github.com/asif10388/synctab/internal/database"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

var dbInstance *Database

const (
	connrepeats = 5
	connbackoff = 5
)

func Init(ctx context.Context, input Input) (*Database, error) {
	if dbInstance == nil {
		db := &Database{
			Database: database.Database{},
			Input:    input,
		}

		var err error
		for i := 0; i < connrepeats; i++ {
			err = db.Database.Init(ctx, db.afterConnect)
			if err == nil {
				break
			}

			log.Error().Err(err).Msgf("failed to connect to database. will attempt again after %d backoff", connbackoff)
			time.Sleep(connbackoff * time.Second)
		}

		if err != nil {
			log.Error().Err(err).Msgf("failed to connect to database after %d attempts. giving up", connrepeats)
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}

		log.Logger.Info().Msg("setting up schema post connection")
		err = db.SetupSchema(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to setup schema in database")
			return nil, fmt.Errorf("failed to setup schema in database: %w", err)
		}

		dbInstance = db
	}

	return dbInstance, nil
}

func (db *Database) afterConnect(ctx context.Context, conn *pgx.Conn) error {
	err := db.PrepareSchema(ctx, conn)
	if err != nil {
		log.Error().Err(err).Msg("failed to setup prepared statements for schema in new connection to configuration database")
		return err
	}

	return nil
}

func (db *Database) Close(ctx context.Context) error {
	db.Database.Close(ctx)
	return nil
}

func (db *Database) execStatementsInTx(ctx context.Context, statements []string, tx pgx.Tx) (err error) {
	for _, statement := range statements {
		_, err := tx.Exec(ctx, statement)
		if err != nil {
			log.Error().Err(err).Msgf("failed to execute sql statement %s", statement)
			break
		}
	}

	return err
}

func (db *Database) execStatements(ctx context.Context, statements []string) (err error) {
	log.Info().Msg("creating transaction to execute statements")

	var tx pgx.Tx

	tx, err = db.CPool.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to start transaction for executing statements")
		return fmt.Errorf("failed to start transaction, error %w", err)
	}

	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("failed to execute statements")
			txErr := tx.Rollback(ctx)
			if txErr != nil {
				log.Error().Err(txErr).Msg("failed to rollback transaction while executing statements")
			}
		} else {
			txErr := tx.Commit(ctx)
			if txErr != nil {
				log.Error().Err(txErr).Msg("failed to commit transaction while executing statements")
				err = txErr
			}
		}
	}()

	err = db.execStatementsInTx(ctx, statements, tx)
	return err
}
