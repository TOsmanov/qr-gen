package qrgen

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepairData(t *testing.T) {
	var list []string
	var err error
	_, err = PrepareData("bad/path.txt")
	assert.NotNil(t, err)
	list, err = PrepareData("../tests/data.txt")
	assert.Nil(t, err)
	expected := []string{
		"12300534643",
		"4735786789",
		"FHKR759Skl6378993",
		"BRGJHMOW525",
		"QWE12GU",
		"QWE12",
		"QWE",
		"QW",
	}
	assert.Equal(t, list, expected)
}

func TestPrepareBackground(t *testing.T) {
	var img image.Image
	var err error
	_, err = PrepareBackground("bad/img.jpg")
	assert.NotNil(t, err)
	img, err = PrepareBackground("../tests/background.jpg")
	assert.Nil(t, err)
	// Test random pixels
	assert.Equal(t, color.YCbCr{Y: 0x1e, Cb: 0x87, Cr: 0x7d}, img.At(0, 0))
	assert.Equal(t, color.YCbCr{Y: 0xa7, Cb: 0x87, Cr: 0x7d}, img.At(5, 5))
	assert.Equal(t, color.YCbCr{Y: 0x12, Cb: 0x7d, Cr: 0x81}, img.At(30, 50))
	assert.Equal(t, color.YCbCr{Y: 0xbb, Cb: 0x88, Cr: 0x7a}, img.At(100, 162))
}
