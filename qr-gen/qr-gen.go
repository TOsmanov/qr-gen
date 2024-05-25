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

func Generation(list []string,
	size int, qr bool, backgroundImg image.Image, font string,
	horizontalAlign int, verticalAlign int, output string, preview bool,
) error {
	hAlign := float64(horizontalAlign) / 100
	vAlign := float64(verticalAlign) / 100
	var err error
	if preview {
		list = []string{
			"https://github.com/TOsmanov/qr-gen",
		}
	}
	for _, data := range list {
		var upperImg image.Image
		var filename string
		if qr {
			upperImg, err = prepareQR(size, data)
			if err != nil {
				return err
			}
		} else {
			upperImg = prepareText(size, font, data)
			if err != nil {
				return err
			}
		}
		x := int(float64(backgroundImg.Bounds().Dx())*hAlign) - size/2
		y := int(float64(backgroundImg.Bounds().Dy())*vAlign) - size/2
		point := image.Point{-x, -y}
		r := image.Rectangle{image.Point{0, 0}, backgroundImg.Bounds().Max}
		rgba := image.NewRGBA(r)
		draw.Draw(rgba, backgroundImg.Bounds(), backgroundImg, image.Point{0, 0}, draw.Src)
		draw.Draw(rgba, backgroundImg.Bounds(), upperImg, point, draw.Src)
		os.Mkdir(output, 0o750)
		if preview {
			filename = "preview"
		} else {
			regex := regexp.MustCompile(`[htps]*://|/|\\|\s`)
			filename = regex.ReplaceAllString(data, "")
		}
		var out *os.File
		out, err = os.Create(fmt.Sprintf("%s/%s.jpg", output, filename))
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

func prepareText(size int, font string, data string) image.Image {
	fontSize := (float64(size) / float64(len([]rune(data)))) * 1.5
	dc := gg.NewContext(size, int(fontSize*1.4))
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace(font, fontSize); err != nil {
		panic(err)
	}
	dc.DrawStringAnchored(data, float64(size/2), fontSize*1.3/2, 0.5, 0.5)
	return dc.Image()
}

func SumSha256(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// TODO Clean Temp files
