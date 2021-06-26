package uuidservice

import (
	"github.com/google/uuid"
)

func NewUUIDService() *UUIDService {
	return &UUIDService{}
}

type UUIDService struct{}

func (s *UUIDService) Generate() string {
	return uuid.New().String()
}
