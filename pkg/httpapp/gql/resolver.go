package gql

//go:generate go run github.com/99designs/gqlgen

import (
	"goboilerplate/pkg/model"
)

type PostService interface {
	Create(p model.CreatePostInput) (string, error)
	Update(UUID string, p model.UpdatePostInput) error
	Delete(UUID string) error
	Get(UUID string) (*model.Post, error)
	GetAll() ([]model.Post, error)
}

func NewResolver(
	postService PostService,
) *Resolver {
	return &Resolver{
		postService,
	}
}

type Resolver struct {
	postService PostService
}
