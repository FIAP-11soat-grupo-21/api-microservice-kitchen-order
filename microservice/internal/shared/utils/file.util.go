package utils

import (
	"mime/multipart"
	"slices"
)

func FileIsImage(file multipart.FileHeader) bool {
	contentType := file.Header.Get("Content-Type")

	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	}

	return slices.Contains(validTypes, contentType)
}
