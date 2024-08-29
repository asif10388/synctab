package database

import (
	"bufio"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	env "github.com/asif10388/synctab/internal/environment"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type schemaSql struct {
	SchemaName string
	statement  []byte
}

func (db *Database) SetupSchema(ctx context.Context) error {
	templateStatements, err := db.getSqlStatementsFromDir()
	if err != nil {
		log.Error().Err(err).Msgf("failed to get statements from directory %s", env.ApiServiceSqlDir)
		return err
	}

	if len(templateStatements) == 0 {
		log.Error().Msg("did not find any SQL statements to initialize schema")
		return nil
	}

	statements := make([]string, len(templateStatements))
	for i, templateStatement := range templateStatements {
		statement, err := db.getStatementFromTemplateString(db.PrimarySchemaName, templateStatement)
		if err != nil {
			log.Error().Err(err).Msgf("failed to get sql statement from template %s", templateStatement)
			return err
		}

		statements[i] = statement
	}

	err = db.execStatements(ctx, statements)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute statements to setup schema")
		return err
	}

	log.Info().Msg("successfully setup schema")
	return nil
}

func (db *Database) PrepareSchema(ctx context.Context, conn *pgx.Conn) error {
	statements := NewStatements()

	if len(statements.schemaSql) == 0 {
		log.Info().Msgf("creating statements to prepare from template in schema %s", db.PrimarySchemaName)

		schemaSql, err := db.getStatementsFromMap(db.PrimarySchemaName, statements.schemaTemplates)
		if err != nil {
			log.Error().Err(err).Msg("failed to get prepared statements for schema")
			return err
		}

		statements.schemaSql = schemaSql

	} else {
		log.Info().Msg("using cached prepared statements for schema")
	}

	if len(statements.schemaSql) == 0 {
		log.Warn().Msg("did not find any prepared statements for schema")
		return nil
	}

	for name, statement := range statements.schemaTemplates {
		_, err := conn.Prepare(ctx, name, statement)
		if err != nil {
			log.Error().Err(err).Msgf("failed to prepare statement %s => %s for schema %s", name, statement, db.PrimarySchemaName)
			return err
		}
	}

	log.Info().Msgf("successfully prepared statements for schema %s in new connection", db.PrimarySchemaName)
	return nil
}

func (db *Database) getSqlStatementsFromDir() ([]string, error) {
	pwd, pwdErr := os.Getwd()
	if pwdErr != nil {
		panic(pwdErr)
	}

	sqlFiles, err := filepath.Glob(filepath.Join(pwd+env.ApiServiceSqlDir) + "/*.sql")
	if err != nil {
		return nil, err
	}

	statements := []string{}
	for _, sqlFile := range sqlFiles {
		fileStatements, err := db.getSqlStatementsFromFile(sqlFile)
		if err != nil {
			log.Error().Err(err).Msgf("failed to get sql statements from file %s", sqlFile)
			continue
		}

		statements = append(statements, fileStatements...)
	}

	return statements, nil

}

func (db *Database) getSqlStatementsFromFile(file string) ([]string, error) {
	statements := []string{}

	fileContent, err := os.Open(file)
	if err != nil {
		log.Error().Err(err).Msgf("failed to open file %s", file)
		return nil, err
	}

	reader := bufio.NewReader(fileContent)
	statement := ""

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			break
		}

		if len(line) > 0 {
			trimmedLine := strings.TrimSpace(string(line))
			if trimmedLine == sqlFileDelimiter {
				statements = append(statements, statement)
				statement = ""
			} else {
				statement += string(line)
			}
		}

		if err == io.EOF {
			if len(statement) > 0 {
				statements = append(statements, statement)
			}

			break
		}
	}

	return statements, nil
}

func (statementSql *schemaSql) Write(s []byte) (int, error) {
	if len(s) == 0 {
		return 0, nil
	}

	d := make([]byte, len(s))
	clen := copy(d, s)
	statementSql.statement = append(statementSql.statement, d...)
	return clen, nil
}

func (db *Database) getStatementsFromMap(schemaName string, templateStatements map[string]string) (map[string]string, error) {
	prepareStatements := make(map[string]string)

	for statementName, statementValue := range templateStatements {
		nameSql, err := db.getStatementFromTemplateString(schemaName, statementName)
		if err != nil {
			log.Error().Err(err).Msgf("failed to get sql statement name %s from map", statementName)
			return nil, err
		}

		statementSql, err := db.getStatementFromTemplateString(schemaName, statementValue)
		if err != nil {
			log.Error().Err(err).Msgf("failed to get sql statement %s from map", statementValue)
			return nil, err
		}

		prepareStatements[nameSql] = string(statementSql)
	}

	return prepareStatements, nil
}

func (db *Database) getStatementFromTemplateString(schemaName string, templateString string) (string, error) {
	if len(templateString) == 0 {
		return "", nil
	}

	newTemplate := template.New("sqlStatement")
	sqlStatement := &schemaSql{SchemaName: schemaName}

	sqlStatementTemplate, err := newTemplate.Parse(templateString)
	if err != nil {
		log.Error().Err(err).Msgf("failed to parse template string %s", templateString)
		return "", err
	}

	err = sqlStatementTemplate.Execute(sqlStatement, sqlStatement)
	if err != nil {
		log.Error().Err(err).Msgf("failed to get sql statement from template string %s", templateString)
		return "", err
	}

	return string(sqlStatement.statement), nil
}
