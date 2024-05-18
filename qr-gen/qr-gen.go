package qrgen

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"

	"github.com/fogleman/gg"
	"github.com/skip2/go-qrcode"
)

func Generation(list []string,
	size int, qr bool, backgroundImg image.Image, font string,
	hAlign float64, vAlign float64, output string) {
	var err error
	for _, data := range list {
		var upperImg image.Image
		if qr {
			upperImg, err = prepareQR(size, data)
			if err != nil {
				log.Fatalf("Error write file: %v", err)
			}
		} else {
			upperImg = prepareText(size, font, data)
			if err != nil {
				log.Fatalf("Error write file: %v", err)
			}
		}
		x := int(float64(backgroundImg.Bounds().Dx())*hAlign) - size/2
		y := int(float64(backgroundImg.Bounds().Dy())*vAlign) - size/2
		point := image.Point{-x, -y}
		r := image.Rectangle{image.Point{0, 0}, backgroundImg.Bounds().Max}
		rgba := image.NewRGBA(r)
		draw.Draw(rgba, backgroundImg.Bounds(), backgroundImg, image.Point{0, 0}, draw.Src)
		draw.Draw(rgba, backgroundImg.Bounds(), upperImg, point, draw.Src)
		os.Mkdir(output, 0750)
		out, err := os.Create(fmt.Sprintf("%s/%s.jpg", output, data))
		if err != nil {
			fmt.Println(err)
		}

		var opt jpeg.Options
		opt.Quality = 80

		jpeg.Encode(out, rgba, &opt)
	}
}

func prepareQR(qrSize int, data string) (image.Image, error) {
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	qrImg := qr.Image(qrSize)
	// if err != nil {
	// 	return nil, err
	// }
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
