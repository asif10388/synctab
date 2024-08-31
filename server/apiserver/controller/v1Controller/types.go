package v1controller

import (
	controller "github.com/asif10388/synctab/apiserver/controller"
	"github.com/asif10388/synctab/apiserver/controller/v1Controller/authcontroller"
	"github.com/asif10388/synctab/apiserver/controller/v1Controller/urlscontroller"
	"github.com/asif10388/synctab/apiserver/middleware"
	"github.com/asif10388/synctab/apiserver/model"
	"github.com/asif10388/synctab/apiserver/model/auth"
	env "github.com/asif10388/synctab/internal/environment"
)

const ApiVersion = "v1"

type V1Controller struct {
	Auth           *auth.Auth
	Model          *model.Model
	Env            *env.Environment
	Middleware     *middleware.Middleware
	Controller     *controller.Controller
	AuthController *authcontroller.AuthController
	UrlsController *urlscontroller.UrlController
}
