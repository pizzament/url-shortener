package main

import (
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/url/save"
	mwLogger "url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/lib/logger/handlers/slogpretty"
	sl "url-shortener/internal/lib/logger/slog"
	"url-shortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	// log.Info("Starting server...", slog.String("env", cfg.Env))
	// log.Debug("Debug started")
	// log.Error("ErrorErrorError")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to create storage", sl.Err(err))
		os.Exit(1)
		// return
	}

	_ = storage

	// TODO: init router: chi, render
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))

	log.Info("Starting server...", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServerConfig.TimeOut,
		WriteTimeout: cfg.HTTPServerConfig.TimeOut,
		IdleTimeout:  cfg.HTTPServerConfig.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start server", sl.Err(err))
		os.Exit(1)
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		// log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
