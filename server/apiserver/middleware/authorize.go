package middleware

import (
	"net/http"
	"strings"

	"github.com/asif10388/synctab/apiserver/controller"
	"github.com/asif10388/synctab/apiserver/model/auth"
	"github.com/gin-gonic/gin"
)

type AutorizationInfo struct {
	updated bool
	token   string
	claims  *auth.Claims
}

func (middleware *Middleware) getTokenFromHeader(ctx *gin.Context) (string, error) {
	authorizationHeader := ctx.Request.Header.Get("Authorization")
	if len(authorizationHeader) == 0 {
		return "", controller.ErrNoTokenFound
	}

	authorizationHeaderParts := strings.Split(authorizationHeader, " ")
	if len(authorizationHeaderParts) != 2 {
		return "", controller.ErrNoTokenFound
	}

	if !strings.EqualFold(authorizationHeaderParts[0], string(controller.BearerAuth)) {
		return "", controller.ErrNoTokenFound
	}

	if authorizationHeaderParts[1] != "" {
		return authorizationHeaderParts[1], nil
	} else {
		return "", controller.ErrInternal
	}
}

func (middleware *Middleware) getAuthorizationInfo(ctx *gin.Context) (*AutorizationInfo, error) {
	currentToken, err := middleware.getTokenFromHeader(ctx)
	if err != nil {
		return nil, err
	}

	authorizationInfo := &AutorizationInfo{
		updated: false,
		token:   currentToken,
	}

	currentClaims, err := middleware.Auth.ValidateJWTToken(ctx, currentToken)
	if err != nil {
		return nil, err
	}

	authorizationInfo.claims = currentClaims
	return authorizationInfo, nil
}

func (middleware *Middleware) authorize(ctx *gin.Context) {
	var err error
	status := http.StatusOK

	defer func() {
		if status >= http.StatusBadRequest {
			ctx.AbortWithStatus(status)
		}
	}()

	_, err = middleware.getAuthorizationInfo(ctx)
	if err != nil {
		status = http.StatusUnauthorized
		return
	}

	ctx.Next()

}

func (middleware *Middleware) UserAuthorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		middleware.authorize(ctx)
	}
}
