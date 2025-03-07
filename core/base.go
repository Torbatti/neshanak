package core

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/Torbatti/neshanak/utils"
	"github.com/go-chi/jwtauth/v5"
	_ "modernc.org/sqlite"
)

const (
	LocalStorageDirName       string = "storage"
	LocalBackupsDirName       string = "backups"
	LocalTempDirName          string = ".df_temp_to_delete" // temp df_data sub directory that will be deleted on each app.Bootstrap()
	LocalAutocertCacheDirName string = ".autocert_cache"
)

var _ App = (*BaseApp)(nil)

// BaseAppConfig defines a BaseApp configuration option
type BaseAppConfig struct {
	DataDir       string
	EncryptionEnv string
	IsDev         bool

	JwtAuth *jwtauth.JWTAuth
}

type BaseApp struct {
	config *BaseAppConfig

	db *sql.DB
	// jwtauth *jwtauth.JWTAuth

	// TODO:
	// logger *slog.Logger
}

func (app *BaseApp) Bootstrap() error {

	var jwt_secret string
	var jwt_auth *jwtauth.JWTAuth

	var db *sql.DB
	var dbPath string

	var err error

	// TODO:
	// ensure that data dir exist
	utils.AssertNotEmptyString(app.DataDir())
	if err = os.MkdirAll(app.DataDir(), os.ModePerm); err != nil {
		// log.Fatal(err)
		return err
	}
	if err = os.MkdirAll(app.DataDir()+"/img", os.ModePerm); err != nil {
		// log.Fatal(err)
		return err
	}

	// Get .env
	jwt_secret = os.Getenv("JWT_SECRET")
	utils.AssertNotEmptyString(jwt_secret)

	// Set Jwt Auth
	jwt_auth = jwtauth.New("HS256", []byte(jwt_secret), nil)
	app.config.JwtAuth = jwt_auth

	// Initiate Sqlite
	dbPath = filepath.Join(app.DataDir(), "data.db?_journal=WAL")
	utils.AssertNotEmptyString(dbPath)

	// if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
	// 	panic(err)
	// }

	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		panic("failed making connection with database!")
		// return err
	}
	app.db = db

	// try to cleanup the df_data temp directory (if any)
	_ = os.RemoveAll(filepath.Join(app.DataDir(), LocalTempDirName))

	return nil
}

func NewBaseApp(config BaseAppConfig) *BaseApp {
	app := &BaseApp{
		config: &config,
	}

	return app
}

func (app *BaseApp) JwtAuth() *jwtauth.JWTAuth {
	return app.config.JwtAuth
}

func (app *BaseApp) IsBootstrapped() bool {
	return app.db != nil
}

func (app *BaseApp) Db() *sql.DB {
	return app.db
}

func (app *BaseApp) DataDir() string {
	return app.config.DataDir
}

func (app *BaseApp) EncryptionEnv() string {
	return app.config.EncryptionEnv
}

func (app *BaseApp) IsDev() bool {
	return app.config.IsDev
}
