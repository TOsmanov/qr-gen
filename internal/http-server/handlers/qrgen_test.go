package handlers

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"log/slog"
	"mime/multipart"
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

	assert.Equal(t, http.StatusOK, w.Code)

	b, err := os.ReadFile("../../../site/index.html")
	assert.Nil(t, err)
	expectedBody := string(b)

	assert.Equal(t, w.Body.String(), expectedBody)
}

func TestGetPreview(t *testing.T) {
	validData := []byte("{}")
	r := httptest.NewRequest(http.MethodGet, "/preview", bytes.NewBuffer(validData))
	w := httptest.NewRecorder()
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	cfg := config.Config{
		Env: "local",
	}
	err := os.Link("../../../tests/expect-preview.jpg", "../../../site/preview.jpg")
	assert.Nil(t, err)

	cfg.PreviewPath = "../../../site/preview.jpg"
	GetPreview(log, w, r, &cfg)

	assert.Equal(t, http.StatusOK, w.Code)

	b, err := os.ReadFile("../../../tests/expect-preview.jpg")
	assert.Nil(t, err)
	expectedBody := string(b)

	assert.Equal(t, w.Body.String(), expectedBody)
}

func TestUploadBackground(t *testing.T) {
	t.Run("Test UploadBackground", func(t *testing.T) {
		pipeReader, pipeWriter := io.Pipe()
		multipartWriter := multipart.NewWriter(pipeWriter)
		go func() {
			defer multipartWriter.Close()

			fileField, err := multipartWriter.CreateFormFile("img", "background.jpg")
			assert.Nil(t, err)

			fileBytes, err := os.ReadFile("../../../tests/background.jpg")
			assert.Nil(t, err)

			reader := bytes.NewReader(fileBytes)
			image, err := jpeg.Decode(reader)
			assert.Nil(t, err)

			err = jpeg.Encode(fileField, image, &jpeg.Options{Quality: 80})
			assert.Nil(t, err)

			fmt.Println(fileField)
		}()
		r := httptest.NewRequest(http.MethodPost, "/background", pipeReader)
		r.Header.Add("content-type", multipartWriter.FormDataContentType())
		w := httptest.NewRecorder()
		log := slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

		UploadBackground(log, w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		expectedBody := "{\"status\":\"OK\",\"data\":\"The background has been successfully upload\"}\n"
		assert.Equal(t, w.Body.String(), expectedBody)
	})
}

func TestPostPreview(t *testing.T) {
	cfg := config.Config{
		Env: "local",
	}

	TestUploadBackground(t)

	cfg.SiteDir = "../../site"
	validData := []byte("{\"size\":120,\"hAlign\":50,\"vAlign\":70}")

	r := httptest.NewRequest(http.MethodPost, "/preview", bytes.NewBuffer(validData))
	w := httptest.NewRecorder()
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	PostPreview(log, w, r, &cfg)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedBody := "{\"status\":\"OK\",\"data\":\"The preview has been successfully generate\"}\n"

	assert.Equal(t, expectedBody, w.Body.String())
}

func TestGenerationQR(t *testing.T) {
	cfg := config.Config{
		Env: "local",
	}
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	validData := []byte(
		"{\"list\":[\"123\"],\"size\":120,\"hAlign\":50,\"vAlign\":75}")
	cfg.TempDir = "../../../temp"
	cfg.SiteDir = "../../../tmp-site"
	err := os.MkdirAll(cfg.SiteDir, os.ModePerm)
	assert.Nil(t, err)

	TestUploadBackground(t)

	r := httptest.NewRequest(http.MethodPost, "/generation", bytes.NewBuffer(validData))
	w := httptest.NewRecorder()

	GenerationQR(log, w, r, &cfg)

	assert.Equal(t, http.StatusOK, w.Code)

	expect, err := os.ReadFile("../../../tests/except-archive_2.zip")
	assert.Nil(t, err)

	assert.Equal(t, len(expect), len(w.Body.Bytes()))
	os.RemoveAll("../../../tmp-site")
	os.RemoveAll("../../../temp")
}
