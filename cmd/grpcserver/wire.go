// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/google/wire"
	"goboilerplate/pkg/grpcapp"
	"goboilerplate/pkg/grpcapp/server"
	"goboilerplate/pkg/service/postservice"
	"goboilerplate/pkg/service/uuidservice"
	"goboilerplate/pkg/storage"
	"goboilerplate/pkg/storage/db"
)

func InitializeGRPCApp(dbc db.Config) (*grpcapp.App, error) {
	wire.Build(
		db.InitDB,
		postservice.NewPostService,
		storage.NewPostRepo,
		uuidservice.NewUUIDService,
		server.NewPostsServer,
		grpcapp.Initialize,
		wire.Bind(new(server.PostService), new(*postservice.PostService)),
		wire.Bind(new(postservice.PostRepo), new(*storage.PostRepo)),
		wire.Bind(new(postservice.UUIDService), new(*uuidservice.UUIDService)),
	)

	return &grpcapp.App{}, nil
}
