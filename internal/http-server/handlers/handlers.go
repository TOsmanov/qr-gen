package handlers

import (
	"log/slog"
	"net/http"

	"github.com/TOsmanov/qr-gen/internal/lib/api/response"
)

type Response struct {
	response.Response
	Body any `json:"data,omitempty"`
}

func IndexHandler(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.IndexHandler"
		log = log.With(
			slog.String("op", op),
		)
		Index(log, w, r)
	}
}

func BackgroundHandler(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.BackgroundHandler"
		log = log.With(
			slog.String("op", op),
		)
		UploadBackground(log, w, r)
	}
}

func PreviewHandler(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.PreviewHandler"
		log = log.With(
			slog.String("op", op),
		)
		switch r.Method {
		case http.MethodGet:
			GetPreview(log, w, r)
		case http.MethodPost:
			PostPreview(log, w, r)
		}
	}
}

func GenerationHandler(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GenerationHandler"
		slog.Info(op)
		log = log.With(
			slog.String("op", op),
		)
		GenerationQR(log, w, r)
	}
}
