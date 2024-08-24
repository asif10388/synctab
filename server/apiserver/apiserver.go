package apiserver

import (
	"github.com/gin-gonic/gin"
)

func NewApiServer(input APIServerInput) *APIServer {
	apiServer := &APIServer{
		APIServerInput: APIServerInput{
			DeploymentType: input.DeploymentType,
		},
	}

	return apiServer
}

func (apiServer *APIServer) initEngine() error {
	engine := gin.Default()
	apiServer.Engine = engine

	return nil
}

func (apiServer *APIServer) Init() error {
	apiServer.initEngine()
	return nil
}

func (apiServer *APIServer) Start() {
	apiServer.Engine.Run(":5000")
}
