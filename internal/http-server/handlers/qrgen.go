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

	"github.com/TOsmanov/qr-gen/internal/lib/api/response"
	qrgen "github.com/TOsmanov/qr-gen/qr-gen"
	"github.com/go-chi/render"
)

type QRParams struct {
	List       []string `json:"list,omitempty"`
	Size       int      `json:"size"`
	Background string   `json:"background"`
	HorizAlign int      `json:"hAlign"`
	VertAlign  int      `json:"vAlign"`
}

func Index(log *slog.Logger, w http.ResponseWriter,
	_ *http.Request,
) {
	tpl := template.Must(template.ParseFiles("site/index.html"))
	tpl.Execute(w, nil)
	log.Info("The main page has been sent successfully")
}

func GetPreview(log *slog.Logger, w http.ResponseWriter,
	r *http.Request,
) {
	const op = "handlers.OrderHandler.GetPreview"
	buf, err := os.ReadFile("site/preview.jpg")
	if err != nil {
		w.WriteHeader(400)
		log.Error(
			"Failed to prepare background",
			op, err)
		render.JSON(w, r,
			response.Error("Failed to get preview"))
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(buf)
	os.Remove("site/preview.jpg")
}

func UploadBackground(log *slog.Logger, w http.ResponseWriter,
	r *http.Request,
) {
	const op = "handlers.OrderHandler.UploadBackground"
	const errorMsg = "Error uploading the file"
	r.ParseMultipartForm(32 << 20) // 32 MB
	file, _, err := r.FormFile("img")
	if err != nil {
		w.WriteHeader(400)
		log.Error("Failed to read image from request", op, err)
		responseFail(w, r, errorMsg)
		return
	}
	defer file.Close()

	var data []byte
	data, err = io.ReadAll(file)
	if err != nil {
		w.WriteHeader(500)
		log.Error("Failed to read data", op, err)
		responseFail(w, r, errorMsg)
		return
	}

	var formats = []string{"image/jpg", "image/jpeg", "image/png"}
	s := http.DetectContentType(data)
	log.Info("s", s, existInSlice(s, formats))

	if !existInSlice(s, formats) {
		w.WriteHeader(400)
		log.Error("Failed to validation file type", op, err)
		responseFail(w, r, errorMsg)
		return
	}

	sum := qrgen.SumSha256(data)
	const tempDir = "./temp/"
	outputJpg := fmt.Sprintf("%s%s.jpg", tempDir, sum)
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		w.WriteHeader(500)
		log.Error("Failed to create temp directory", op, err)
		responseFail(w, r, errorMsg)
		return
	}

	if _, err = os.Stat(outputJpg); errors.Is(err, os.ErrNotExist) {
		var tempFile *os.File
		tempFile, err = os.Create(outputJpg)
		if err != nil {
			w.WriteHeader(500)
			log.Error("Failed to create temp file", op, err)
			responseFail(w, r, errorMsg)
			return
		}
		defer tempFile.Close()
		_, err = tempFile.Write(data)
		if err != nil {
			w.WriteHeader(500)
			log.Error("Failed to write temp file", op, err)
			responseFail(w, r, errorMsg)
			return
		}
	}
	log.Info("The background image has been uploaded successfully")
	responseOK(w, r, sum)
}

func PostPreview(log *slog.Logger, w http.ResponseWriter,
	r *http.Request,
) {
	const op = "handlers.OrderHandler.PostPreview"
	var params QRParams
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		log.Error("Failed to read request", op, err)
		render.JSON(w, r, response.Error("Failed to read request"))
		return
	}
	defer r.Body.Close()

	json.Unmarshal(b, &params)

	backgroundImg, err := qrgen.PrepareBackground("./temp/" + params.Background + ".jpg")
	if err != nil {
		w.WriteHeader(500)
		log.Error(
			"Failed to prepare background",
			op, err)
		render.JSON(w, r,
			response.Error("Failed to prepare preview"))
		return
	}

	var list []string

	qrgen.Generation(list, params.Size, true, backgroundImg, "", params.HorizAlign, params.VertAlign, "site", true)
	buf, err := os.ReadFile("temp/preview/preview.jpg")
	if err != nil {
		w.WriteHeader(400)
		log.Error(
			"Failed to read preview",
			op, err)
		render.JSON(w, r,
			response.Error("Failed to prepare preview"))
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(buf)
	responseOK(w, r, "The preview has been successfully generate")

	// TODO: Clean()
}

func GenerationQR(log *slog.Logger, w http.ResponseWriter,
	r *http.Request,
) {
	const op = "handlers.OrderHandler.Generation"
	slog.Info(op)

	var params QRParams
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		log.Error("Failed to read request", op, err)
		render.JSON(w, r, response.Error("Failed to read request"))
		return
	}
	defer r.Body.Close()

	json.Unmarshal(b, &params)
	backgroundImg, err := qrgen.PrepareBackground("./temp/" + params.Background + ".jpg")
	if err != nil {
		w.WriteHeader(400)
		log.Error(
			"Failed to prepare background",
			op, err)
		render.JSON(w, r,
			response.Error("Failed to prepare preview"))
		return
	}

	outputZip := fmt.Sprintf("site/%s.zip", params.Background)

	tempDir := "output/" + params.Background
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		w.WriteHeader(500)
		log.Error("Failed to create temp directory", op, err)
		render.JSON(w, r, response.Error("Failed to create temp directory"))
		return
	}

	qrgen.Generation(
		params.List, params.Size, true, backgroundImg,
		"", params.HorizAlign, params.VertAlign, tempDir, false)
	qrgen.Archive(tempDir, outputZip)

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
	os.RemoveAll(tempDir)
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
