package schema

import (
	"NeliQuiz/pkg/utils"
	"gorm.io/gorm"
)

type Schema struct {
	ID string `gorm:"primaryKey;type:char(26)"`
}

func (s *Schema) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = utils.GenerateULID()
	}

	return
}
