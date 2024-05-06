package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"

	"github.com/fogleman/gg"
	"github.com/skip2/go-qrcode"
)

var (
	data            string
	background      string
	size            int
	horizontalAlign int
	verticalAlign   int
	output          string
	qr              bool
	font            string
)

func init() {
	if flag.Lookup("file") == nil {
		flag.StringVar(&data, "file", "data.txt", "The path to the data file for QR codes")
	}
	if flag.Lookup("background") == nil {
		flag.StringVar(&background, "background", "background.jpg", "The path to the background image")
	}
	if flag.Lookup("size") == nil {
		flag.IntVar(&size, "size", 200, "Size of the upper image")
	}
	if flag.Lookup("horizontalAlign") == nil {
		flag.IntVar(&horizontalAlign, "h-align", 50, "Horizontal alignment as a percentage")
	}
	if flag.Lookup("verticalAlign") == nil {
		flag.IntVar(&verticalAlign, "v-align", 50, "Vertical alignment as a percentage")
	}
	if flag.Lookup("output") == nil {
		flag.StringVar(&output, "output", "output", "Folder with images")
	}
	if flag.Lookup("qr") == nil {
		flag.BoolVar(&qr, "qr", false, "Insert qr or text")
	}
	if flag.Lookup("font") == nil {
		flag.StringVar(&font, "font", "./font/DroidSansMono.ttf", "Font")
	}
}

func main() {
	flag.Parse()
	hAlign := float64(horizontalAlign) / 100
	vAlign := float64(verticalAlign) / 100
	file, err := os.Open(background)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	backgroundImg, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Error %v", err)
	}

	list, err := prepareData(data)
	if err != nil {
		log.Fatalf("Error in data preparation: %v", err)
		os.Exit(1)
	}

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

func prepareData(path string) ([]string, error) {
	var list []string
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		list = append(list, line)
	}
	return list, nil
}

func prepareQR(qrSize int, data string) (image.Image, error) {
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	qrImg := qr.Image(qrSize)
	if err != nil {
		return nil, err
	}
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
