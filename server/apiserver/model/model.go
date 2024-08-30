package model

import (
	"context"
	"time"

	controller "github.com/asif10388/synctab/apiserver/controller"
	"github.com/asif10388/synctab/apiserver/model/database"
	env "github.com/asif10388/synctab/internal/environment"
	"github.com/rs/zerolog/log"
)

var modelInstance *Model

func NewModel() (*Model, error) {
	if modelInstance != nil {
		return modelInstance, nil
	}

	model := &Model{
		Env: env.GetEnvironment(),
	}

	if model.Env == nil {
		log.Error().Msg("failed to initialize environment")
		return nil, controller.ErrEnvironmentNotInitialized
	}

	model.ModelDefaults = ModelDefaults{
		SchemaName:  model.Env.GetStrEnv("SYNCTAB_ACCOUNT_PRIMARY_SCHEMA_NAME"),
		TokenSecret: model.Env.GetEnv("SYNCTAB_API_JWT_SECRET"),

		MaxUrlLen:   model.Env.GetIntEnv("SYNCTAB_MAXURLLEN"),
		MaxNameLen:  model.Env.GetIntEnv("SYNCTAB_MAXNAMELEN"),
		MaxEmailLen: model.Env.GetIntEnv("SYNCTAB_MAXEMAILLEN"),

		MaxUsernameLen: model.Env.GetIntEnv("SYNCTAB_MAXUSERNAMELEN"),
		MaxPasswordLen: model.Env.GetIntEnv("SYNCTAB_MAXPASSWORDLEN"),
		MinPasswordLen: model.Env.GetIntEnv("SYNCTAB_MINPASSWORDLEN"),

		MaxPaginationEntries: model.Env.GetIntEnv("SYNCTAB_MAXPAGINATION_ENTRIES"),
		UserTimeout:          model.Env.GetDurationEnv("SYNCTAB_USER_IDLE_TIMEOUT"),

		TokenExpiresDuration: model.Env.GetDurationEnv("SYNCTAB_JWT_EXPIRATION"),
		TokenUpdateDuration:  model.Env.GetDurationEnv("SYNCTAB_JWT_RENEWAL_BEFORE"),
	}

	model.UserTimeoutSecs = int64(model.UserTimeout / time.Second)
	model.TokenUpdateDurationSecs = int64(model.TokenUpdateDuration / time.Second)
	model.TokenExpiresDurationSecs = int64(model.TokenExpiresDuration / time.Second)

	modelInstance = model
	return modelInstance, nil
}

func (model *Model) Init() error {
	ctx := context.Background()

	log.Info().Msg("initializing configuration database")

	db, err := database.Init(ctx, database.Input{
		SchemaName: model.SchemaName,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize model configuration database handler")
	}

	model.Database = db

	log.Info().Msg("model infrastructure initialized")
	return nil
}

func (model *Model) SetController(controller *controller.Controller) {
	model.Controller = controller

}

func (model *Model) Start() error {
	return nil
}

func (model *Model) Stop() error {
	ctx := context.Background()
	err := model.Database.Close(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to close database connection")
	}

	return nil
}
