package util

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func StaticPath(imgPath string) string {
	return strings.Replace(imgPath, os.Getenv("IMAGE_PATH"), "/img", -1)
}

func ValidateImageFile(file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	defer src.Close()

	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return fmt.Errorf("unable to read file: %v", err)
	}

	contentType := http.DetectContentType(buffer)
	if contentType != "image/jpeg" && contentType != "image/png" {
		return fmt.Errorf("invalid file type: %s. Only JPEG and PNG images are allowed", contentType)
	}

	// Reset the read pointer to the beginning of the file
	if _, err := src.Seek(0, 0); err != nil {
		return fmt.Errorf("unable to reset file pointer: %v", err)
	}

	log.Printf("##################################################\n")
	log.Printf("# validated file!                                #\n")
	log.Printf("##################################################\n")

	return nil
}

func SaveFile(file *multipart.FileHeader) (string, error) {
	err := ValidateImageFile(file)
	if err != nil {
		return "", err
	}

	// Destination
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

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

func IsDevelopment() bool {
	return os.Getenv("RUN_MODE") == "development"
}

// Name description
func IsProduction() bool {
	return os.Getenv("RUN_MODE") == "production"
}

func IsEmptyOrWhitespace(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
