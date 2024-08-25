package v1controller

import (
	"fmt"
	"net/http"

	controller "github.com/asif10388/synctab/apiserver/controller"
	"github.com/asif10388/synctab/apiserver/controller/v1Controller/authcontroller"
	env "github.com/asif10388/synctab/internal/environment"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var v1Controller *V1Controller

func NewV1Controller(controller *controller.Controller) (controller.ApiVersion, error) {
	if v1Controller == nil {

		authController, err := authcontroller.NewAuthController(controller)
		if err != nil {
			log.Error().Err(err).Msg("failed to initialize v1 auth controller")
			return nil, err
		}

		v1Controller = &V1Controller{
			Controller:     controller,
			Env:            env.GetEnvironment(),
			AuthController: authController,
		}
	}

	return v1Controller, nil
}

func (v1Controller *V1Controller) Init() error {
	apiPrefix := v1Controller.Env.GetStrEnv("SYNCTAB_API_PREFIX") + "/"
	apiPrefix += ApiVersion

	fmt.Println(apiPrefix)

	publicGroup := v1Controller.Controller.Engine.Group(apiPrefix, cors.New(v1Controller.Controller.GetCorsConfig()))

	publicGroup.GET(controller.PingPath, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	v1Controller.AuthController.Init(publicGroup)

	return nil
}

func (v1Controller *V1Controller) Start() error {
	return nil
}
