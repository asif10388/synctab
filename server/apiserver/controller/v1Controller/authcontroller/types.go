package authcontroller

import (
	controller "github.com/asif10388/synctab/apiserver/controller"
	env "github.com/asif10388/synctab/internal/environment"
)

const (
	authPrefix = "/auth"
	loginPath  = "/login"
	logoutPath = "/logout"
)

type AuthController struct {
	Controller *controller.Controller
	Env        *env.Environment
}
