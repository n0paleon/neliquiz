package schema

import "NeliQuiz/internal/domain/entities"

type Category struct {
	ID   string `gorm:"type:uuid;primary_key;default:uuid_generate_v4();primary_key"`
	Name string `gorm:"unique;not null"`
}

func (s *Category) TableName() string {
	return "categories"
}

func (s *Category) ToEntity() *entities.Category {
	return &entities.Category{
		ID:   s.ID,
		Name: s.Name,
	}
}

func ToCategorySchema(e *entities.Category) *Category {
	return &Category{
		ID:   e.ID,
		Name: e.Name,
	}
}
