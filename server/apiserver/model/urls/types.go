package urls

import (
	"github.com/asif10388/synctab/apiserver/model"
	"github.com/asif10388/synctab/apiserver/model/auth"
)

type Urls struct {
	*auth.Auth
	*model.Model
}
