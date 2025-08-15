package mocks

import (
	"NeliQuiz/internal/features/category/domain"
	"github.com/stretchr/testify/mock"
)

type CategoryUseCaseMock struct {
	mock.Mock
}

func (m *CategoryUseCaseMock) GetListCategories(query string, limit int) ([]domain.Category, error) {
	args := m.Called(query, limit)
	return args.Get(0).([]domain.Category), args.Error(1)
}

func (m *CategoryUseCaseMock) FindOrCreateBatch(categories []domain.Category) ([]domain.Category, error) {
	args := m.Called(categories)
	return args.Get(0).([]domain.Category), args.Error(1)
}

func (m *CategoryUseCaseMock) FindCategoryByName(name string) (*domain.Category, error) {
	args := m.Called(name)
	return args.Get(0).(*domain.Category), args.Error(1)
}
