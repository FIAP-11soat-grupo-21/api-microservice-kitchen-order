package handlers

import "tech_challenge/internal/shared/interfaces"

type FileHandler struct {
	fileProvider interfaces.IFileProvider
}

func NewFileHandler(fileProvider interfaces.IFileProvider) *FileHandler {
	return &FileHandler{
		fileProvider: fileProvider,
	}
}

func (h *FileHandler) FindFile(fileName string) (string, error) {
	fileUrl, err := h.fileProvider.GetPresignedURL(fileName)

	if err != nil {
		return "", err
	}

	return fileUrl, nil
}
