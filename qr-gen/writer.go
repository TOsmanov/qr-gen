package qrgen

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func Archive(inputFolder string, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		filename := filepath.Base(path)
		f, err := w.Create(filename)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}
		return nil
	}
	err = filepath.Walk(inputFolder, walker)
	if err != nil {
		return err
	}
	return nil
}
