package auth

import (
	"context"
	"fmt"
	"time"

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
	Password     string `json:"password"`
	PasswordHash string `json:"passwordhash"`
	model.ModelCommon
}

func init() {
	userFields := "email, username, passwordhash"

	userTemplates := map[string]string{
		"add_user_v1":          "select _id from main.add_user_v1($1, $2, $3)",
		"get_user_by_email_v1": fmt.Sprintf("select %s from main.users where email = $1", userFields),
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
	addUserSQL := database.NewStatements().GetSchemaTemplate("add_user_v1")
	if addUserSQL == "" {
		return fmt.Errorf("add_user_v1 SQL template not found")
	}

	err = tx.QueryRow(ctx, addUserSQL, user.Email, user.Username, user.PasswordHash).Scan(&user.Id)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return model.ErrUserAddFailed
	}

	return nil
}

func (auth *Auth) CreateUser(ctx *gin.Context) error {
	user := &User{}

	err := ctx.ShouldBindJSON(user)
	if err != nil {
		return err
	}

	passwordHash, err := auth.GetPasswordHash(user.Password)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create a password hash")
		return err
	}

	user.PasswordHash = passwordHash

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

func (user *User) login(ctx context.Context, auth *Auth, tx pgx.Tx) error {
	getUserSQL := database.NewStatements().GetSchemaTemplate("get_user_by_email_v1")
	if getUserSQL == "" {
		return fmt.Errorf("get_user_by_email_v1 SQL template not found")
	}

	err := tx.QueryRow(ctx, getUserSQL, user.Email).Scan(
		&user.Email,
		&user.Username,
		&user.PasswordHash,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return controller.ErrInvalidCreds
		}

		return fmt.Errorf("failed to fetch user: %w", err)
	}

	return nil

}

func (auth *Auth) LoginUser(ctx *gin.Context) (*LoginResponse, error) {

	user := &User{}
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		return nil, controller.ErrInvalidInput
	}

	tx, err := auth.Database.CPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Error().Err(err).Msg("failed to start transaction")
		return nil, controller.ErrInternal
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

	err = user.login(ctx, auth, tx)
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch user")
		return nil, err
	}

	if !auth.CompareAdaptiveHashString(user.PasswordHash, user.Password) {
		log.Error().Err(model.ErrUserAuthFailed).Msg("user authentication failed")
		return nil, model.ErrUserAuthFailed
	}

	claims := Claims{
		Version: 1,
		Email:   user.Email,
	}

	creds := Credentials{
		Email:    user.Email,
		Password: user.Password,
	}

	token, err := auth.GetJWTToken(ctx, &creds, &claims, time.Now().UTC())
	if err != nil {
		log.Error().Err(err).Msg("failed to get JWT")
		return nil, err
	}

	response := LoginResponse{
		Email:    user.Email,
		Username: user.Username,
		Token:    token,
	}

	log.Info().Str("username", user.Username).Msg("user logged in successfully")
	return &response, nil
}
