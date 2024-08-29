package auth

import "github.com/asif10388/synctab/apiserver/model"

var authInstance *Auth

func NewAuthModel(model *model.Model) (*Auth, error) {
	if authInstance == nil {
		authInstance = &Auth{
			Model: model,
		}
	}

	return authInstance, nil
}
