package qrgen

import (
	"os"
	"testing"

	"github.com/TOsmanov/qr-gen/internal/lib/utils"
	"github.com/stretchr/testify/assert"
)

const ext = ".jpg"

func TestGenerationQR(t *testing.T) {
	list, err := PrepareData("../tests/data.txt")
	assert.Nil(t, err)
	img, err := PrepareBackground("../tests/background.jpg")
	assert.Nil(t, err)

	params := Params{
		Data:            list,
		Size:            120,
		QRmode:          true,
		BackgroundImg:   img,
		HorizontalAlign: 50,
		VerticalAlign:   75,
		Output:          "../tests/output/",
		Preview:         false,
	}

	err = Generation(params)
	assert.Nil(t, err)

	// Comparing the number of files in a folder
	output, err := os.ReadDir("../tests/output/")
	assert.Nil(t, err)

	expectOutput, err := os.ReadDir("../tests/output/")
	assert.Nil(t, err)
	assert.Equal(t, len(output), len(expectOutput))

	// Comparing sums
	for _, file := range list {
		// Output
		var f1 *os.File
		f1, err = os.Open("../tests/output/" + file + ext)
		assert.Nil(t, err)

		f1.Seek(0, 0)
		var sum1 string
		sum1, err = utils.FileSumSha256(f1)
		assert.Nil(t, err)

		// Expected output
		var f2 *os.File
		f2, err = os.Open("../tests/expect-output/" + file + ext)
		assert.Nil(t, err)

		f2.Seek(0, 0)
		var sum2 string
		sum2, err = utils.FileSumSha256(f2)
		assert.Nil(t, err)

		// Compare
		assert.Equal(t, sum2, sum1)
	}

	if err == nil {
		os.RemoveAll("../tests/output/")
	}
}

func TestGenerationQRPreview(t *testing.T) {
	img, err := PrepareBackground("../tests/background.jpg")
	assert.Nil(t, err)

	params := Params{
		Size:            120,
		QRmode:          true,
		BackgroundImg:   img,
		HorizontalAlign: 50,
		VerticalAlign:   75,
		Output:          "../tests/",
		Preview:         true,
	}

	err = Generation(params)
	assert.Nil(t, err)

	// Output
	var f1 *os.File
	f1, err = os.Open("../tests/preview.jpg")
	assert.Nil(t, err)

	f1.Seek(0, 0)
	var sum1 string
	sum1, err = utils.FileSumSha256(f1)
	assert.Nil(t, err)

	// Expected output
	var f2 *os.File
	f2, err = os.Open("../tests/expect-preview.jpg")
	assert.Nil(t, err)

	f2.Seek(0, 0)
	var sum2 string
	sum2, err = utils.FileSumSha256(f2)
	assert.Nil(t, err)

	// Compare
	assert.Equal(t, sum2, sum1)

	if err == nil {
		os.Remove("../tests/preview.jpg")
	}
}

func TestGenerationQRErrors(t *testing.T) {
	list, err := PrepareData("../tests/data.txt")
	assert.Nil(t, err)
	img, err := PrepareBackground("../tests/background.jpg")
	assert.Nil(t, err)

	params := Params{
		Data:            list,
		Size:            -1,
		QRmode:          true,
		BackgroundImg:   img,
		HorizontalAlign: 50,
		VerticalAlign:   75,
		Output:          "../tests/output/",
		Preview:         false,
	}

	err = Generation(params)
	assert.NotNil(t, err)
}

func TestGenerationText(t *testing.T) {
	list, err := PrepareData("../tests/data.txt")
	assert.Nil(t, err)
	img, err := PrepareBackground("../tests/background.jpg")
	assert.Nil(t, err)

	params := Params{
		Data:            list,
		Size:            120,
		QRmode:          false,
		BackgroundImg:   img,
		HorizontalAlign: 50,
		VerticalAlign:   75,
		Font:            "../tests/font/DroidSansMono.ttf",
		Output:          "../tests/output/",
		Preview:         false,
	}

	err = Generation(params)
	assert.Nil(t, err)

	// Comparing the number of files in a folder
	output, err := os.ReadDir("../tests/output/")
	assert.Nil(t, err)

	expectOutput, err := os.ReadDir("../tests/output/")
	assert.Nil(t, err)
	assert.Equal(t, len(output), len(expectOutput))

	// Comparing sums
	for _, file := range list {
		// Output
		var f1 *os.File
		f1, err = os.Open("../tests/output/" + file + ext)
		assert.Nil(t, err)

		f1.Seek(0, 0)
		var sum1 string
		sum1, err = utils.FileSumSha256(f1)
		assert.Nil(t, err)

		// Expected output
		var f2 *os.File
		f2, err = os.Open("../tests/expect-output-text/" + file + ext)
		assert.Nil(t, err)

		f2.Seek(0, 0)
		var sum2 string
		sum2, err = utils.FileSumSha256(f2)
		assert.Nil(t, err)

		// Compare
		assert.Equal(t, sum2, sum1)
	}

	if err == nil {
		os.RemoveAll("../tests/output/")
	}
}
