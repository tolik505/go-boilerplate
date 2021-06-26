package server

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"goboilerplate/pkg/grpcapp/pb"
	"goboilerplate/pkg/grpcapp/server/mocks"
	"goboilerplate/pkg/model"
	"testing"
)

func Test_postsServer_GetAllPosts(t *testing.T) {
	ctx := context.Background()

	t.Run("It Returns Posts", func(t *testing.T) {
		srvc := &mocks.PostService{}
		psts := []model.Post{
			{
				UUID:    "uuid-1",
				Content: "my post",
			},
		}
		srvc.On("GetAll").Return(psts, nil)
		req := &pb.GetAllPostsRequest{}
		s := NewPostsServer(srvc)
		res, err := s.GetAllPosts(ctx, req)

		expected := &pb.GetAllPostsResponse{
			Posts: []*pb.Post{
				{
					Uuid:    "uuid-1",
					Content: "my post",
				},
			},
		}
		assert.Equal(t, expected, res)
		assert.Nil(t, err)
	})

	t.Run("It Returns Empty Posts", func(t *testing.T) {
		srvc := &mocks.PostService{}
		var psts []model.Post
		srvc.On("GetAll").Return(psts, nil)
		req := &pb.GetAllPostsRequest{}
		s := NewPostsServer(srvc)
		res, err := s.GetAllPosts(ctx, req)

		expected := &pb.GetAllPostsResponse{
			Posts: []*pb.Post(nil),
		}
		assert.Equal(t, expected, res)
		assert.Nil(t, err)
	})

	t.Run("It Handles Service Error", func(t *testing.T) {
		srvc := &mocks.PostService{}
		var psts []model.Post
		srvc.On("GetAll").Return(psts, errors.New("error"))
		req := &pb.GetAllPostsRequest{}
		s := NewPostsServer(srvc)
		res, err := s.GetAllPosts(ctx, req)

		assert.Nil(t, res)
		assert.EqualError(t, err, "error")
	})
}
