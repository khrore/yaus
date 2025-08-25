package main

import (
	"log/slog"
	"net/http"
	"os"
	"yaus/internal/config"
	mwLogger "yaus/internal/http-server/middleware/logger"
	"yaus/internal/http-server/save"
	"yaus/internal/lib/logger/slogext"
	"yaus/internal/storage/sqlite"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.MustLoad()

	// logger
	log := slogext.SetupLogger(cfg.Env)
	log.Info("starting yaus", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// sqlite setup
	db, err := sqlite.New(cfg.DBPath)
	if err != nil {
		log.Error("failed to init data base", slogext.Err(err))
		os.Exit(1)
	}

	// router
	router := chi.NewRouter()

	router.Use(middleware.RequestID) // to better request handling
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// start server
	router.Post("/url", save.New(log, db))

	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}
