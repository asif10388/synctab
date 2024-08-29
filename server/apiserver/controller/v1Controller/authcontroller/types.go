package authcontroller

import (
	controller "github.com/asif10388/synctab/apiserver/controller"
	"github.com/asif10388/synctab/apiserver/model"
	"github.com/asif10388/synctab/apiserver/model/auth"
	env "github.com/asif10388/synctab/internal/environment"
)

const (
	authPrefix   = "/auth"
	loginPath    = "/login"
	logoutPath   = "/logout"
	registerPath = "/register"
)

type AuthController struct {
	Controller *controller.Controller
	Env        *env.Environment
	Model      *model.Model
	Auth       *auth.Auth
}
