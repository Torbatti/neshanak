package neshanak

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/Torbatti/neshanak/cmd"
	"github.com/Torbatti/neshanak/core"
	"github.com/Torbatti/neshanak/utils"
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
	utils.AssertNotEmptyString(config.DefaultDataDir)

	executableName = filepath.Base(os.Args[0])
	utils.AssertNotEmptyString(executableName)

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
		devFlag:           config.DefaultDev,
		dataDirFlag:       config.DefaultDataDir,
		queryTimeout:      int(config.DefaultQueryTimeout),
		encryptionEnvFlag: config.DefaultEncryptionEnv,
		hideStartBanner:   config.HideStartBanner,
	}

	nk.App = core.NewBaseApp(core.BaseAppConfig{
		IsDev:         nk.devFlag,
		DataDir:       nk.dataDirFlag,
		EncryptionEnv: nk.encryptionEnvFlag,
	})

	return nk
}

func (nk *Neshanak) Start() error {
	nk.RootCmd.AddCommand(cmd.NewSuperuserCommand(nk))
	nk.RootCmd.AddCommand(cmd.NewServeCommand(nk, !nk.hideStartBanner))

	return nk.Execute()
}

func (nk *Neshanak) Execute() error {
	// if !nk.skipBootstrap() {
	if err := nk.Bootstrap(); err != nil {
		log.Fatal(err)
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
		nk.RootCmd.Execute()

		done <- true
	}()

	<-done

	return nil
}

// TODO:
func (nk *Neshanak) skipBootstrap() bool {
	// flags := []string{
	// 	"-h",
	// 	"--help",
	// 	"-v",
	// 	"--version",
	// }

	if nk.IsBootstrapped() {
		return true // already bootstrapped
	}

	// cmd, _, err := nk.RootCmd.Find(os.Args[1:])
	// if err != nil {
	// 	return true // unknown command
	// }

	// for _, arg := range os.Args {
	// 	if !list.ExistInSlice(arg, flags) {
	// 		continue
	// 	}

	// 	// ensure that there is no user defined flag with the same name/shorthand
	// 	trimmed := strings.TrimLeft(arg, "-")
	// 	if len(trimmed) > 1 && cmd.Flags().Lookup(trimmed) == nil {
	// 		return true
	// 	}
	// 	if len(trimmed) == 1 && cmd.Flags().ShorthandLookup(trimmed) == nil {
	// 		return true
	// 	}
	// }

	return false
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
