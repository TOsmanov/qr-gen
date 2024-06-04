package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func SumSha256(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

func FileSumSha256(f *os.File) (string, error) {
	file1Sum := sha256.New()
	if _, err := io.Copy(file1Sum, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", file1Sum.Sum(nil)), nil
}
