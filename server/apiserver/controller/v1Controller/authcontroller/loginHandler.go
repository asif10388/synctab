package authcontroller

import (
	"net/http"

	"github.com/asif10388/synctab/apiserver/controller"
	"github.com/gin-gonic/gin"
)

func (authController *AuthController) loginHandler(ctx *gin.Context) {

	defer func() {
		ctx.JSON(http.StatusOK, controller.Response{Message: "Success"})
	}()
}
