Batch QR Code Generator

This application places a QR code on the background image, the data for QR codes is taken from a text file. The output will be a folder with the finished images in jpg format. 
Each image will contain a qr code with an encrypted string from the data file.

## Quick start

### CLI

For generate images with QR-code, run:

```bash
go run ./cmd/  -file ./tests/data.txt -background ./tests/background.jpg -size 120 -h-align 50 -v-align 75 -output output
```

#### Arguments

- `-background` – The path to the background image (default "background.jpg").
- `-file` – The path to the data file for QR codes (default "data.txt").
- `-h-align` – Horizontal alignment as a percentage (default 50).
- `-output` – Folder with images (default "output").
- `-qr` – Insert qr or text (default true).
- `-size` – Size of the upper image (default 200).
- `-v-align` – Vertical alignment as a percentage (default 50).

### Text

To insert text, you need a font file in the `ttf` format.

- `-font` – Font for text (default "./fonts/DroidSansMono.ttf").

## Build

Build for Windows:
```bash
env GOOS=windows GOARCH=386 go build -o qr-gen.exe
```