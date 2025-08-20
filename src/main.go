package main

import (
	"log/slog"
	"os"
	"yaus/internal/config"
	mwLogger "yaus/internal/http-server/middleware/mvLogger"
	"yaus/internal/logger/sl"
	"yaus/internal/storage/sqlite"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	// logger
	log := setupLogger(cfg.Env)
	log.Info("starting yaus", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// sqlite setup
	db, err := sqlite.New(cfg.DBPath)
	if err != nil {
		log.Error("failed to init data base", sl.Err(err))
		os.Exit(1)
	}

	_ = db

	// router
	router := chi.NewRouter()

	router.Use(middleware.RequestID) // to better request handling
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
