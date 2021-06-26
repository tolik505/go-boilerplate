package server

import (
	"context"
	"goboilerplate/pkg/grpcapp/pb"
	"goboilerplate/pkg/model"
)

type PostService interface {
	GetAll() ([]model.Post, error)
}

func NewPostsServer(service PostService) pb.PostsServiceServer {
	return &postsServer{service}
}

type postsServer struct {
	service PostService
}

func (s *postsServer) GetAllPosts(ctx context.Context, request *pb.GetAllPostsRequest) (*pb.GetAllPostsResponse, error) {
	p, err := s.service.GetAll()
	if err != nil {
		return nil, err
	}
	var rp []*pb.Post
	// map internal type to pb type
	for i := range p {
		pbp := &pb.Post{
			Uuid:    p[i].UUID,
			Content: p[i].Content,
		}
		rp = append(rp, pbp)
	}
	resp := &pb.GetAllPostsResponse{
		Posts: rp,
	}

	return resp, nil
}
