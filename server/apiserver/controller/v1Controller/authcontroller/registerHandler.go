package authcontroller

import (
	"net/http"

	"github.com/asif10388/synctab/apiserver/controller"
	"github.com/gin-gonic/gin"
)

func (authController *AuthController) registerHandler(ctx *gin.Context) {
	var err error
	// var res *auth.AuthenticationResponse

	// var requestBody any
	// jsonData := ctx.BindJSON(requestBody)

	defer func() {
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusUnauthorized, controller.Response{Message: "Authentication failed"})
		} else {
			ctx.JSON(http.StatusOK, controller.Response{Message: "Success my ass"})
		}
	}()

	err = authController.Auth.UserRegister(ctx)
	if err != nil {
		return
	}
}
