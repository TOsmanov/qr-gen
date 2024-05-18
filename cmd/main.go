package main

import (
	"flag"
	"log"
	"os"

	qrgen "github.com/TOsmanov/qr-gen/qr-gen"
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
		flag.BoolVar(&qr, "qr", true, "Insert qr or text")
	}
	if flag.Lookup("font") == nil {
		flag.StringVar(&font, "font", "./fonts/DroidSansMono.ttf", "Font for text")
	}
}

func main() {
	flag.Parse()
	hAlign := float64(horizontalAlign) / 100
	vAlign := float64(verticalAlign) / 100

	backgroundImg, err := qrgen.PrepareBackground(background)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	list, err := qrgen.PrepareData(data)
	if err != nil {
		log.Fatalf("Error in data preparation: %v", err)
		os.Exit(1)
	}

	qrgen.Generation(list, size, qr, backgroundImg, font, hAlign, vAlign, output)
}
