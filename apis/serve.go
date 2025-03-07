package apis

import (
	"log"
	"net/http"
	"time"

	"github.com/Torbatti/neshanak/core"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

// ServeConfig defines a configuration struct for apis.Serve().
type ServeConfig struct {
	// ShowStartBanner indicates whether to show or hide the server start console message.
	ShowStartBanner bool

	// HttpAddr is the TCP address to listen for the HTTP server (eg. "127.0.0.1:80").
	HttpAddr string

	// HttpsAddr is the TCP address to listen for the HTTPS server (eg. "127.0.0.1:443").
	HttpsAddr string

	// Optional domains list to use when issuing the TLS certificate.
	//
	// If not set, the host from the bound server address will be used.
	//
	// For convenience, for each "non-www" domain a "www" entry and
	// redirect will be automatically added.
	CertificateDomains []string

	// AllowedOrigins is an optional list of CORS origins (default to "*").
	AllowedOrigins []string
}

func NewRouter(app core.App) (*chi.Mux, error) {

	var router *chi.Mux

	router = chi.NewRouter()

	// Default Middlewares
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{},
		ExposedHeaders:   []string{},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	router.Use(middleware.Compress(5, "text/html", "text/css", "text/javascript"))
	router.Use(httprate.LimitByIP(100, 1*time.Minute))

	return router, nil

}

func Serve(app core.App, config ServeConfig) error {
	var server http.Server
	var router *chi.Mux
	var mainAddr string

	var err error

	if len(config.AllowedOrigins) == 0 {
		config.AllowedOrigins = []string{"*"}
	}

	// ensure that the latest migrations are applied before starting the server
	// err := app.RunAllMigrations()
	// if err != nil {
	// 	return err
	// }

	router, err = NewRouter(app)
	if err != nil {
		log.Fatal(err)
		// return err
	}

	mainAddr = config.HttpAddr
	if config.HttpsAddr != "" {
		mainAddr = config.HttpsAddr
	}

	server = http.Server{
		Addr:    mainAddr,
		Handler: router,
	}
	log.Println("Starting server on :", server.Addr)

	if err = server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	// TODO:  graceful shutdown

	return nil
}
