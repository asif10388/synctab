package model

import (
	"errors"
	"time"

	"github.com/asif10388/synctab/apiserver/controller"
	database "github.com/asif10388/synctab/apiserver/model/database"
	env "github.com/asif10388/synctab/internal/environment"
	"github.com/jackc/pgx/v5/pgtype"
)

type ModelDefaults struct {
	SchemaName  string
	TokenSecret string

	MaxUrlLen      int
	MaxNameLen     int
	MaxEmailLen    int
	MaxUsernameLen int
	MaxPasswordLen int
	MinPasswordLen int

	MaxPaginationEntries int

	UserTimeout     time.Duration
	UserTimeoutSecs int64

	TokenExpiresDuration time.Duration
	TokenUpdateDuration  time.Duration

	TokenExpiresDurationSecs int64
	TokenUpdateDurationSecs  int64
}

type Model struct {
	Env        *env.Environment
	Database   *database.Database
	Controller *controller.Controller

	ModelDefaults
}

type ModelCommon struct {
	CreatedAt pgtype.Timestamptz `json:"-"`
	UpdatedAt pgtype.Timestamptz `json:"-"`
}

// Errors
var (
	// User
	ErrUserInvalid      = errors.New("invalid user")
	ErrUserAddFailed    = errors.New("failed to add user")
	ErrUserNotFound     = errors.New("failed to find user")
	ErrUserUpdateFailed = errors.New("failed to update user")
	ErrUserDeleteFailed = errors.New("failed to delete user")
	ErrUserAuthFailed   = errors.New("failed to authenticate user")

	// Token
	ErrInvalidToken = errors.New("invalid token")
	ErrUpdateToken  = errors.New("update token")
)
