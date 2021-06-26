package postservice

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"goboilerplate/pkg/model"
	"goboilerplate/pkg/service/postservice/mocks"
	"testing"
)

func TestPostService_Create(t *testing.T) {
	t.Run("It Successfully Creates A Post", func(t *testing.T) {
		us := &mocks.UuidService{}
		repo := &mocks.PostRepo{}
		us.On("Generate").Return("uuid-1").Once()
		pc := model.Post{
			UUID:    "uuid-1",
			Content: "my post",
		}
		repo.On("Create", pc).Return(nil).Once()

		s := NewPostService(repo, us)
		p := model.CreatePostInput{
			Content: "my post",
		}
		uuid, err := s.Create(p)

		mock.AssertExpectationsForObjects(t, us, repo)
		assert.Equal(t, "uuid-1", uuid)
		assert.Nil(t, err)
	})

	t.Run("It Handles Repo Error", func(t *testing.T) {
		us := &mocks.UuidService{}
		repo := &mocks.PostRepo{}
		us.On("Generate").Return("uuid-1")
		pc := model.Post{
			UUID:    "uuid-1",
			Content: "my post",
		}
		repo.On("Create", pc).Return(errors.New("error"))

		s := NewPostService(repo, us)
		p := model.CreatePostInput{
			Content: "my post",
		}
		uuid, err := s.Create(p)

		assert.Empty(t, uuid)
		assert.EqualError(t, err, "Couldn't create a post: error")
	})
}
