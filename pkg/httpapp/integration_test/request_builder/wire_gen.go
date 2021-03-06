// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package request_builder

import (
	"goboilerplate/pkg/api/minio"
	mocks2 "goboilerplate/pkg/api/minio/mocks"
	"goboilerplate/pkg/httpapp"
	"goboilerplate/pkg/httpapp/gql"
	mocks3 "goboilerplate/pkg/httpapp/gql/mocks"
	mocks4 "goboilerplate/pkg/httpapp/mocks"
	"goboilerplate/pkg/service/fileservice"
	"goboilerplate/pkg/service/postservice"
	"goboilerplate/pkg/service/postservice/mocks"
	"goboilerplate/pkg/storage"
	"gorm.io/gorm"
)

// Injectors from wire.go:

func InitializeHTTPAppIntegration(DB *gorm.DB, us *mocks.UuidService, mcl *mocks2.Client, mc minio.Config, buck minio.Bucket, port httpapp.Port) (*httpapp.App, error) {
	postRepo := storage.NewPostRepo(DB)
	postService := postservice.NewPostService(postRepo, us)
	resolver := gql.NewResolver(postService)
	api := minio.NewMinioAPI(mcl, buck)
	fileService := fileservice.NewFileService(api, us)
	handler := httpapp.NewHandler(fileService)
	app := httpapp.Initialize(resolver, handler, port)
	return app, nil
}

func InitializeHTTPAppFunctional(ps *mocks3.PostService, fs *mocks4.FileService, port httpapp.Port) (*httpapp.App, error) {
	resolver := gql.NewResolver(ps)
	handler := httpapp.NewHandler(fs)
	app := httpapp.Initialize(resolver, handler, port)
	return app, nil
}
