package urlscontroller

import (
	"net/http"

	controller "github.com/asif10388/synctab/apiserver/controller"
	"github.com/asif10388/synctab/apiserver/middleware"
	model "github.com/asif10388/synctab/apiserver/model"
	urls "github.com/asif10388/synctab/apiserver/model/urls"
	env "github.com/asif10388/synctab/internal/environment"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var urlControllerInstance *UrlController

func NewUrlsController(controller *controller.Controller, model *model.Model, middleware *middleware.Middleware) (*UrlController, error) {
	if urlControllerInstance == nil {
		urlsModel, err := urls.NewUrlsModel(model)
		if err != nil {
			log.Error().Err(err).Msg("failed to initialize v1 urls model")
			return nil, err
		}

		urlControllerInstance = &UrlController{
			Model:      model,
			Controller: controller,
			Urls:       urlsModel,
			Middleware: middleware,
			Env:        env.GetEnvironment(),
		}
	}

	return urlControllerInstance, nil
}

func (urlController *UrlController) Init(publicGroup *gin.RouterGroup) {
	urlsGroup := publicGroup.Group(urlsPrefix, urlController.Middleware.UserAuthorize())

	urlsGroup.GET(validatePath, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "valid token",
		})
	})

	urlsGroup.POST(urlGroupPath, urlController.addUrlGroupHandler)
	urlsGroup.GET(urlGroupPath, urlController.getUrlsByUserHandler)
	urlsGroup.DELETE(urlIdPrefix, urlController.deleteUrlByIdHandler)
}
