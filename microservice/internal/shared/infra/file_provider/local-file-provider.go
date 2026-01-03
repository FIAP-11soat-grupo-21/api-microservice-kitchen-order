package file_service

import (
	"os"
	"path/filepath"
)

type LocalFileProvider struct {
	basePath string
}

func NewLocalFileProvider() *LocalFileProvider {
	uploadedFilesPath := filepath.Join(".", "uploads")

	if err := os.MkdirAll(uploadedFilesPath, os.ModePerm); err != nil {
		panic("Failed to create uploads directory: " + err.Error())
	}

	if err := os.Chmod(uploadedFilesPath, 0755); err != nil {
		panic("Failed to set permissions for uploads directory: " + err.Error())
	}

	return &LocalFileProvider{
		basePath: uploadedFilesPath,
	}
}

func (l *LocalFileProvider) UploadFile(fileName string, fileContent []byte) error {
	fileExits := l.fileExists(fileName)

	if fileExits {
		return os.ErrExist
	}

	filePath := filepath.Join(l.basePath, fileName)

	return os.WriteFile(filePath, fileContent, 0644)
}

func (l *LocalFileProvider) DeleteFile(fileName string) error {
	fileExits := l.fileExists(fileName)

	if !fileExits {
		return os.ErrNotExist
	}

	filePath := filepath.Join(l.basePath, fileName)

	return os.Remove(filePath)
}

func (l *LocalFileProvider) fileExists(fileName string) bool {
	filePath := filepath.Join(l.basePath, fileName)

	_, err := os.Stat(filePath)

	return !os.IsNotExist(err)
}
