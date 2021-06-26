package fileservice

import (
	"github.com/pkg/errors"
	"goboilerplate/pkg/service/postservice"
	"io"
)

type FileStorage interface {
	Put(objectName string, reader io.Reader, contentType string) error
}

type FileService struct {
	fileStorage FileStorage
	uuidS       postservice.UUIDService
}

func NewFileService(
	fs FileStorage,
	uuidS postservice.UUIDService,
) *FileService {
	return &FileService{fs, uuidS}
}

func (f *FileService) Upload(src io.Reader, contentType string) (string, error) {
	UUID := f.uuidS.Generate()
	err := f.fileStorage.Put(UUID, src, contentType)
	if err != nil {
		return "", errors.Wrap(err, "Couldn't upload file into file storage")
	}

	return UUID, nil
}
