// Package classification of Turnkey API
//
// Documentation for Turnkey API
//
//	Schemes: http,https
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/turnkeyca/monolith/auth"
	"github.com/turnkeyca/monolith/bitly"
	"github.com/turnkeyca/monolith/db"
	"github.com/turnkeyca/monolith/employment"
	"github.com/turnkeyca/monolith/pet"
	"github.com/turnkeyca/monolith/reference"
	"github.com/turnkeyca/monolith/roommate"
	"github.com/turnkeyca/monolith/server"
	"github.com/turnkeyca/monolith/shorturl"
	"github.com/turnkeyca/monolith/user"
)

func configureDocRoutes(router *mux.Router) {
	getRouter := router.Methods(http.MethodGet).Subrouter()
	opts := middleware.RedocOpts{SpecURL: "./swagger.yml"}
	getRouter.Handle("/docs", middleware.Redoc(opts, nil))
	getRouter.Handle("/swagger.yml", http.FileServer(http.Dir("./")))
}

func configureRoutes(logger *log.Logger) (*mux.Router, error) {
	router := mux.NewRouter()
	bitly := bitly.NewClient(logger)
	database, err := db.New(logger)
	if err != nil {
		logger.Fatalf("failed to create database %#v\n", err)
	}
	authenticator := auth.New(logger, database)

	configureDocRoutes(router)
	shorturl.ConfigureShortUrlRoutes(router, logger, bitly, authenticator)
	user.ConfigureUserRoutes(router, logger, database, authenticator)
	roommate.ConfigureRoommateRoutes(router, logger, database, authenticator)
	reference.ConfigureReferenceRoutes(router, logger, database, authenticator)
	pet.ConfigurePetRoutes(router, logger, database, authenticator)
	employment.ConfigureEmploymentRoutes(router, logger, database, authenticator)
	auth.ConfigureAuthRoutes(router, logger, database)

	return router, nil
}

func serve(logger *log.Logger, s *http.Server) {
	err := s.ListenAndServe()
	if err != nil {
		logger.Fatalf("failed to start %#v\n", err)
	}
}

func shutdown(logger *log.Logger, httpServer *http.Server, sc chan os.Signal) {
	signal.Notify(sc, os.Interrupt)
	s := <-sc
	logger.Printf("termination signal received - trying to shutdown gracefully: %#v\n", s)
	c, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	cancelFunc()
	httpServer.Shutdown(c)
}

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	err := godotenv.Load(".env")
	if err != nil {
		logger.Printf("failed to load environment from .env: %#v\n", err)
	}
	sm, err := configureRoutes(logger)
	if err != nil {
		logger.Fatalf("failed to configure routes: %#v\n", err)
	}
	logger.Println("starting server")
	srv := server.New(logger)
	httpServer := srv.NewHttpServer(sm)
	go serve(logger, httpServer)
	sc := make(chan os.Signal, 1)
	shutdown(logger, httpServer, sc)
}
