package auth

import (
	"fmt"

	"github.com/asif10388/synctab/apiserver/controller"
	"github.com/asif10388/synctab/apiserver/model"
	"github.com/asif10388/synctab/apiserver/model/database"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type User struct {
	Id           string `json:"id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	PasswordHash string `json:"passwordhash"`
	model.ModelCommon
}

func init() {
	// userFields := "id, email, username, password, created_at, updated_at, deleted_at"

	userTemplates := map[string]string{
		"add_user_v1": "select _id from main.add_user_v1($1, $2, $3)",
	}

	database.NewStatements().AddSchemaTemplateMap(userTemplates)
}

func (auth *Auth) checkRegisterCreds(creds *UserRegisterCredentials) error {
	if len(creds.Username) == 0 || len(creds.Username) > auth.MaxUsernameLen {
		return controller.ErrInvalidCreds
	}

	if len(creds.Email) == 0 || len(creds.Email) > auth.MaxEmailLen {
		return controller.ErrInvalidCreds
	}

	if len(creds.Password) < auth.MinPasswordLen || len(creds.Password) > auth.MaxPasswordLen {
		return controller.ErrInvalidCreds
	}

	return nil
}

func (auth *Auth) getRegisterCreds(ctx *gin.Context) (*UserRegisterCredentials, error) {
	registerCreds := &UserRegisterCredentials{}

	err := ctx.ShouldBindJSON(registerCreds)
	if err != nil {
		return nil, err
	}

	err = auth.checkRegisterCreds(registerCreds)
	if err != nil {
		return nil, err
	}

	return registerCreds, nil
}

func (user *User) create(ctx *gin.Context, auth *Auth, tx pgx.Tx) (err error) {
	fmt.Println(database.Statements{})

	err = tx.QueryRow(ctx, "add_user_v1", user.Email, user.Username, user.PasswordHash).Scan(&user.Id)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return model.ErrUserAddFailed
	}

	return nil
}

func (auth *Auth) UserRegister(ctx *gin.Context) error {
	user := &User{}

	err := ctx.ShouldBindJSON(user)
	if err != nil {
		return err
	}

	tx, err := auth.Database.CPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Error().Err(err).Msg("failed to start transaction")
		return controller.ErrInternal
	}

	defer func() {
		if err != nil {
			txErr := tx.Rollback(ctx)
			if txErr != nil {
				log.Error().Err(txErr).Msg("failed to rollback transaction")
			}
		} else {
			txErr := tx.Commit(ctx)
			if txErr != nil {
				log.Error().Err(txErr).Msg("failed to commit transaction")
				err = txErr
			} else {
				log.Info().Msgf("successfully created user")
			}
		}
	}()

	err = user.create(ctx, auth, tx)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return err
	}

	return nil
}
