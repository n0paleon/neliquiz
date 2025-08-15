package mocks

import (
	"NeliQuiz/internal/features/category/domain"
	"github.com/stretchr/testify/mock"
)

type CategoryRepositoryMock struct {
	mock.Mock
}

func (m *CategoryRepositoryMock) FindOrCreateBatch(categories []domain.Category) ([]domain.Category, error) {
	args := m.Called(categories)
	return args.Get(0).([]domain.Category), args.Error(1)
}

func (m *CategoryRepositoryMock) FindOrCreateCategoryByName(name string) (*domain.Category, error) {
	args := m.Called(name)
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *CategoryRepositoryMock) FindCategoryByName(name string) (*domain.Category, error) {
	args := m.Called(name)
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *CategoryRepositoryMock) FindAll() ([]domain.Category, error) {
	args := m.Called()
	return args.Get(0).([]domain.Category), args.Error(1)
}

func (m *CategoryRepositoryMock) SearchCategoriesByName(query string, limit int) ([]domain.Category, error) {
	args := m.Called(query, limit)
	return args.Get(0).([]domain.Category), args.Error(1)
}
