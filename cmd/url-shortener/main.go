package main

import (
	"log/slog"
	"os"
	"url-shortener/internal/config"
	sl "url-shortener/internal/lib/logger/slog"
	"url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	// fmt.Println(cfg)

	log := setupLogger(cfg.Env)

	// log = log.With(slog.String("env", cfg.Env))
	log.Info("Starting server...", slog.String("env", cfg.Env))
	log.Debug("Debug started")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to create storage", sl.Err(err))
		os.Exit(1)
		// return
	}

	_ = storage

	// TODO: init router: chi, render

	// TODO: run server
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
