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

type UserModel struct {
	Id           string `json:"id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	PasswordHash string `json:"passwordhash"`
	model.ModelCommon
}

func init() {
	userFields := "id, email, username, passwordhash"

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

func (userModel *UserModel) create(ctx *gin.Context, auth *Auth, tx pgx.Tx) (err error) {
	addUserSQL := database.NewStatements().GetSchemaTemplate("add_user_v1")
	if addUserSQL == "" {
		return fmt.Errorf("add_user_v1 SQL template not found")
	}

	err = tx.QueryRow(ctx, addUserSQL, userModel.Email, userModel.Username, userModel.PasswordHash).Scan(&userModel.Id)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return model.ErrUserAddFailed
	}

	return nil
}

func (auth *Auth) CreateUser(ctx *gin.Context) error {
	userModel := &UserModel{}

	err := ctx.ShouldBindJSON(userModel)
	if err != nil {
		return err
	}

	passwordHash, err := auth.GetPasswordHash(userModel.Password)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create a password hash")
		return err
	}

	userModel.PasswordHash = passwordHash

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

	err = userModel.create(ctx, auth, tx)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return err
	}

	return nil
}

func (userModel *UserModel) login(ctx context.Context, auth *Auth, tx pgx.Tx) error {
	getUserSQL := database.NewStatements().GetSchemaTemplate("get_user_by_email_v1")
	if getUserSQL == "" {
		return fmt.Errorf("get_user_by_email_v1 SQL function not found")
	}

	err := tx.QueryRow(ctx, getUserSQL, userModel.Email).Scan(
		&userModel.Id,
		&userModel.Email,
		&userModel.Username,
		&userModel.PasswordHash,
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
	userModel := &UserModel{}

	err := ctx.ShouldBindJSON(&userModel)
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

	err = userModel.login(ctx, auth, tx)
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch user")
		return nil, err
	}

	if !auth.CompareAdaptiveHashString(userModel.PasswordHash, userModel.Password) {
		log.Error().Err(model.ErrUserAuthFailed).Msg("user authentication failed")
		return nil, model.ErrUserAuthFailed
	}

	claims := Claims{
		Version: 1,
		Email:   userModel.Email,
	}

	creds := Credentials{
		Email:    userModel.Email,
		Password: userModel.Password,
	}

	token, err := auth.GetJWTToken(ctx, &creds, &claims, time.Now().UTC())
	if err != nil {
		log.Error().Err(err).Msg("failed to get JWT")
		return nil, err
	}

	response := LoginResponse{
		Token:    token,
		Email:    userModel.Email,
		Username: userModel.Username,
	}

	auth.UserId = userModel.Id

	log.Info().Str("user_id", auth.UserId).Msg("user logged in successfully")
	return &response, nil
}
