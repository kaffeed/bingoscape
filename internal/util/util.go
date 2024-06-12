package util

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func StaticPath(imgPath string) string {
	return strings.Replace(imgPath, os.Getenv("IMAGE_PATH"), "/img", -1)
}

func SaveFile(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Destination

	p := filepath.Join(os.Getenv("IMAGE_PATH"), "tiles")
	dst, err := os.CreateTemp(p, fmt.Sprintf("bingoscape-*%s", path.Ext(file.Filename)))

	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}
	return StaticPath(dst.Name()), nil
}
