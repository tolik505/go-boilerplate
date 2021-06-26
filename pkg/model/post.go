package model

import (
	"github.com/pkg/errors"
	"gorm.io/datatypes"
)

type Post struct {
	UUID     string          `gorm:"primaryKey" json:"uuid"`
	Content  string          `json:"content"`
	Metadata *datatypes.JSON `json:"metadata"`
}

func (p *Post) MetadataString() *string {
	if p.Metadata == nil {
		return nil
	}

	s := p.Metadata.String()

	return &s
}

func (p *Post) SetMetadataFromString(m *string) error {
	if m == nil {
		return nil
	}

	j := &datatypes.JSON{}

	if err := j.UnmarshalJSON([]byte(*m)); err != nil {
		return errors.Wrap(err, "Couldn't unmarshal json")
	}

	p.Metadata = j

	return nil
}
