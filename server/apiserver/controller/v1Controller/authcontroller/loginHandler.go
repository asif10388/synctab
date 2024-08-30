package authcontroller

import (
	"net/http"

	"github.com/asif10388/synctab/apiserver/controller"
	"github.com/asif10388/synctab/apiserver/model/auth"
	"github.com/gin-gonic/gin"
)

func (authController *AuthController) loginHandler(ctx *gin.Context) {
	var err error
	var response *auth.LoginResponse

	defer func() {
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusUnauthorized, controller.Response{Message: err.Error()})
		} else {
			ctx.JSON(http.StatusOK, LoginResponse{
				Email:    response.Email,
				Username: response.Username,
				Token:    response.Token,
			})
		}
	}()

	response, err = authController.Auth.LoginUser(ctx)
	if err != nil {
		return
	}
}
