// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"goboilerplate/pkg/grpcapp"
	"goboilerplate/pkg/grpcapp/server"
	"goboilerplate/pkg/service/postservice"
	"goboilerplate/pkg/service/uuidservice"
	"goboilerplate/pkg/storage"
	"goboilerplate/pkg/storage/db"
)

// Injectors from wire.go:

func InitializeGRPCApp(dbc db.Config) (*grpcapp.App, error) {
	gormDB, err := db.InitDB(dbc)
	if err != nil {
		return nil, err
	}
	postRepo := storage.NewPostRepo(gormDB)
	uuidService := uuidservice.NewUUIDService()
	postService := postservice.NewPostService(postRepo, uuidService)
	postsServiceServer := server.NewPostsServer(postService)
	app := grpcapp.Initialize(postsServiceServer)
	return app, nil
}