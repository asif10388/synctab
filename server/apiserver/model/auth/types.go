package auth

import (
	"github.com/asif10388/synctab/apiserver/model"
	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	*model.Model
	UserId string
}

type Claims struct {
	Version int    `json:"version"`
	Email   string `json:"email"`
	jwt.RegisteredClaims
}

type Credentials struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRegisterCredentials struct {
	Username string `json:"username" binding:"required"`
	*Credentials
}

type UserLoginCredentials struct {
	*Credentials
}

type LoginResponse struct {
	Token    string
	Username string
	Email    string
}
