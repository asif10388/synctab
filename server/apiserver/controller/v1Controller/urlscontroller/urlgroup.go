package urlscontroller

import (
	"net/http"

	controller "github.com/asif10388/synctab/apiserver/controller"
	"github.com/gin-gonic/gin"
)

func (urlsController *UrlController) addUrlGroupHandler(ctx *gin.Context) {
	var err error

	defer func() {
		if err != nil {

			ctx.Error(err)
			if err == controller.ErrUserNotAuthorized {
				ctx.JSON(http.StatusUnauthorized, controller.Response{Message: err.Error()})
			}

		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"Message": "OK",
			})
		}
	}()
}
