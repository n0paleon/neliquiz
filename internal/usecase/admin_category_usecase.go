package usecase

import (
	"NeliQuiz/internal/domain"
	"NeliQuiz/internal/domain/entities"
)

type AdminCategoryUseCase struct {
	categoryRepo domain.CategoryRepository
}

func (u *AdminCategoryUseCase) GetListCategories(query string, limit int) ([]entities.Category, error) {
	if query != "" && limit != 0 {
		return u.categoryRepo.SearchCategoriesByName(query, limit)
	}

	return u.categoryRepo.FindAll()
}

func NewAdminCategoryUseCase(categoryRepo domain.CategoryRepository) *AdminCategoryUseCase {
	return &AdminCategoryUseCase{categoryRepo}
}
