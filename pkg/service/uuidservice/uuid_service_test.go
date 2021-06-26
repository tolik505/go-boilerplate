package uuidservice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_uuidService_Generate(t *testing.T) {
	s := NewUUIDService()
	uuid := s.Generate()
	assert.IsType(t, uuid, "string")
	assert.Len(t, uuid, 36)
}
