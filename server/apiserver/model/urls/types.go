package urls

import (
	"github.com/asif10388/synctab/apiserver/model"
	"github.com/asif10388/synctab/apiserver/model/auth"
)

type UrlModel struct {
	Id      string `json:"id"`
	Url     string `json:"url"`
	Title   string `json:"title"`
	UserId  string `json:"user_id,omitempty"`
	GroupId string `json:"group_id"`
	model.ModelCommon
}

type Urls struct {
	*auth.Auth
	*model.Model
	UrlResponse []UrlModel
	UrlRequest  []UrlModel
}

type Tabs struct {
	Id    string `json:"id"`
	Url   string `json:"url"`
	Title string `json:"title"`
}

type TransformUrls struct {
	GroupId string `json:"group_id"`
	Tabs    []Tabs `json:"tabs"`
}
