package usecase

import (
	"NeliQuiz/internal/features/category/domain"
	"NeliQuiz/internal/shared/errorx"
)

type CategoryUseCase struct {
	categoryRepo domain.CategoryRepository
}

func (u *CategoryUseCase) GetListCategories(query string, limit int) ([]domain.Category, error) {
	if query != "" && limit != 0 {
		return u.categoryRepo.SearchCategoriesByName(query, limit)
	}

	return u.categoryRepo.FindAll()
}

func (u *CategoryUseCase) FindOrCreateBatch(categories []domain.Category) ([]domain.Category, error) {
	if len(categories) == 0 {
		return categories, nil
	}

	for _, category := range categories {
		if err := category.Validate(); err != nil {
			return categories, errorx.BadRequest(err.Error())
		}
	}

	return u.categoryRepo.FindOrCreateBatch(categories)
}

func (u *CategoryUseCase) FindCategoryByName(name string) (*domain.Category, error) {
	if name == "" {
		return nil, errorx.BadRequest("category name is required")
	}

	category := domain.Category{
		Name: name,
	}
	if err := category.Validate(); err != nil {
		return nil, errorx.BadRequest(err.Error())
	}

	return u.categoryRepo.FindCategoryByName(category.Name)
}

func NewCategoryUseCase(categoryRepo domain.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{categoryRepo}
}
