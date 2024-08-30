package urlscontroller

import (
	controller "github.com/asif10388/synctab/apiserver/controller"
	model "github.com/asif10388/synctab/apiserver/model"
	urls "github.com/asif10388/synctab/apiserver/model/urls"
	env "github.com/asif10388/synctab/internal/environment"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var urlControllerInstance *UrlController

func NewUrlsController(controller *controller.Controller, model *model.Model) (*UrlController, error) {
	if urlControllerInstance == nil {
		urlsModel, err := urls.NewUrlsModel(model)
		if err != nil {
			log.Error().Err(err).Msg("failed to initialize v1 urls model")
			return nil, err
		}

		urlControllerInstance = &UrlController{
			Controller: controller,
			Model:      model,
			Urls:       urlsModel,
			Env:        env.GetEnvironment(),
		}
	}

	return urlControllerInstance, nil
}

func (urlController *UrlController) Init(publicGroup *gin.RouterGroup) {
	urlsGroup := publicGroup.Group(urlsPrefix, urlController.Middleware.UserAuthorize())
	urlsGroup.POST(urlGroupPath, urlController.addUrlGroupHandler)
}
