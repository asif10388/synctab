package apiserver

import "github.com/gin-gonic/gin"

type APIServer struct {
	Engine *gin.Engine

	APIServerInput
}

type APIServerInput struct {
	DeploymentType string
}
