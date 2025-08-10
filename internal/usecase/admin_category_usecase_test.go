package usecase

import (
	"NeliQuiz/internal/domain/entities"
	"NeliQuiz/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdminCategoryUseCase_GetListCategories(t *testing.T) {
	categoryRepo := new(mocks.CategoryRepositoryMock)
	usecase := NewAdminCategoryUseCase(categoryRepo)

	categoryRepo.On("FindAll").Return([]entities.Category{}, nil)

	categories, err := usecase.GetListCategories()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(categories))
}
