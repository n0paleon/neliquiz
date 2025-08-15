package repository

import (
	"NeliQuiz/internal/features/category/domain"
	"NeliQuiz/internal/shared/strutil"
	"gorm.io/gorm"
)

type CategorySchema struct {
	ID   string `gorm:"primaryKey;type:char(26)"`
	Name string `gorm:"unique;not null"`
}

func (s *CategorySchema) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = strutil.GenerateULID()
	}

	return
}

func (s *CategorySchema) TableName() string {
	return "categories"
}

func (s *CategorySchema) ToEntity() *domain.Category {
	return &domain.Category{
		ID:   s.ID,
		Name: s.Name,
	}
}

func ToCategorySchema(e *domain.Category) *CategorySchema {
	return &CategorySchema{
		ID:   e.ID,
		Name: e.Name,
	}
}
