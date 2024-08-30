package middleware

import (
	"github.com/asif10388/synctab/apiserver/controller"
	"github.com/asif10388/synctab/apiserver/model"
	"github.com/asif10388/synctab/apiserver/model/auth"
	env "github.com/asif10388/synctab/internal/environment"
	"github.com/rs/zerolog/log"
)

var middlewareInstance *Middleware

func NewMiddleware() (*Middleware, error) {
	if middlewareInstance != nil {
		return middlewareInstance, nil
	}

	middleware := &Middleware{
		Env: env.GetEnvironment(),
	}

	if middleware.Env == nil {
		log.Error().Msg("failed to initialize environment")
		return nil, controller.ErrEnvironmentNotInitialized
	}

	log.Info().Msg("successfully initialized middleware")
	middlewareInstance = middleware
	return middlewareInstance, nil
}

func (middleware *Middleware) Init(model *model.Model) error {
	middleware.Model = model

	authModel, err := auth.NewAuthModel(model)
	if err != nil {
		return err
	}

	middleware.Auth = authModel

	return nil
}

func (middleware *Middleware) SetController(controller *controller.Controller) {
	middleware.Controller = controller
}

func (middleware *Middleware) Start() error {
	return nil
}

func (middleware *Middleware) Stop() error {
	return nil
}
