package v1controller

import (
	controller "github.com/asif10388/synctab/apiserver/controller"
	"github.com/asif10388/synctab/apiserver/controller/v1Controller/authcontroller"
	env "github.com/asif10388/synctab/internal/environment"
)

const ApiVersion = "v1"

type V1Controller struct {
	Env            *env.Environment
	Controller     *controller.Controller
	AuthController *authcontroller.AuthController
}
