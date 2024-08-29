package model

import (
	"errors"
	"regexp"
	"time"

	"github.com/asif10388/synctab/apiserver/controller"
	database "github.com/asif10388/synctab/apiserver/model/database"
	env "github.com/asif10388/synctab/internal/environment"
	"github.com/jackc/pgx/v5/pgtype"
)

type ModelDefaults struct {
	PrimarySchemaName        string
	PrimarySchemaNamePattern string
	PrimarySchemaNameRegexp  *regexp.Regexp

	SchemaSetupTimeout time.Duration

	MaxNameLen      int
	MaxEmailLen     int
	MaxUsernameLen  int
	MaxUrlLen       int
	MaxIPLen        int
	MaxPasswordLen  int
	MinPasswordLen  int
	MaxSchemaLen    int
	MaxHostnameLen  int
	MaxEndpointsLen int

	EndpointsSeparator string
	HostPortSeparator  string

	HostnamePattern string
	HostnameRegexp  *regexp.Regexp

	LoginIdPattern string
	LoginIdRegexp  *regexp.Regexp

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
	CreatedAt   pgtype.Timestamptz `json:"-"`
	CreatedAtTm *time.Time         `json:"created_at,omitempty"`
	UpdatedAt   pgtype.Timestamptz `json:"-"`
	UpdatedAtTm *time.Time         `json:"updated_at,omitempty"`
	DeletedAt   pgtype.Timestamptz `json:"-"`
	DeletedAtTm *time.Time         `json:"-"`
}

// Errors
var (
	// User
	ErrUserInvalid           = errors.New("invalid user")
	ErrUserNotFound          = errors.New("failed to find user")
	ErrUserAddFailed         = errors.New("failed to add user")
	ErrUserUpdateFailed      = errors.New("failed to update user")
	ErrUserDeleteFailed      = errors.New("failed to delete user")
	ErrUserInitFailed        = errors.New("failed to initialize user")
	ErrUserSetPasswordFailed = errors.New("failed to set user password")
	ErrUserAuthFailed        = errors.New("user authentication failure")

	// Token
	ErrInvalidToken = errors.New("invalid token")
	ErrUpdateToken  = errors.New("update token")
)
