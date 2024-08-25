package apiserver

import (
	controller "github.com/asif10388/synctab/apiserver/controller"
	v1controller "github.com/asif10388/synctab/apiserver/controller/v1Controller"
	env "github.com/asif10388/synctab/internal/environment"
	"github.com/gin-gonic/gin"
)

type ApiVersionAllocator func(*controller.Controller) (controller.ApiVersion, error)

var ApiVersionAllocators = map[string]ApiVersionAllocator{
	v1controller.ApiVersion: v1controller.NewV1Controller,
}

type APIServer struct {
	Engine     *gin.Engine
	Env        *env.Environment
	Controller *controller.Controller
	Versions   map[string]controller.ApiVersion

	APIServerInput
}

type APIServerInput struct {
	DeploymentType string
}
