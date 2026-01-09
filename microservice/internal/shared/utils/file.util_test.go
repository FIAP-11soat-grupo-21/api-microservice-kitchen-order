package utils

import (
	"mime/multipart"
	"net/textproto"
	"testing"
)

func TestFileIsImage_ValidImageTypes(t *testing.T) {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	}

	for _, contentType := range validTypes {
		header := make(textproto.MIMEHeader)
		header.Set("Content-Type", contentType)
		
		file := multipart.FileHeader{
			Filename: "test.jpg",
			Header:   header,
		}

		result := FileIsImage(file)

		if !result {
			t.Errorf("Expected FileIsImage to return true for %s, got false", contentType)
		}
	}
}

func TestFileIsImage_InvalidTypes(t *testing.T) {
	invalidTypes := []string{
		"text/plain",
		"application/pdf",
		"video/mp4",
		"audio/mp3",
		"application/json",
	}

	for _, contentType := range invalidTypes {
		header := make(textproto.MIMEHeader)
		header.Set("Content-Type", contentType)
		
		file := multipart.FileHeader{
			Filename: "test.txt",
			Header:   header,
		}

		result := FileIsImage(file)

		if result {
			t.Errorf("Expected FileIsImage to return false for %s, got true", contentType)
		}
	}
}

func TestFileIsImage_EmptyContentType(t *testing.T) {
	header := make(textproto.MIMEHeader)
	
	file := multipart.FileHeader{
		Filename: "test.txt",
		Header:   header,
	}

	result := FileIsImage(file)

	if result {
		t.Error("Expected FileIsImage to return false for empty content type, got true")
	}
}