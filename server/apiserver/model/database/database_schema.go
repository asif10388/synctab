package database

import (
	"bufio"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	env "github.com/asif10388/synctab/internal/environment"
	"github.com/rs/zerolog/log"
)

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

	err = db.execStatements(ctx, templateStatements)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute statements to setup schema")
		return err
	}

	log.Info().Msg("successfully setup schema")
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
