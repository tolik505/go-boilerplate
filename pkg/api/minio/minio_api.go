package minio

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"io"
)

type Client interface {
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64,
		opts minio.PutObjectOptions) (info minio.UploadInfo, err error)
}

type Config struct {
	Host   string
	Key    string
	Secret string
}

type Bucket string

type API struct {
	client Client
	bucket Bucket
}

func NewMinioAPI(client Client, bucket Bucket) *API {
	return &API{client, bucket}
}

func (s *API) Put(objectName string, reader io.Reader, contentType string) error {
	var buf bytes.Buffer
	size, err := io.Copy(&buf, reader)
	if err != nil {
		return err
	}
	_, err = s.client.PutObject(
		context.Background(),
		string(s.bucket),
		objectName,
		reader,
		size,
		minio.PutObjectOptions{ContentType: contentType},
	)

	return err
}
