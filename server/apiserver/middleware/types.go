package middleware

import (
	"github.com/asif10388/synctab/apiserver/controller"
	"github.com/asif10388/synctab/apiserver/model"
	"github.com/asif10388/synctab/apiserver/model/auth"
	env "github.com/asif10388/synctab/internal/environment"
)

type Middleware struct {
	Auth       *auth.Auth
	Model      *model.Model
	Env        *env.Environment
	Controller *controller.Controller
}
