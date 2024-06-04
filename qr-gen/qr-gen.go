package qrgen

import (
	"crypto/sha256"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"regexp"

	"github.com/fogleman/gg"
	"github.com/skip2/go-qrcode"
)

type Params struct {
	Data            []string `json:"list,omitempty"`
	Size            int      `json:"size"`
	HorizontalAlign int      `json:"hAlign"`
	VerticalAlign   int      `json:"vAlign"`
	BackgroundImg   image.Image
	QRmode          bool
	Font            string
	Output          string
	Preview         bool
}

func Generation(params Params) error {
	if (params.Size <= 0) && !(params.HorizontalAlign < 0) && !(params.VerticalAlign < 0) {
		return fmt.Errorf("numeric parameters must be greater than zero")
	}
	hAlign := float64(params.HorizontalAlign) / 100
	vAlign := float64(params.VerticalAlign) / 100
	var err error

	if params.Preview {
		params.Data = []string{
			"https://github.com/TOsmanov/qr-gen",
		}
	}
	for _, data := range params.Data {
		var upperImg image.Image
		var filename string
		if params.QRmode {
			upperImg, err = prepareQR(params.Size, data)
			if err != nil {
				return err
			}
		} else {
			upperImg, err = prepareText(params.Size, params.Font, data)
			if err != nil {
				return err
			}
		}
		x := int(float64(params.BackgroundImg.Bounds().Dx())*hAlign) - params.Size/2
		y := int(float64(params.BackgroundImg.Bounds().Dy())*vAlign) - params.Size/2
		point := image.Point{-x, -y}
		r := image.Rectangle{image.Point{0, 0}, params.BackgroundImg.Bounds().Max}
		rgba := image.NewRGBA(r)
		draw.Draw(rgba, params.BackgroundImg.Bounds(), params.BackgroundImg, image.Point{0, 0}, draw.Src)
		draw.Draw(rgba, params.BackgroundImg.Bounds(), upperImg, point, draw.Src)
		os.Mkdir(params.Output, 0o750)
		if params.Preview {
			filename = "preview"
		} else {
			regex := regexp.MustCompile(`[htps]*://|/|\\|\s`)
			filename = regex.ReplaceAllString(data, "")
		}
		var out *os.File
		out, err = os.Create(fmt.Sprintf("%s/%s.jpg", params.Output, filename))
		if err != nil {
			return err
		}

		var opt jpeg.Options
		opt.Quality = 80

		jpeg.Encode(out, rgba, &opt)
	}
	return nil
}

func prepareQR(qrSize int, data string) (image.Image, error) {
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	qrImg := qr.Image(qrSize)
	return qrImg, nil
}

func prepareText(size int, font string, data string) (image.Image, error) {
	fontSize := (float64(size) / float64(len([]rune(data)))) * 1.5
	dc := gg.NewContext(size, int(fontSize*1.4))
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace(font, fontSize); err != nil {
		return nil, err
	}
	dc.DrawStringAnchored(data, float64(size/2), fontSize*1.3/2, 0.5, 0.5)
	return dc.Image(), nil
}

func SumSha256(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}
