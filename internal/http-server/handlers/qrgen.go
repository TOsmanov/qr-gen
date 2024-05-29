package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/TOsmanov/qr-gen/internal/config"
	"github.com/TOsmanov/qr-gen/internal/lib/api/response"
	qrgen "github.com/TOsmanov/qr-gen/qr-gen"
	"github.com/go-chi/render"
)

func Index(log *slog.Logger, w http.ResponseWriter,
	_ *http.Request, cfg *config.Config,
) {
	tpl := template.Must(template.ParseFiles(cfg.MainPage))
	tpl.Execute(w, nil)
	log.Info("The main page has been sent successfully")
}

func GetPreview(log *slog.Logger, w http.ResponseWriter,
	r *http.Request, cfg *config.Config,
) {
	const op = "handlers.OrderHandler.GetPreview"
	buf, err := os.ReadFile(cfg.PreviewPath)
	if err != nil {
		w.WriteHeader(400)
		log.Error(
			"Failed to prepare background",
			op, err)
		render.JSON(w, r,
			response.Error("Failed to get preview"))
		return
	}

	w.Header().Set("Content-Type", "image/jpg")
	w.Write(buf)
	os.Remove(cfg.PreviewPath)
}

func UploadBackground(log *slog.Logger, w http.ResponseWriter,
	r *http.Request, cfg *config.Config,
) {
	const op = "handlers.OrderHandler.UploadBackground"
	r.ParseMultipartForm(32 << 20) // 32 MB
	file, _, err := r.FormFile("img")
	if err != nil {
		w.WriteHeader(400)
		log.Error("Failed to read image from request", op, err)
		responseFail(w, r, "Failed to read image from request")
		return
	}
	defer file.Close()

	var data []byte
	data, err = io.ReadAll(file)
	if err != nil {
		w.WriteHeader(500)
		log.Error("Failed to read data", op, err)
		responseFail(w, r, "Failed to read data")
		return
	}

	// Format validation
	formats := []string{"image/jpg", "image/jpeg", "image/png"}
	s := http.DetectContentType(data)

	if !existInSlice(s, formats) {
		w.WriteHeader(400)
		log.Error("Failed to validation file type", op, err)
		responseFail(w, r, "Failed to validation file type")
		return
	}

	// Prepare temp file
	sum := qrgen.SumSha256(data)
	outputJpg := fmt.Sprintf("%s/%s.jpg", cfg.TempDir, sum)
	err = os.MkdirAll(cfg.TempDir, os.ModePerm)
	if err != nil {
		w.WriteHeader(500)
		log.Error("Failed to create temp directory", op, err)
		responseFail(w, r, "Failed to create temp directory")
		return
	}

	// Write background file (if not exist)
	if _, err = os.Stat(outputJpg); errors.Is(err, os.ErrNotExist) {
		var tempFile *os.File
		tempFile, err = os.Create(outputJpg)
		if err != nil {
			w.WriteHeader(500)
			log.Error("Failed to create temp file", op, err)
			responseFail(w, r, "Failed to create temp file")
			return
		}
		defer tempFile.Close()
		_, err = tempFile.Write(data)
		if err != nil {
			w.WriteHeader(500)
			log.Error("Failed to write temp file", op, err)
			responseFail(w, r, "Failed to write temp file")
			return
		}
	}
	log.Info("The background image has been uploaded successfully")
	responseOK(w, r, sum)
}

func PostPreview(log *slog.Logger, w http.ResponseWriter,
	r *http.Request, cfg *config.Config,
) {
	const op = "handlers.OrderHandler.PostPreview"
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

	params.BackgroundImg = fmt.Sprintf("%s/%s.jpg", cfg.TempDir, params.BackgroundImg)
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
	const op = "handlers.OrderHandler.Generation"
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

	backgroundHash := params.BackgroundImg

	params.BackgroundImg = fmt.Sprintf("%s/%s.jpg", cfg.TempDir, params.BackgroundImg)

	// Make temp directory
	params.Output = fmt.Sprintf("%s/%s", cfg.TempDir, params.BackgroundImg)
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
	outputZip := fmt.Sprintf("%s/%s.zip", cfg.SiteDir, backgroundHash)
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

	// Clean
	os.RemoveAll(params.Output)
	os.Remove(outputZip)
}

func existInSlice(s string, list []string) bool {
	for _, format := range list {
		if s == format {
			return true
		}
	}
	return false
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
