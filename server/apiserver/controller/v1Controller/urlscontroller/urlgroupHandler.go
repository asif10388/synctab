package urlscontroller

import (
	"net/http"

	controller "github.com/asif10388/synctab/apiserver/controller"
	"github.com/asif10388/synctab/apiserver/model/urls"
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
				"message": "successfully created urls",
			})
		}
	}()

	_, err = urlsController.Urls.CreateUrls(ctx)
	if err != nil {
		return
	}

}

func (urlsController *UrlController) getUrlsByUserHandler(ctx *gin.Context) {
	var err error
	var response *[]urls.TransformUrls

	defer func() {
		if err != nil {
			ctx.Error(err)
			if err == controller.ErrUserNotAuthorized {
				ctx.JSON(http.StatusUnauthorized, controller.Response{Message: err.Error()})
			} else {
				ctx.JSON(http.StatusInternalServerError, controller.Response{Message: err.Error()})
			}

		} else {
			if response == nil {
				ctx.JSON(http.StatusOK, gin.H{
					"message": "No urls were found",
				})
			} else {
				ctx.JSON(http.StatusOK, response)
			}

		}
	}()

	response, err = urlsController.Urls.GetUrlsByUserId(ctx)
	if err != nil {
		return
	}

}

func (urlsController *UrlController) deleteUrlByIdHandler(ctx *gin.Context) {
	var err error

	defer func() {
		if err != nil {
			ctx.Error(err)
			if err == controller.ErrUserNotAuthorized {
				ctx.JSON(http.StatusUnauthorized, controller.Response{Message: err.Error()})
			} else {
				ctx.JSON(http.StatusInternalServerError, controller.Response{Message: err.Error()})
			}

		} else {

			ctx.JSON(http.StatusOK, gin.H{
				"message": "URL Deleted successfully",
			})
		}
	}()

	err = urlsController.Urls.DeleteUrlById(ctx)
	if err != nil {
		return
	}
}
