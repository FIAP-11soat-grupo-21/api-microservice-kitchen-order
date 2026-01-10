package handlers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFileProvider é um mock do file provider
type MockFileProvider struct {
	mock.Mock
}

func (m *MockFileProvider) UploadFile(fileName string, fileContent []byte) error {
	args := m.Called(fileName, fileContent)
	return args.Error(0)
}

func (m *MockFileProvider) DeleteFile(fileName string) error {
	args := m.Called(fileName)
	return args.Error(0)
}

func (m *MockFileProvider) GetPresignedURL(fileName string) (string, error) {
	args := m.Called(fileName)
	return args.String(0), args.Error(1)
}

func TestNewFileHandler(t *testing.T) {
	// Arrange
	mockProvider := &MockFileProvider{}

	// Act
	handler := NewFileHandler(mockProvider)

	// Assert
	assert.NotNil(t, handler)
	assert.Equal(t, mockProvider, handler.fileProvider)
}

func TestFileHandler_FindFile_Success(t *testing.T) {
	// Arrange
	mockProvider := &MockFileProvider{}
	handler := NewFileHandler(mockProvider)

	fileName := "test-file.jpg"
	expectedURL := "https://example.com/presigned-url/test-file.jpg"

	mockProvider.On("GetPresignedURL", fileName).Return(expectedURL, nil)

	// Act
	result, err := handler.FindFile(fileName)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedURL, result)
	mockProvider.AssertExpectations(t)
}

func TestFileHandler_FindFile_Error(t *testing.T) {
	// Arrange
	mockProvider := &MockFileProvider{}
	handler := NewFileHandler(mockProvider)

	fileName := "non-existent-file.jpg"
	expectedError := errors.New("file not found")

	mockProvider.On("GetPresignedURL", fileName).Return("", expectedError)

	// Act
	result, err := handler.FindFile(fileName)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, result)
	mockProvider.AssertExpectations(t)
}

func TestFileHandler_FindFile_EmptyFileName(t *testing.T) {
	// Arrange
	mockProvider := &MockFileProvider{}
	handler := NewFileHandler(mockProvider)

	fileName := ""
	expectedError := errors.New("invalid file name")

	mockProvider.On("GetPresignedURL", fileName).Return("", expectedError)

	// Act
	result, err := handler.FindFile(fileName)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, result)
	mockProvider.AssertExpectations(t)
}

func TestFileHandler_FindFile_DifferentFileTypes(t *testing.T) {
	// Arrange
	mockProvider := &MockFileProvider{}
	handler := NewFileHandler(mockProvider)

	testCases := []struct {
		fileName    string
		expectedURL string
	}{
		{"image.jpg", "https://example.com/image.jpg"},
		{"document.pdf", "https://example.com/document.pdf"},
		{"video.mp4", "https://example.com/video.mp4"},
		{"archive.zip", "https://example.com/archive.zip"},
	}

	for _, tc := range testCases {
		mockProvider.On("GetPresignedURL", tc.fileName).Return(tc.expectedURL, nil)

		// Act
		result, err := handler.FindFile(tc.fileName)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, tc.expectedURL, result)
	}

	mockProvider.AssertExpectations(t)
}

func TestFileHandler_FindFile_NetworkError(t *testing.T) {
	// Arrange
	mockProvider := &MockFileProvider{}
	handler := NewFileHandler(mockProvider)

	fileName := "test-file.jpg"
	networkError := errors.New("network timeout")

	mockProvider.On("GetPresignedURL", fileName).Return("", networkError)

	// Act
	result, err := handler.FindFile(fileName)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, networkError, err)
	assert.Empty(t, result)
	mockProvider.AssertExpectations(t)
}

func TestFileHandler_FindFile_PermissionError(t *testing.T) {
	// Arrange
	mockProvider := &MockFileProvider{}
	handler := NewFileHandler(mockProvider)

	fileName := "restricted-file.jpg"
	permissionError := errors.New("access denied")

	mockProvider.On("GetPresignedURL", fileName).Return("", permissionError)

	// Act
	result, err := handler.FindFile(fileName)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, permissionError, err)
	assert.Empty(t, result)
	mockProvider.AssertExpectations(t)
}