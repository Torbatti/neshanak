package neshanak

import (
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/Torbatti/neshanak/core"
	"github.com/spf13/cobra"
)

var _ core.App = (*Neshanak)(nil)

// Version of Neshanak
var Version = "(untracked)"

type Neshanak struct {
	core.App

	devFlag           bool
	dataDirFlag       string
	encryptionEnvFlag string
	queryTimeout      int
	hideStartBanner   bool

	// RootCmd is the main console command
	RootCmd *cobra.Command
}

type Config struct {
	// hide the default console server info on app startup
	HideStartBanner bool

	// optional default values for the console flags
	DefaultDev           bool
	DefaultDataDir       string // if not set, it will fallback to "./pb_data"
	DefaultEncryptionEnv string
	DefaultQueryTimeout  time.Duration // default to core.DefaultQueryTimeout (in seconds)

}

func New() *Neshanak {
	_, isUsingGoRun := inspectRuntime()

	return NewWithConfig(Config{
		DefaultDev: isUsingGoRun,
	})

}

func NewWithConfig(config Config) *Neshanak {

	var nk *Neshanak
	var executableName string

	// initialize a default data directory based on the executable baseDir
	if config.DefaultDataDir == "" {
		baseDir, _ := inspectRuntime()
		config.DefaultDataDir = filepath.Join(baseDir, "nk_data")
	}

	executableName = filepath.Base(os.Args[0])

	nk = &Neshanak{
		RootCmd: &cobra.Command{
			Use:     executableName,
			Short:   executableName + " CLI",
			Version: Version,
			FParseErrWhitelist: cobra.FParseErrWhitelist{
				UnknownFlags: true,
			},
			// no need to provide the default cobra completion command
			CompletionOptions: cobra.CompletionOptions{
				DisableDefaultCmd: true,
			},
		},
	}

	nk.App = core.NewBaseApp(core.BaseAppConfig{
		IsDev:         nk.devFlag,
		DataDir:       nk.dataDirFlag,
		EncryptionEnv: nk.encryptionEnvFlag,
	})

	return nk
}

func (dst *Neshanak) Start() error {
	// dst.RootCmd.AddCommand(cmd.NewSuperuserCommand(dst))
	// dst.RootCmd.AddCommand(cmd.NewServeCommand(dst))

	return dst.Execute()
}

func (dst *Neshanak) Execute() error {
	// if !dst.skipBootstrap() {
	if err := dst.Bootstrap(); err != nil {
		return err
	}
	// }

	done := make(chan bool, 1)

	// listen for interrupt signal to gracefully shutdown the application
	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
		<-sigch

		done <- true
	}()

	// execute the root command
	go func() {
		// note: leave to the commands to decide whether to print their error
		dst.RootCmd.Execute()

		done <- true
	}()

	<-done

	return nil
}

func inspectRuntime() (baseDir string, withGoRun bool) {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		// probably ran with go run
		withGoRun = true
		baseDir, _ = os.Getwd()
	} else {
		// probably ran with go build
		withGoRun = false
		baseDir = filepath.Dir(os.Args[0])
	}
	return
}
