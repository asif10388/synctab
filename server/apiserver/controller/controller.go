package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DefaultResponse struct {
	Name string
	Age  float64
}

func MakeController() *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, DefaultResponse{
			Name: "BLEH BLEH",
			Age:  32,
		})
	})

	fmt.Println("HI from controller")

	return router
}
