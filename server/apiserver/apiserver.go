package apiserver

import (
	"fmt"

	"github.com/rs/zerolog/log"

	controller "github.com/asif10388/synctab/apiserver/controller"
	env "github.com/asif10388/synctab/internal/environment"
	"github.com/gin-gonic/gin"
)

func NewApiServer(input APIServerInput) (*APIServer, error) {
	apiServer := &APIServer{
		APIServerInput: APIServerInput{
			DeploymentType: input.DeploymentType,
		},
	}

	newEnv, err := env.NewEnvironment(apiServer.DeploymentType)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize environment: %w", err)
	}
	apiServer.Env = newEnv

	newController, err := controller.NewController()
	if err != nil {
		fmt.Println("failed to allocate controller")
		return nil, fmt.Errorf("failed to allocate controller: %w", err)
	}
	apiServer.Controller = newController

	configuredApiVersions := newController.GetApiVersions()
	apiServer.Versions = make(map[string]controller.ApiVersion)

	for apiVersion, apiAllocatorFunc := range ApiVersionAllocators {
		supported, exists := configuredApiVersions[apiVersion]
		if !supported || !exists {
			fmt.Println("ignoring unsupported API version")
			continue
		}

		allocatedServer, err := apiAllocatorFunc(apiServer.Controller)
		if err != nil {
			fmt.Println("failed to allocate API version")
			return nil, err
		}

		apiServer.Versions[apiVersion] = allocatedServer
	}

	return apiServer, nil
}

func (apiServer *APIServer) initEngine() error {
	engine := gin.Default()
	apiServer.Engine = engine

	return nil
}

func (apiServer *APIServer) Init() error {

	apiServer.initEngine()

	err := apiServer.Controller.Init(apiServer.Engine)
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize controller")
	}

	for apiVersion, apiController := range apiServer.Versions {
		err = apiController.Init()
		if err != nil {
			log.Error().Err(err).Msgf("failed to initialize API version %s", apiVersion)
			return fmt.Errorf("failed to initialize API version %s: %v", apiVersion, err)
		}
	}

	return nil
}

func (apiServer *APIServer) Start() {

	for apiVersion, apiController := range apiServer.Versions {
		err := apiController.Start()
		if err != nil {
			log.Error().Err(err).Msgf("failed to start API version %s", apiVersion)
		}

		log.Info().Msgf("starting api version %s", apiVersion)
	}

	port := apiServer.Env.GetStrEnv("SYNCTAB_API_PORTS")

	apiServer.Engine.Run(":" + port)
}
