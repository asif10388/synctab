package authcontroller

import (
	"net/http"

	"github.com/asif10388/synctab/apiserver/controller"
	"github.com/gin-gonic/gin"
)

func (authController *AuthController) registerHandler(ctx *gin.Context) {
	var err error

	defer func() {
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusUnauthorized, controller.Response{Message: err.Error()})
		} else {
			ctx.JSON(http.StatusOK, controller.Response{Message: "Successfully created user"})
		}
	}()

	err = authController.Auth.CreateUser(ctx)
	if err != nil {
		return
	}
}
