package core

import (
	"database/sql"
)

type App interface {
	// Logger() *slog.Logger

	IsBootstrapped() bool

	Bootstrap() error

	DataDir() string

	EncryptionEnv() string

	IsDev() bool

	Db() *sql.DB
}
