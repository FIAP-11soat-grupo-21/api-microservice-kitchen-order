package interfaces

type IFileProvider interface {
	UploadFile(fileName string, fileContent []byte) error
	DeleteFile(fileName string) error
	GetPresignedURL(fileName string) (string, error)
}
