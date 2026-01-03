package factories

import (
	file_provider "tech_challenge/internal/shared/infra/file_provider"
	"tech_challenge/internal/shared/interfaces"
)

func NewFileProvider() interfaces.IFileProvider {
	return file_provider.NewS3FileProvider("")
}
