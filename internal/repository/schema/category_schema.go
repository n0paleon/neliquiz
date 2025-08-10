package schema

import "NeliQuiz/internal/domain/entities"

type Category struct {
	Schema
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
		Schema: Schema{
			ID: e.ID,
		},
		Name: e.Name,
	}
}
