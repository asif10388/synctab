package model

import (
	"context"
	"regexp"
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
		PrimarySchemaName: model.Env.GetStrEnv("SYNCTAB_ACCOUNT_PRIMARY_SCHEMA_NAME"),

		MaxIPLen:        model.Env.GetIntEnv("SYNCTAB_MAXIPLEN"),
		MaxUrlLen:       model.Env.GetIntEnv("SYNCTAB_MAXURLLEN"),
		MaxNameLen:      model.Env.GetIntEnv("SYNCTAB_MAXNAMELEN"),
		MaxEmailLen:     model.Env.GetIntEnv("SYNCTAB_MAXEMAILLEN"),
		MaxSchemaLen:    model.Env.GetIntEnv("SYNCTAB_MAXSCHEMALEN"),
		MaxUsernameLen:  model.Env.GetIntEnv("SYNCTAB_MAXUSERNAMELEN"),
		MaxPasswordLen:  model.Env.GetIntEnv("SYNCTAB_MAXPASSWORDLEN"),
		MinPasswordLen:  model.Env.GetIntEnv("SYNCTAB_MINPASSWORDLEN"),
		MaxHostnameLen:  model.Env.GetIntEnv("SYNCTAB_MAXHOSTNAMELEN"),
		MaxEndpointsLen: model.Env.GetIntEnv("SYNCTAB_MAXENDPOINTSLEN"),

		EndpointsSeparator: model.Env.GetStrEnv("SYNCTAB_ENDPOINTS_SEPARATOR"),
		HostPortSeparator:  model.Env.GetStrEnv("SYNCTAB_HOSTPORT_SEPARATOR"),

		HostnamePattern: model.Env.GetStrEnv("SYNCTAB_HOSTNAME_VALIDATION_PATTERN"),

		LoginIdPattern: model.Env.GetStrEnv("SYNCTAB_USER_LOGINID_PATTERN"),

		MaxPaginationEntries: model.Env.GetIntEnv("SYNCTAB_MAXPAGINATION_ENTRIES"),
		UserTimeout:          model.Env.GetDurationEnv("SYNCTAB_USER_IDLE_TIMEOUT"),

		TokenExpiresDuration: model.Env.GetDurationEnv("SYNCTAB_JWT_EXPIRATION"),
		TokenUpdateDuration:  model.Env.GetDurationEnv("SYNCTAB_JWT_RENEWAL_BEFORE"),
	}

	hostR, err := regexp.Compile(model.HostnamePattern)
	if err != nil {
		log.Error().Err(err).Msgf("failed to compile hostname pattern %s", model.HostnamePattern)
		return nil, err
	}
	model.HostnameRegexp = hostR

	loginR, err := regexp.Compile(model.LoginIdPattern)
	if err != nil {
		log.Error().Err(err).Msgf("failed to compile login id pattern %s", model.LoginIdPattern)
		return nil, err
	}
	model.LoginIdRegexp = loginR

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
		PrimarySchemaName: model.PrimarySchemaName,
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
