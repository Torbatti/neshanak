package core

import (
	"database/sql"

	"github.com/go-chi/jwtauth/v5"
)

type App interface {
	// Logger() *slog.Logger

	IsBootstrapped() bool

	Bootstrap() error

	DataDir() string

	EncryptionEnv() string

	IsDev() bool

	Db() *sql.DB

	JwtAuth() *jwtauth.JWTAuth
}
