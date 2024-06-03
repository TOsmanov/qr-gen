package qrgen

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		f1, err = os.Open("../tests/output/" + file + ".jpg")
		assert.Nil(t, err)

		f1.Seek(0, 0)
		var sum1 string
		sum1, err = fileSumSha256(f1)
		assert.Nil(t, err)

		// Expected output
		var f2 *os.File
		f2, err = os.Open("../tests/expect-output/" + file + ".jpg")
		assert.Nil(t, err)

		f2.Seek(0, 0)
		var sum2 string
		sum2, err = fileSumSha256(f2)
		assert.Nil(t, err)

		// Compare
		assert.Equal(t, sum2, sum1)
	}

	os.RemoveAll("../tests/output/")
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
	sum1, err = fileSumSha256(f1)
	assert.Nil(t, err)

	// Expected output
	var f2 *os.File
	f2, err = os.Open("../tests/expect-preview.jpg")
	assert.Nil(t, err)

	f2.Seek(0, 0)
	var sum2 string
	sum2, err = fileSumSha256(f2)
	assert.Nil(t, err)

	// Compare
	assert.Equal(t, sum2, sum1)

	os.Remove("../tests/preview.jpg")
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

func fileSumSha256(f *os.File) (string, error) {
	file1Sum := sha256.New()
	if _, err := io.Copy(file1Sum, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", file1Sum.Sum(nil)), nil
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
		f1, err = os.Open("../tests/output/" + file + ".jpg")
		assert.Nil(t, err)

		f1.Seek(0, 0)
		var sum1 string
		sum1, err = fileSumSha256(f1)
		assert.Nil(t, err)

		// Expected output
		var f2 *os.File
		f2, err = os.Open("../tests/expect-output-text/" + file + ".jpg")
		assert.Nil(t, err)

		f2.Seek(0, 0)
		var sum2 string
		sum2, err = fileSumSha256(f2)
		assert.Nil(t, err)

		// Compare
		assert.Equal(t, sum2, sum1)
	}

	os.RemoveAll("../tests/output/")
}
