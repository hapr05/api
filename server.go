package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	// cfg is the global configuration for the server. It's read in at startup from
	// the config.json file and enviornment variables, see config.go for more info.
	cfg *config
	// log output
	logger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
	// application database connection
	appDB *sql.DB
)

// NewServerRoutes returns a Muxer that has all API routes.
// This makes for easy testing using httptest, see server_test.go
func NewServerRoutes() *http.ServeMux {
	m := http.NewServeMux()

	m.HandleFunc("/", NotFoundHandler)
	m.Handle("/status", middleware(HealthCheckHandler))
	m.HandleFunc("/.well-known/acme-challenge/", CertbotHandler)

	// m.Handle("/v0/users", middleware(UserHandler))
	// m.Handle("/v0/users/", middleware(UsersHandler))

	m.Handle("/v0/primers", middleware(PrimersHandler))
	m.Handle("/v0/primers/", middleware(PrimerHandler))

	m.Handle("/v0/sources", middleware(SourcesHandler))
	m.Handle("/v0/sources/", middleware(SourceHandler))

	m.Handle("/v0/urls", middleware(UrlsHandler))
	m.Handle("/v0/urls/", middleware(UrlHandler))

	m.Handle("/v0/uncrawlables", middleware(UncrawlablesHandler))
	m.Handle("/v0/uncrawlables/", middleware(UncrawlableHandler))

	// m.Handle("/v0/links", middleware(UrlHandler))
	// m.Handle("/v0/links/", middleware(UrlsHandler))

	// m.Handle("/v0/snapshots", middleware())
	// m.Handle("/v0/snapshots/", middleware())

	// m.Handle("/v0/content", middleware())
	// m.Handle("/v0/content/", middleware())

	// m.Handle("/v0/metadata", middleware())
	// m.Handle("/v0/metadata/", middleware())

	// m.Handle("/v0/consensus", middleware())
	// m.Handle("/v0/consensus/", middleware())

	// m.Handle("/v0/collections", middleware())
	// m.Handle("/v0/collections/", middleware())

	return m
}

// main app entry point
func main() {
	var err error
	cfg, err = initConfig(os.Getenv("GOLANG_ENV"))
	if err != nil {
		// panic if the server is missing a vital configuration detail
		panic(fmt.Errorf("server configuration error: %s", err.Error()))
	}

	// TODO - run this in a goroutine & support reporting "oh snap no DB"
	// while waiting for a connection
	connectToAppDb()

	// base server
	s := &http.Server{}
	// connect mux routes to server
	s.Handler = NewServerRoutes()

	// print notable config settings
	printConfigInfo()

	// fire it up!
	fmt.Println("starting server on port", cfg.Port)

	// start server wrapped in a log.Fatal b/c http.ListenAndServe will not
	// return unless there's an error
	logger.Fatal(StartServer(cfg, s))
}
