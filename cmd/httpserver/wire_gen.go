// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"goboilerplate/pkg/api/minio"
	"goboilerplate/pkg/httpapp"
	"goboilerplate/pkg/httpapp/gql"
	"goboilerplate/pkg/service/fileservice"
	"goboilerplate/pkg/service/postservice"
	"goboilerplate/pkg/service/uuidservice"
	"goboilerplate/pkg/storage"
	"goboilerplate/pkg/storage/db"
)

// Injectors from wire.go:

func InitializeHTTPApp(port httpapp.Port, dbc db.Config, mc minio.Config, bucket minio.Bucket) (*httpapp.App, error) {
	gormDB, err := db.InitDB(dbc)
	if err != nil {
		return nil, err
	}
	postRepo := storage.NewPostRepo(gormDB)
	uuidService := uuidservice.NewUUIDService()
	postService := postservice.NewPostService(postRepo, uuidService)
	resolver := gql.NewResolver(postService)
	client, err := minio.NewClient(mc)
	if err != nil {
		return nil, err
	}
	api := minio.NewMinioAPI(client, bucket)
	fileService := fileservice.NewFileService(api, uuidService)
	handler := httpapp.NewHandler(fileService)
	app := httpapp.Initialize(resolver, handler, port)
	return app, nil
}
