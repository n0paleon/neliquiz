package domain

import "NeliQuiz/internal/domain/entities"

type CategoryRepository interface {
	FindOrCreateCategoryByName(name string) (*entities.Category, error)
}
