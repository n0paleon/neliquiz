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

	t.Run("should return all categories when query is empty", func(t *testing.T) {
		categoryRepo.On("FindAll").Return([]entities.Category{
			{ID: "1", Name: "Math"},
			{ID: "2", Name: "Science"},
		}, nil).Once()

		categories, err := usecase.GetListCategories("", 0)
		assert.NoError(t, err)
		assert.Len(t, categories, 2)
		assert.Equal(t, "Math", categories[0].Name)
		assert.Equal(t, "Science", categories[1].Name)

		categoryRepo.AssertExpectations(t)
	})

	t.Run("should return search result when query and limit are provided", func(t *testing.T) {
		query := "ma"
		limit := 10
		categoryRepo.On("SearchCategoriesByName", query, limit).Return([]entities.Category{
			{ID: "1", Name: "Matematika"},
		}, nil).Once()

		categories, err := usecase.GetListCategories(query, limit)
		assert.NoError(t, err)
		assert.Len(t, categories, 1)
		assert.Equal(t, "Matematika", categories[0].Name)

		categoryRepo.AssertExpectations(t)
	})

	t.Run("should propagate error from repository", func(t *testing.T) {
		query := "fi"
		limit := 5
		mockErr := assert.AnError

		categoryRepo.On("SearchCategoriesByName", query, limit).Return([]entities.Category(nil), mockErr).Once()

		categories, err := usecase.GetListCategories(query, limit)
		assert.Nil(t, categories)
		assert.Equal(t, mockErr, err)

		categoryRepo.AssertExpectations(t)
	})
}
