package authcontroller

import (
	controller "github.com/asif10388/synctab/apiserver/controller"
	env "github.com/asif10388/synctab/internal/environment"
	"github.com/gin-gonic/gin"
)

var authControllerInstance *AuthController

func NewAuthController(controller *controller.Controller) (*AuthController, error) {
	if authControllerInstance == nil {
		authControllerInstance = &AuthController{
			Controller: controller,
			Env:        env.GetEnvironment(),
		}
	}

	return authControllerInstance, nil
}

func (authController *AuthController) Init(publicGroup *gin.RouterGroup) {
	authGroup := publicGroup.Group(authPrefix)

	authGroup.POST(loginPath, authController.loginHandler)
}
