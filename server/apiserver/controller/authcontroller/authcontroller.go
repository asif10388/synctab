package authcontroller

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (ac *AuthController) loginHandler(ctx *gin.Context) {
	fmt.Println("Hello from auth")
}
