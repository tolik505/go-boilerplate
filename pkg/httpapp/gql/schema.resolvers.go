package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"goboilerplate/pkg/httpapp/gql/generated"
	"goboilerplate/pkg/model"
)

func (r *mutationResolver) CreatePost(ctx context.Context, data model.CreatePostInput) (string, error) {
	UUID, err := r.postService.Create(data)
	if err != nil {
		return "", err
	}

	return UUID, nil
}

func (r *mutationResolver) UpdatePost(ctx context.Context, uuid string, data model.UpdatePostInput) (bool, error) {
	err := r.postService.Update(uuid, data)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeletePost(ctx context.Context, uuid string) (bool, error) {
	err := r.postService.Delete(uuid)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) Posts(ctx context.Context) ([]model.Post, error) {
	p, err := r.postService.GetAll()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *queryResolver) Post(ctx context.Context, uuid string) (*model.Post, error) {
	p, err := r.postService.Get(uuid)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
