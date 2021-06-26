// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package integration

import (
	"github.com/google/wire"
	"goboilerplate/pkg/grpcapp"
	"goboilerplate/pkg/grpcapp/server"
	"goboilerplate/pkg/service/postservice"
	"goboilerplate/pkg/service/postservice/mocks"
	"goboilerplate/pkg/storage"
	"gorm.io/gorm"
)

func InitializeGRPCApp(DB *gorm.DB, us *mocks.UuidService) (*grpcapp.App, error) {
	wire.Build(
		postservice.NewPostService,
		storage.NewPostRepo,
		server.NewPostsServer,
		grpcapp.Initialize,
		wire.Bind(new(server.PostService), new(*postservice.PostService)),
		wire.Bind(new(postservice.PostRepo), new(*storage.PostRepo)),
		wire.Bind(new(postservice.UUIDService), new(*mocks.UuidService)),
	)

	return &grpcapp.App{}, nil
}
