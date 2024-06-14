package qrgen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/TOsmanov/qr-gen/internal/lib/utils"
	"github.com/stretchr/testify/assert"
)

func TestArchive(t *testing.T) {
	inputAbsPath, err := filepath.Abs("../tests/expect-output")
	assert.Nil(t, err)

	outputPath := "../tests/archive.zip"
	expectOutputPath := "../tests/except-archive.zip"

	err = Archive(inputAbsPath, outputPath)
	assert.Nil(t, err)

	// Output archive
	var f1 *os.File
	f1, err = os.Open(outputPath)
	assert.Nil(t, err)

	f1.Seek(0, 0)
	var sum1 string
	sum1, err = utils.FileSumSha256(f1)
	assert.Nil(t, err)

	// Expected archive
	var f2 *os.File
	f2, err = os.Open(expectOutputPath)
	assert.Nil(t, err)

	f2.Seek(0, 0)
	var sum2 string
	sum2, err = utils.FileSumSha256(f2)
	assert.Nil(t, err)

	// Compare
	assert.Equal(t, sum2, sum1)
	if err == nil {
		os.Remove("../tests/archive.zip")
	}
}
