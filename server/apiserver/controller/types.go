package controller

import (
	"errors"

	env "github.com/asif10388/synctab/internal/environment"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	env    *env.Environment
	Engine *gin.Engine
}

type ApiVersion interface {
	Init() error
	Start() error
	// Stop() error
}

const (
	PingPath     = "/ping"
	ApiModelTime = "api-model-time"
)

type Response struct {
	Message string `json:"message"`
}

type ModelFunc func(*gin.Context) error

var (
	ErrInvalidName               = errors.New("invalid name")
	ErrInvalidInput              = errors.New("invalid input")
	ErrInvalidPassword           = errors.New("invalid password")
	ErrNoTokenFound              = errors.New("did not find token")
	ErrInvalidCreds              = errors.New("invalid credentials")
	ErrUserNotAuthorized         = errors.New("user not authorized")
	ErrInternal                  = errors.New("internal server error")
	ErrInvalidEmail              = errors.New("invalid email address")
	ErrIncorrectEmail            = errors.New("incorrect email address")
	ErrEnvironmentNotInitialized = errors.New("failed to initialize environment")
)

type RoutePrefixes struct {
	AuthPrefix      string
	UsersPrefix     string
	UserIdPrefix    string
	AccountsPrefix  string
	ResourcesPrefix string
	AccountIdPrefix string
}

type RoutePaths struct {
	PingPath           string
	LoginPath          string
	LogoutPath         string
	RegisterPath       string
	ValidatePath       string
	ChangePasswordPath string
}

type AuthType string

const BearerAuth AuthType = "bearer"
