package database

import (
	"errors"

	"github.com/asif10388/synctab/internal/database"
)

var (
	ErrInvalidSchema           = errors.New("invalid schema name")
	ErrUnknownStatement        = errors.New("unknown statement")
	ErrStatementsInvalid       = errors.New("invalid statements instance")
	ErrStatementsUninitialized = errors.New("database statements uninitialized")
)

type Input struct {
	PrimarySchemaName string
}

type Database struct {
	database.Database
	Input
}

type Statements struct {
	primarySchemaTemplates map[string]string
	primarySchemaSql       map[string]string
}
