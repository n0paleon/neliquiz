package mocks

import (
	"NeliQuiz/internal/domain/entities"
	"github.com/stretchr/testify/mock"
)

type QuestionRepositoryMock struct {
	mock.Mock
}

func (m *QuestionRepositoryMock) Create(q *entities.Question) (*entities.Question, error) {
	args := m.Called(q)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Question), args.Error(1)
}

func (m *QuestionRepositoryMock) FindById(id string) (*entities.Question, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Question), args.Error(1)
}

func (m *QuestionRepositoryMock) DeleteById(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *QuestionRepositoryMock) GetRandom() (*entities.Question, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Question), args.Error(1)
}

func (m *QuestionRepositoryMock) PaginateQuestions(page, limit int) ([]entities.Question, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]entities.Question), args.Get(1).(int64), args.Error(2)
}
