// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/google/wire"
	minio2 "github.com/minio/minio-go/v7"
	"goboilerplate/pkg/api/minio"
	"goboilerplate/pkg/httpapp"
	"goboilerplate/pkg/httpapp/gql"
	"goboilerplate/pkg/httpapp/gql/generated"
	"goboilerplate/pkg/service/fileservice"
	"goboilerplate/pkg/service/postservice"
	"goboilerplate/pkg/service/uuidservice"
	"goboilerplate/pkg/storage"
	"goboilerplate/pkg/storage/db"
)

func InitializeHTTPApp(
	port httpapp.Port,
	dbc db.Config,
	mc minio.Config,
	bucket minio.Bucket,
) (*httpapp.App, error) {
	wire.Build(
		db.InitDB,
		minio.NewClient,
		minio.NewMinioAPI,
		httpapp.NewHandler,
		gql.NewResolver,
		httpapp.Initialize,
		postservice.NewPostService,
		storage.NewPostRepo,
		fileservice.NewFileService,
		uuidservice.NewUUIDService,
		wire.Bind(new(generated.ResolverRoot), new(*gql.Resolver)),
		wire.Bind(new(gql.PostService), new(*postservice.PostService)),
		wire.Bind(new(httpapp.FileService), new(*fileservice.FileService)),
		wire.Bind(new(fileservice.FileStorage), new(*minio.API)),
		wire.Bind(new(postservice.PostRepo), new(*storage.PostRepo)),
		wire.Bind(new(postservice.UUIDService), new(*uuidservice.UUIDService)),
		wire.Bind(new(minio.Client), new(*minio2.Client)),
	)

	return &httpapp.App{}, nil
}
