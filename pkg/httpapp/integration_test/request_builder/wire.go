// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package request_builder

import (
	"github.com/google/wire"
	"goboilerplate/pkg/api/minio"
	mocks2 "goboilerplate/pkg/api/minio/mocks"
	"goboilerplate/pkg/httpapp"
	"goboilerplate/pkg/httpapp/gql"
	"goboilerplate/pkg/httpapp/gql/generated"
	mocks3 "goboilerplate/pkg/httpapp/gql/mocks"
	mocks4 "goboilerplate/pkg/httpapp/mocks"
	"goboilerplate/pkg/service/fileservice"
	"goboilerplate/pkg/service/postservice"
	"goboilerplate/pkg/service/postservice/mocks"
	"goboilerplate/pkg/storage"
	"gorm.io/gorm"
)

func InitializeHTTPAppIntegration(
	DB *gorm.DB,
	us *mocks.UuidService,
	mcl *mocks2.Client,
	mc minio.Config,
	buck minio.Bucket,
	port httpapp.Port,
) (*httpapp.App, error) {

	wire.Build(
		minio.NewMinioAPI,
		httpapp.NewHandler,
		gql.NewResolver,
		httpapp.Initialize,
		postservice.NewPostService,
		storage.NewPostRepo,
		fileservice.NewFileService,
		wire.Bind(new(generated.ResolverRoot), new(*gql.Resolver)),
		wire.Bind(new(gql.PostService), new(*postservice.PostService)),
		wire.Bind(new(httpapp.FileService), new(*fileservice.FileService)),
		wire.Bind(new(fileservice.FileStorage), new(*minio.API)),
		wire.Bind(new(postservice.PostRepo), new(*storage.PostRepo)),
		wire.Bind(new(postservice.UUIDService), new(*mocks.UuidService)),
		wire.Bind(new(minio.Client), new(*mocks2.Client)),
	)

	return &httpapp.App{}, nil
}

func InitializeHTTPAppFunctional(
	ps *mocks3.PostService,
	fs *mocks4.FileService,
	port httpapp.Port,
) (*httpapp.App, error) {

	wire.Build(
		httpapp.NewHandler,
		gql.NewResolver,
		httpapp.Initialize,
		wire.Bind(new(generated.ResolverRoot), new(*gql.Resolver)),
		wire.Bind(new(gql.PostService), new(*mocks3.PostService)),
		wire.Bind(new(httpapp.FileService), new(*mocks4.FileService)),
	)

	return &httpapp.App{}, nil
}
