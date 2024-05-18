package qrgen

import (
	"crypto/md5"
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
	backgroundImg, err := PrepareBackground("../tests/background.jpg")
	assert.Nil(t, err)

	Generation(list, 200, true, backgroundImg,
		"../tests/fonts/DroidSansMono.ttf", 50, 50,
		"../tests/output")

	// Comparing the number of files in a folder
	output, err := os.ReadDir("../tests/output/")
	assert.Nil(t, err)
	expectOutput, err := os.ReadDir("../tests/output/")
	assert.Nil(t, err)
	assert.Equal(t, len(output), len(expectOutput))

	// MD5 sum comparison
	for _, file := range list {
		// Output
		f1, err := os.Open("../tests/output/" + file + ".jpg")
		assert.Nil(t, err)
		f1.Seek(0, 0)
		sum1, err := getMD5SumString(f1)
		assert.Nil(t, err)
		// Expected output
		f2, err := os.Open("../tests/expect-output/" + file + ".jpg")
		assert.Nil(t, err)
		f2.Seek(0, 0)
		sum2, err := getMD5SumString(f2)
		assert.Nil(t, err)
		// Compare
		assert.Equal(t, sum2, sum1)
	}

	os.RemoveAll("../tests/output/")
}

func getMD5SumString(f *os.File) (string, error) {
	file1Sum := md5.New()
	_, err := io.Copy(file1Sum, f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", file1Sum.Sum(nil)), nil
}

// TODO: func TestGenerationText(t *testing.T) {}
