package handlers

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/TOsmanov/qr-gen/internal/config"
	"github.com/TOsmanov/qr-gen/internal/lib/api/response"
	qrgen "github.com/TOsmanov/qr-gen/qr-gen"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

var static struct {
	MainPage   []byte
	Preview    []byte
	Background image.Image
}

func Index(log *slog.Logger, w http.ResponseWriter,
	r *http.Request, cfg *config.Config,
) {
	const op = "handlers.IndexHandler.Index"
	if len(static.MainPage) == 0 {
		var err error
		static.MainPage, err = os.ReadFile(cfg.MainPage)
		if err != nil {
			w.WriteHeader(500)
			log.Error(
				"Failed to read main page",
				op, err)
			render.JSON(w, r,
				response.Error("Failed to get main page"))
			return
		}
	}
	w.Write(static.MainPage)
	log.Info("The main page has been sent successfully")
}

func GetPreview(log *slog.Logger, w http.ResponseWriter,
	r *http.Request, cfg *config.Config,
) {
	const op = "handlers.PreviewHandler.GetPreview"
	var err error
	static.Preview, err = os.ReadFile(cfg.PreviewPath)
	if err != nil {
		w.WriteHeader(500)
		log.Error(
			"Failed to prepare background",
			op, err)
		render.JSON(w, r,
			response.Error("Failed to get preview"))
		return
	}

	w.Header().Set("Content-Type", "image/jpg")
	w.Write(static.Preview)
	os.Remove(cfg.PreviewPath)
}

func UploadBackground(log *slog.Logger, w http.ResponseWriter,
	r *http.Request,
) {
	const op = "handlers.BackgroundHandler.UploadBackground"
	r.ParseMultipartForm(32 << 20) // 32 MB
	file, _, err := r.FormFile("img")
	if err != nil {
		w.WriteHeader(400)
		log.Error("Failed to read image from request", op, err)
		responseFail(w, r, "Failed to read image from request")
		return
	}
	defer file.Close()

	static.Background, _, err = image.Decode(file)
	if err != nil {
		w.WriteHeader(400)
		log.Error("Failed to decode data", op, err)
		responseFail(w, r, "Failed to decode data")
		return
	}

	log.Info("The background image has been uploaded successfully")
	responseOK(w, r, "The background has been successfully upload")
}

func PostPreview(log *slog.Logger, w http.ResponseWriter,
	r *http.Request, cfg *config.Config,
) {
	const op = "handlers.PreviewHandler.PostPreview"
	var params qrgen.Params

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		log.Error("Failed to read request", op, err)
		render.JSON(w, r, response.Error("Failed to read request"))
		return
	}
	defer r.Body.Close()

	json.Unmarshal(b, &params)

	params.BackgroundImg = static.Background
	params.QRmode, params.Preview = true, true
	params.Output = cfg.SiteDir

	err = qrgen.Generation(params)
	if err != nil {
		w.WriteHeader(500)
		log.Error(
			"Failed to generate preview",
			op, err)
		render.JSON(w, r,
			response.Error("Failed to generate preview"))
		return
	}
	responseOK(w, r, "The preview has been successfully generate")
}

func GenerationQR(log *slog.Logger, w http.ResponseWriter,
	r *http.Request, cfg *config.Config,
) {
	const op = "handlers.GenerationHandler.Generation"
	var params qrgen.Params

	// Parameters
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		log.Error("Failed to read request", op, err)
		render.JSON(w, r, response.Error("Failed to read request"))
		return
	}
	defer r.Body.Close()

	json.Unmarshal(b, &params)

	id := uuid.New()
	params.BackgroundImg = static.Background

	// Make temp directory
	params.Output = fmt.Sprintf("%s/%s", cfg.TempDir, id.String())
	err = os.MkdirAll(params.Output, os.ModePerm)
	if err != nil {
		w.WriteHeader(500)
		log.Error("Failed to create temp directory", op, err)
		render.JSON(w, r, response.Error("Failed to create temp directory"))
		return
	}

	// Generate images
	params.QRmode, params.Preview = true, false
	err = qrgen.Generation(params)
	if err != nil {
		w.WriteHeader(500)
		log.Error(
			"Failed to generation",
			op, err)
		render.JSON(w, r,
			response.Error("Failed to generation"))
		return
	}

	// Archiving
	outputZip := fmt.Sprintf("%s/%s.zip", cfg.SiteDir, id.String())
	qrgen.Archive(params.Output, outputZip)

	buf, err := os.ReadFile(outputZip)
	if err != nil {
		w.WriteHeader(500)
		log.Error(
			"Failed to read archive",
			op, err)
		render.JSON(w, r,
			response.Error("Failed to get archive"))
		return
	}
	w.Write(buf)
	os.Remove(outputZip)
}

func responseOK(
	w http.ResponseWriter,
	r *http.Request,
	body string,
) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Body:     body,
	})
}

func responseFail(
	w http.ResponseWriter,
	r *http.Request,
	msg string,
) {
	render.JSON(w, r, Response{
		Response: response.Error(msg),
	})
}
