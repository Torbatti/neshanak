package core

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/Torbatti/neshanak/utils"
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

	JWT_SECRET string
}

type BaseApp struct {
	config *BaseAppConfig

	db *sql.DB
	// jwtauth *jwtauth.JWTAuth

	// TODO:
	// logger *slog.Logger
}

func (app *BaseApp) Bootstrap() error {

	utils.AssertNotEmptyString(app.DataDir())

	// TODO:
	// ensure that data dir exist
	if err := os.MkdirAll(app.DataDir(), os.ModePerm); err != nil {
		// log.Fatal(err)
		return err
	}
	if err := os.MkdirAll(app.DataDir()+"/img", os.ModePerm); err != nil {
		// log.Fatal(err)
		return err
	}

	// Initiate Sqlite
	dbPath := filepath.Join(app.DataDir(), "data.db?_journal=WAL")
	utils.AssertNotEmptyString(dbPath)

	// if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
	// 	panic(err)
	// }

	db, err := sql.Open("sqlite", dbPath)
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
