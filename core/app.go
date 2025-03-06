package core

import (
	"database/sql"
	"log/slog"
)

type App interface {
	Logger() *slog.Logger

	IsBootstrapped() bool

	Bootstrap() error

	DataDir() string

	EncryptionEnv() string

	IsDev() bool

	DB() *sql.DB
}
