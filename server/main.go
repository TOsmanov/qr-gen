package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/TOsmanov/qr-gen/internal/config"
	"github.com/TOsmanov/qr-gen/internal/http-server/handlers"
	mwLogger "github.com/TOsmanov/qr-gen/internal/http-server/middleware/logger"
	"github.com/TOsmanov/qr-gen/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	var err error
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Starting qr-generation service", slog.String("env", cfg.Env))
	log.Debug("DEBUG messages are enabled", slog.String("env", cfg.Env))

	router := chi.NewRouter()

	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.HandleFunc("/", handlers.IndexHandler(log))
	router.HandleFunc("/preview", handlers.PreviewHandler(log))
	router.HandleFunc("/qrgen", handlers.GenerationHandler(log))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("Starting server", slog.String("address", srv.Addr))

	go func() {
		if err = srv.ListenAndServe(); err != nil {
			log.Error("Failed to serve server", sl.Err(err))
		}
	}()

	log.Info("Server started")

	<-done
	log.Info("Stopping server")

	Clean()
	log.Info("Server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case "local":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func Clean() {

}
