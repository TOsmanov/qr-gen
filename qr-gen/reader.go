package qrgen

import (
	"bufio"
	"image"
	"os"
)

func PrepareBackground(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func PrepareData(path string) ([]string, error) {
	var list []string
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		list = append(list, line)
	}
	return list, nil
}
