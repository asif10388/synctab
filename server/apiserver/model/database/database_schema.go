package database

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	env "github.com/asif10388/synctab/internal/environment"
	"github.com/rs/zerolog/log"
)

type schemaSql struct {
	SchemaName string
	statement  []byte
}

func (db *Database) SetupPrimarySchema(ctx context.Context) error {
	templateStatements, err := db.getSqlStatementsFromDir(env.ApiServiceSqlDir)
	if err != nil {
		log.Error().Err(err).Msgf("failed to get statements from directory %s", env.ApiServiceSqlPrimaryDir)
		return err
	}

	if len(templateStatements) == 0 {
		log.Error().Msg("did not find any SQL statements to initialize primary schema")
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
		log.Error().Err(err).Msg("failed to execute statements to setup primary schema")
		return err
	}

	log.Info().Msg("successfully setup primary schema")
	return nil
}

func (db *Database) getSqlStatementsFromDir(dir string) ([]string, error) {
	if !strings.HasPrefix(dir, env.ApiServiceSqlDir) {
		log.Error().Msgf("%s is not a sub-directory of %s", dir, env.ApiServiceSqlDir)
		return nil, fmt.Errorf("%s is not a sub-directory of %s", dir, env.ApiServiceSqlDir)
	}

	pwd, pwdErr := os.Getwd()
	if pwdErr != nil {
		panic(pwdErr)
	}

	sqlFiles, err := filepath.Glob(filepath.Join(pwd+"/apiserver/model/sql/synctabdb") + "/*.sql")
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

	fmt.Println(len(statements))
	return statements, nil

}

func (db *Database) getSqlStatementsFromFile(file string) ([]string, error) {
	// if !strings.HasPrefix(file, env.ApiServiceSqlDir) {
	// 	log.Error().Msgf("file %s is not inside Apiservice sql directory %s", file, env.ApiServiceSqlDir)
	// 	return nil, fmt.Errorf("%s is not a file under %s", file, env.ApiServiceSqlDir)
	// }

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
