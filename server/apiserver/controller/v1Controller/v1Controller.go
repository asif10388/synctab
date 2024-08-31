package v1controller

import (
	"fmt"
	"net/http"

	controller "github.com/asif10388/synctab/apiserver/controller"
	controllerPkg "github.com/asif10388/synctab/apiserver/controller"
	authcontrollerPkg "github.com/asif10388/synctab/apiserver/controller/v1Controller/authcontroller"
	urlscontrollerPkg "github.com/asif10388/synctab/apiserver/controller/v1Controller/urlscontroller"
	middlewarePkg "github.com/asif10388/synctab/apiserver/middleware"
	modelPkg "github.com/asif10388/synctab/apiserver/model"

	envPkg "github.com/asif10388/synctab/internal/environment"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var v1Controller *V1Controller

func NewV1Controller(controller *controllerPkg.Controller) (controllerPkg.ApiVersion, error) {
	if v1Controller == nil {

		middleware, err := middlewarePkg.NewMiddleware()
		if err != nil {
			return nil, err
		}

		model, err := modelPkg.NewModel()
		if err != nil {
			return nil, err
		}

		authController, err := authcontrollerPkg.NewAuthController(controller, model, middleware)
		if err != nil {
			log.Error().Err(err).Msg("failed to initialize v1 auth controller")
			return nil, err
		}

		urlsController, err := urlscontrollerPkg.NewUrlsController(controller, model, middleware)
		if err != nil {
			log.Error().Err(err).Msg("failed to initialize v1 urls controller")
			return nil, err
		}

		v1Controller = &V1Controller{
			Model:          model,
			Controller:     controller,
			Middleware:     middleware,
			AuthController: authController,
			UrlsController: urlsController,
			Env:            envPkg.GetEnvironment(),
		}
	}

	return v1Controller, nil
}

func (v1Controller *V1Controller) Init() error {
	log.Info().Msg("initializing v1 APIs")

	err := v1Controller.Model.Init()
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize model")
		return fmt.Errorf("failed to initialize model: %w", err)
	}

	err = v1Controller.Middleware.Init(v1Controller.Model)
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize middleware")
		return fmt.Errorf("failed to initialize middleware: %w", err)
	}

	apiPrefix := v1Controller.Env.GetStrEnv("SYNCTAB_API_PREFIX") + "/"
	apiPrefix += ApiVersion

	publicGroup := v1Controller.Controller.Engine.Group(apiPrefix, cors.New(v1Controller.Controller.GetCorsConfig()))

	publicGroup.GET(controller.PingPath, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	v1Controller.AuthController.Init(publicGroup)
	v1Controller.UrlsController.Init(publicGroup)

	return nil
}

func (v1Controller *V1Controller) Start() error {
	log.Info().Msg("starting model")
	v1Controller.Model.Start()

	log.Info().Msg("starting middleware")
	v1Controller.Middleware.Start()

	return nil
}
