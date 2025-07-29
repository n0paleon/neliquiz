package usecase_test

import (
	"NeliQuiz/internal/domain/entities"
	"NeliQuiz/internal/repository/mocks"
	"NeliQuiz/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRandomQuestion(t *testing.T) {
	mockRepo := new(mocks.QuestionRepositoryMock)
	useCase := usecase.NewUserQuestionUseCase(mockRepo)

	expectedQuestion := &entities.Question{
		ID:      "123",
		Content: "Apa ibu kota Indonesia?",
		Options: []entities.Option{
			entities.Option{
				ID:        "1",
				Content:   "Jakarta",
				IsCorrect: true,
			},
			entities.Option{
				ID:        "2",
				Content:   "Bandung",
				IsCorrect: false,
			},
			entities.Option{
				ID:        "3",
				Content:   "Surabaya",
				IsCorrect: false,
			},
			entities.Option{
				ID:        "4",
				Content:   "Yogyakarta",
				IsCorrect: false,
			},
		},
	}

	mockRepo.On("GetRandom").Return(expectedQuestion, nil)

	result, err := useCase.GetRandomQuestion()

	assert.NoError(t, err)
	assert.Equal(t, expectedQuestion, result)
	mockRepo.AssertExpectations(t)
}

func TestGetRandomQuestion_Error(t *testing.T) {
	mockRepo := new(mocks.QuestionRepositoryMock)
	useCase := usecase.NewUserQuestionUseCase(mockRepo)

	mockRepo.On("GetRandom").Return(nil, assert.AnError)

	result, err := useCase.GetRandomQuestion()

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestCheckAnswer_Correct(t *testing.T) {
	mockRepo := new(mocks.QuestionRepositoryMock)
	useCase := usecase.NewUserQuestionUseCase(mockRepo)

	option := entities.Option{
		ID:        "option1",
		Content:   "Jakarta",
		IsCorrect: true,
	}

	question := &entities.Question{
		ID:      "question1",
		Content: "Ibu kota Indonesia?",
		Options: []entities.Option{option},
	}

	mockRepo.On("FindById", "question1").Return(question, nil)

	// Call usecase
	isCorrect, selected, err := useCase.CheckAnswer("question1", "option1")

	assert.NoError(t, err)
	assert.True(t, isCorrect)
	assert.Equal(t, &option, selected)
	mockRepo.AssertExpectations(t)
}

func TestCheckAnswer_InvalidQuestion(t *testing.T) {
	mockRepo := new(mocks.QuestionRepositoryMock)
	useCase := usecase.NewUserQuestionUseCase(mockRepo)

	mockRepo.On("FindById", "invalid-id").Return(nil, assert.AnError)

	isCorrect, selected, err := useCase.CheckAnswer("invalid-id", "option1")

	assert.Error(t, err)
	assert.False(t, isCorrect)
	assert.Nil(t, selected)
	mockRepo.AssertExpectations(t)
}

func TestCheckAnswer_WrongOption(t *testing.T) {
	mockRepo := new(mocks.QuestionRepositoryMock)
	useCase := usecase.NewUserQuestionUseCase(mockRepo)

	option := entities.Option{
		ID:        "option1",
		Content:   "Surabaya",
		IsCorrect: false,
	}

	question := &entities.Question{
		ID:      "question1",
		Content: "Ibu kota Indonesia?",
		Options: []entities.Option{option},
	}

	mockRepo.On("FindById", "question1").Return(question, nil)

	isCorrect, selected, err := useCase.CheckAnswer("question1", "option1")

	assert.NoError(t, err)
	assert.False(t, isCorrect)
	assert.Equal(t, &option, selected)
	mockRepo.AssertExpectations(t)
}
