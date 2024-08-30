package authcontroller

import (
	controller "github.com/asif10388/synctab/apiserver/controller"
	model "github.com/asif10388/synctab/apiserver/model"
	auth "github.com/asif10388/synctab/apiserver/model/auth"
	env "github.com/asif10388/synctab/internal/environment"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var authControllerInstance *AuthController

func NewAuthController(controller *controller.Controller, model *model.Model) (*AuthController, error) {
	if authControllerInstance == nil {
		authModel, err := auth.NewAuthModel(model)
		if err != nil {
			log.Error().Err(err).Msg("failed to initialize v1 auth model")
			return nil, err
		}

		authControllerInstance = &AuthController{
			Controller: controller,
			Model:      model,
			Auth:       authModel,
			Env:        env.GetEnvironment(),
		}
	}

	return authControllerInstance, nil
}

func (authController *AuthController) Init(publicGroup *gin.RouterGroup) {
	authGroup := publicGroup.Group(authPrefix)

	authGroup.POST(loginPath, authController.loginHandler)
	authGroup.POST(registerPath, authController.registerHandler)
}
