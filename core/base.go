package core

import (
	"database/sql"
)

var _ App = (*BaseApp)(nil)

// BaseAppConfig defines a BaseApp configuration option
type BaseAppConfig struct {
	DataDir       string
	EncryptionEnv string
	IsDev         bool

	JWT_SECRET string
}

type BaseApp struct {
	config *BaseAppConfig

	db *sql.DB
	// jwtauth *jwtauth.JWTAuth

	// TODO:
	// logger *slog.Logger
}
