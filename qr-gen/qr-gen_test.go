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
	os.RemoveAll("../tests/output/")
	list, err := PrepareData("../tests/data.txt")
	assert.Nil(t, err)

	params := Params{
		Data:            list,
		Size:            120,
		QRmode:          true,
		BackgroundImg:   "../tests/background.jpg",
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

func fileSumSha256(f *os.File) (string, error) {
	file1Sum := sha256.New()
	if _, err := io.Copy(file1Sum, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", file1Sum.Sum(nil)), nil
}

// TODO: func TestGenerationText(t *testing.T) {}
