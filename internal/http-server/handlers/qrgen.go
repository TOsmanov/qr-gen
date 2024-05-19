package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"io"
	"log/slog"
	"net/http"

	"github.com/TOsmanov/qr-gen/internal/lib/api/response"
	qrgen "github.com/TOsmanov/qr-gen/qr-gen"
	"github.com/go-chi/render"
)

type QRParams struct {
	list       []string
	size       int
	background image.Image
	hAlign     int
	vAlign     int
}

func Index(log *slog.Logger, w http.ResponseWriter,
	r *http.Request,
) {
	var tpl = template.Must(template.ParseFiles("site/index.html"))
	tpl.Execute(w, nil)
}

func GetPreview(log *slog.Logger, w http.ResponseWriter,
	r *http.Request,
) {
	const op = "handlers.OrderHandler.GetPreview"
	var params QRParams
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("Failed to read request", fmt.Errorf("%s: %w", op, err))
		render.JSON(w, r, response.Error("Failed to read request"))
		return
	}
	defer r.Body.Close()
	json.Unmarshal(b, &params)
	qrgen.Generation(params.list, params.size, false, params.background, "", params.hAlign, params.vAlign, "site", true)
	responseOK(w, r, "The order has been successfully added")
}

func GetDownlaodLink(log *slog.Logger, w http.ResponseWriter,
	r *http.Request,
) {
	const op = "handlers.OrderHandler.GetDownlaodLink"
	slog.Info(op)
	responseOK(w, r, "The order has been successfully added")
}

func responseOK(w http.ResponseWriter, r *http.Request, msg string) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Message:  msg,
	})
}
