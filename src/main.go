package main

import (
	"log/slog"
	"os"
	"yaus/internal/config"
	"yaus/internal/logger/sl"
	"yaus/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting yaus", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	db, err := sqlite.New(cfg.DBPath)
	if err != nil {
		log.Error("failed to init data base", sl.Err(err))
		os.Exit(1)
	}

	err = db.SaveURL("https://google.com", "google")
	if err != nil {
		log.Error("failed to save URL", sl.Err(err))
	}

	err = db.SaveURL("https://google.com", "google")
	if err != nil {
		log.Error("failed to save URL", sl.Err(err))
	}

	_ = db
	// init router: chi github.com/go-chi/chi
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
