package mocks

import (
	"NeliQuiz/internal/domain/entities"
	"github.com/stretchr/testify/mock"
)

type CategoryRepositoryMock struct {
	mock.Mock
}

func (m *CategoryRepositoryMock) FindOrCreateCategoryByName(name string) (*entities.Category, error) {
	args := m.Called(name)
	return args.Get(0).(*entities.Category), args.Error(1)
}
