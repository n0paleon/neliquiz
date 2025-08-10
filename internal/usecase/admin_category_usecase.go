package usecase

import (
	"NeliQuiz/internal/domain"
	"NeliQuiz/internal/domain/entities"
)

type AdminCategoryUseCase struct {
	categoryRepo domain.CategoryRepository
}

func (u *AdminCategoryUseCase) GetListCategories() ([]entities.Category, error) {
	return u.categoryRepo.FindAll()
}

func NewAdminCategoryUseCase(categoryRepo domain.CategoryRepository) *AdminCategoryUseCase {
	return &AdminCategoryUseCase{categoryRepo}
}
