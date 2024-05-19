package handlers

import (
	"log/slog"
	"net/http"

	"github.com/TOsmanov/qr-gen/internal/lib/api/response"
)

type Response struct {
	response.Response
	Message any `json:"data,omitempty"`
}

func IndexHandler(log *slog.Logger) http.HandlerFunc {
	slog.Info("IndexHandler")
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.IndexHandler"
		slog.Info(op)
		log = log.With(
			slog.String("op", op),
		)
		Index(log, w, r)
	}
}

func PreviewHandler(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.PreviewHandler"
		slog.Info(op)
		log = log.With(
			slog.String("op", op),
		)
		GetPreview(log, w, r)
	}
}

func GenerationHandler(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GenerationHandler"
		slog.Info(op)
		log = log.With(
			slog.String("op", op),
		)
		GetDownlaodLink(log, w, r)
	}
}
