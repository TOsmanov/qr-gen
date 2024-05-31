package main

import (
	"flag"
	"log"

	qrgen "github.com/TOsmanov/qr-gen/qr-gen"
)

var (
	params        qrgen.Params
	data          string
	backgroundImg string
)

func init() {
	if flag.Lookup("file") == nil {
		flag.StringVar(&data, "file", "data.txt", "The path to the data file for QR codes")
	}
	if flag.Lookup("background") == nil {
		flag.StringVar(&backgroundImg, "background", "background.jpg", "The path to the background image")
	}
	if flag.Lookup("size") == nil {
		flag.IntVar(&params.Size, "size", 200, "Size of the upper image")
	}
	if flag.Lookup("horizontalAlign") == nil {
		flag.IntVar(&params.HorizontalAlign, "h-align", 50, "Horizontal alignment as a percentage")
	}
	if flag.Lookup("verticalAlign") == nil {
		flag.IntVar(&params.VerticalAlign, "v-align", 50, "Vertical alignment as a percentage")
	}
	if flag.Lookup("output") == nil {
		flag.StringVar(&params.Output, "output", "output", "Folder with images")
	}
	if flag.Lookup("qr") == nil {
		flag.BoolVar(&params.QRmode, "qr", true, "Insert qr or text")
	}
	if flag.Lookup("font") == nil {
		flag.StringVar(&params.Font, "font", "./fonts/DroidSansMono.ttf", "Font for text")
	}
}

func main() {
	flag.Parse()

	list, err := qrgen.PrepareData(data)
	if err != nil {
		log.Fatalf("Error in data preparation: %v", err)
	}
	params.BackgroundImg, err = qrgen.PrepareBackground(backgroundImg)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	params.Data = list
	params.Preview = false

	err = qrgen.Generation(params)
	if err != nil {
		log.Fatalf("Error write file: %v", err)
	}
}
