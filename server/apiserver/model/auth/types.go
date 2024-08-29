package auth

import "github.com/asif10388/synctab/apiserver/model"

type Auth struct {
	*model.Model
}

type Credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRegisterCredentials struct {
	Email string `json:"email" binding:"required"`
	*Credentials
}

type UserLoginCredentials struct {
	*Credentials
}

type AuthenticationResponse struct {
	Token  string
	UserId string
}
