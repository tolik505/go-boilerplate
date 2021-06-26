package storage

import (
	"github.com/pkg/errors"
	"goboilerplate/pkg/model"
	"gorm.io/gorm"
)

func NewPostRepo(db *gorm.DB) *PostRepo {
	return &PostRepo{db}
}

type PostRepo struct {
	db *gorm.DB
}

func (r *PostRepo) Create(p model.Post) error {
	return r.db.Create(&p).Error
}

func (r *PostRepo) Update(p model.Post) error {
	return r.db.Model(&p).Where("uuid = ?", p.UUID).Updates(&p).Error
}

func (r *PostRepo) Get(UUID string) (*model.Post, error) {
	p := &model.Post{}
	err := r.db.Where("uuid = ?", UUID).Take(p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *PostRepo) GetAll() ([]model.Post, error) {
	var p []model.Post

	if err := r.db.Find(&p).Error; err != nil {
		return nil, err
	}

	return p, nil
}

func (r *PostRepo) Delete(UUID string) error {
	return r.db.Where("uuid = ?", UUID).Delete(model.Post{}).Error
}
