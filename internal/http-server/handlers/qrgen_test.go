package handlers

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/TOsmanov/qr-gen/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	validData := []byte("{}")
	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(validData))
	w := httptest.NewRecorder()
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	cfg := config.Config{
		Env: "local",
	}
	cfg.MainPage = "../../../site/index.html"
	Index(log, w, r, &cfg)

	assert.Equal(t, w.Code, http.StatusOK)

	b, err := os.ReadFile("../../../site/index.html")
	assert.Nil(t, err)
	expectedBody := string(b)

	if w.Body.String() != expectedBody {
		t.Errorf("expected response body %s, but got %s", expectedBody, w.Body.String())
	}
}
