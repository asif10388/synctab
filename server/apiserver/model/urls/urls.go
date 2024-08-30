package urls

import (
	"fmt"

	"github.com/asif10388/synctab/apiserver/model"
	"github.com/asif10388/synctab/apiserver/model/auth"
)

var urlsInstance *Urls

func NewUrlsModel(model *model.Model) (*Urls, error) {
	if urlsInstance == nil {
		auth, err := auth.NewAuthModel(model)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize url model with auth: %w", err)
		}

		urlsInstance = &Urls{
			Model: model,
			Auth:  auth,
		}
	}

	return urlsInstance, nil
}
