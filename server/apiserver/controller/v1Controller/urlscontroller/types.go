package urlscontroller

import (
	controller "github.com/asif10388/synctab/apiserver/controller"
	"github.com/asif10388/synctab/apiserver/middleware"
	"github.com/asif10388/synctab/apiserver/model"
	"github.com/asif10388/synctab/apiserver/model/urls"
	env "github.com/asif10388/synctab/internal/environment"
)

const (
	urlsPrefix   = "/urls"
	urlGroupPath = "/url-group"
	urlIdPrefix  = urlGroupPath + "/:id"
)

type UrlController struct {
	Controller *controller.Controller
	Middleware *middleware.Middleware
	Env        *env.Environment
	Model      *model.Model
	Urls       *urls.Urls
}
