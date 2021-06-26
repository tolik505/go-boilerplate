package minio

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewClient(conf Config) (*minio.Client, error) {
	client, err := minio.New(conf.Host, &minio.Options{
		Creds: credentials.NewStaticV4(conf.Key, conf.Secret, ""),
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}
