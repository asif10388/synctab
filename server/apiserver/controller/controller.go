package controller

import (
	env "github.com/asif10388/synctab/internal/environment"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var controllerInstance *Controller

func NewController() (*Controller, error) {
	if controllerInstance == nil {
		controllerInstance = &Controller{
			env: env.GetEnvironment(),
		}
	}

	return controllerInstance, nil
}

func (controller *Controller) Init(engine *gin.Engine) error {
	controller.Engine = engine
	return nil
}

func (controller *Controller) GetApiVersions() map[string]bool {
	apiVersions := make(map[string]bool)

	for _, version := range controller.env.GetTrimmedStringListEnv("SYNCTAB_API_VERSIONS") {
		apiVersions[version] = true
	}

	return apiVersions
}

func (controller *Controller) GetCorsConfig() cors.Config {
	allowedOrigins := make(map[string]bool)
	allowAllOrigins := false

	for _, origin := range controller.env.GetTrimmedStringListEnv("SYNCTAB_API_ALLOW_ORIGINS") {
		if origin == "*" {
			allowAllOrigins = true
			break
		}

		allowedOrigins[origin] = true
	}

	corsConfig := cors.Config{
		AllowMethods:     controller.env.GetTrimmedStringListEnv("SYNCTAB_API_ALLOW_METHODS"),
		AllowHeaders:     controller.env.GetTrimmedStringListEnv("SYNCTAB_API_ALLOW_HEADERS"),
		AllowCredentials: controller.env.GetBoolEnv("SYNCTAB_API_ALLOW_CREDENTIALS"),
		AllowWebSockets:  controller.env.GetBoolEnv("SYNCTAB_API_ALLOW_WEBSOCKETS"),
	}

	if allowAllOrigins {
		corsConfig.AllowOriginFunc = func(origin string) bool {
			return true
		}
	} else {
		corsConfig.AllowOriginFunc = func(origin string) bool {
			if _, exists := allowedOrigins[origin]; exists {
				return true
			}

			return false
		}
	}

	return corsConfig
}
