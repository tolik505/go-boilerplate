package postservice

import (
	"github.com/pkg/errors"
	"goboilerplate/pkg/model"
)

type PostRepo interface {
	Create(p model.Post) error
	Update(p model.Post) error
	Delete(UUID string) error
	Get(UUID string) (*model.Post, error)
	GetAll() ([]model.Post, error)
}

type UUIDService interface {
	Generate() string
}

type PostService struct {
	repo  PostRepo
	uuidS UUIDService
}

func NewPostService(
	repo PostRepo,
	uuidS UUIDService,
) *PostService {
	return &PostService{repo, uuidS}
}

func (s *PostService) Create(i model.CreatePostInput) (string, error) {
	UUID := s.uuidS.Generate()
	p := model.Post{
		UUID:    UUID,
		Content: i.Content,
	}
	if err := p.SetMetadataFromString(i.Metadata); err != nil {
		return "", err
	}
	if err := s.repo.Create(p); err != nil {
		return "", errors.Wrap(err, "Couldn't create a post")
	}

	return UUID, nil
}

func (s *PostService) Update(UUID string, i model.UpdatePostInput) error {
	p := model.Post{
		UUID:    UUID,
		Content: i.Content,
	}
	if err := p.SetMetadataFromString(i.Metadata); err != nil {
		return err
	}
	if err := s.repo.Update(p); err != nil {
		return errors.Wrapf(err, "Couldn't update the post [%s]", p.UUID)
	}

	return nil
}

func (s *PostService) Get(UUID string) (*model.Post, error) {
	p, err := s.repo.Get(UUID)
	if err != nil {
		return nil, errors.Wrapf(err, "Couldn't get a post [%s]", UUID)
	}

	return p, nil
}

func (s *PostService) GetAll() ([]model.Post, error) {
	p, err := s.repo.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't all posts")
	}

	return p, nil
}

func (s *PostService) Delete(UUID string) error {
	if err := s.repo.Delete(UUID); err != nil {
		return errors.Wrapf(err, "Couldn't delete the post [%s]", UUID)
	}

	return nil
}
