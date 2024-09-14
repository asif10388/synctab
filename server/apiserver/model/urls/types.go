package urls

import (
	"time"

	"github.com/asif10388/synctab/apiserver/model"
	"github.com/asif10388/synctab/apiserver/model/auth"
)

type UrlModel struct {
	Id      string `json:"id"`
	Url     string `json:"url"`
	Title   string `json:"title"`
	GroupId string `json:"group_id"`
	UserId  string `json:"user_id,omitempty"`
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
	Tabs      []Tabs    `json:"tabs"`
	GroupId   string    `json:"group_id"`
	CreatedAt time.Time `json:"created_at"`
}
