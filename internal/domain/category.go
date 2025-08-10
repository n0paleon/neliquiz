package domain

import "NeliQuiz/internal/domain/entities"

type CategoryRepository interface {
	FindOrCreateBatch(categories []entities.Category) ([]entities.Category, error)
	FindOrCreateCategoryByName(name string) (*entities.Category, error)
	FindCategoryByName(name string) (*entities.Category, error)
	FindAll() ([]entities.Category, error)
}

type AdminCategoryUseCase interface {
	GetListCategories() ([]entities.Category, error)
}
