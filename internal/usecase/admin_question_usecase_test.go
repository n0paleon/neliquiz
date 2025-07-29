package usecase

import (
	"NeliQuiz/internal/domain/entities"
	"NeliQuiz/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAdminQuestionUsecase_CreateQuestion(t *testing.T) {
	questionRepo := new(mocks.QuestionRepositoryMock)

	usecase := NewAdminQuestionUseCase(questionRepo)

	questionRepo.On("Create", mock.AnythingOfType("*entities.Question")).Return(&entities.Question{}, nil)

	questionPayload, err := entities.NewQuestion("siapa presiden pertama indonesia??")
	assert.NoError(t, err)

	options := []entities.Option{
		entities.Option{
			Content:   "jokowi",
			IsCorrect: false,
		},
		entities.Option{
			Content:   "prabowo",
			IsCorrect: false,
		},
		entities.Option{
			Content:   "soekarno",
			IsCorrect: true,
		},
	}

	t.Run("should create question correctly", func(t *testing.T) {
		err = questionPayload.AddOption(options...)
		assert.NoError(t, err)

		err := usecase.CreateQuestion(questionPayload)
		assert.NoError(t, err)
	})

	t.Run("should error when options is reach the maximum", func(t *testing.T) {
		options := []entities.Option{
			entities.Option{
				Content:   "mohammad hatta",
				IsCorrect: false,
			},
			entities.Option{
				Content:   "tan malaka",
				IsCorrect: false,
			},
			entities.Option{
				Content:   "sby",
				IsCorrect: true,
			},
		}

		err = questionPayload.AddOption(options...)
		assert.Error(t, err)

		assert.True(t, len(questionPayload.Options) < 5) // 5 is options maximum number
	})
}

func TestAdminQuestionUseCase_ListQuestions(t *testing.T) {
	questionRepo := new(mocks.QuestionRepositoryMock)
	usecase := NewAdminQuestionUseCase(questionRepo)

	questionRepo.On("PaginateQuestions", mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return([]entities.Question{}, int64(10), nil)

	_, total, err := usecase.GetListQuestions(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), total)
}

func TestAdminQuestionUseCase_DeleteQuestion(t *testing.T) {
	questionRepo := new(mocks.QuestionRepositoryMock)
	usecase := NewAdminQuestionUseCase(questionRepo)

	questionRepo.On("DeleteById", mock.AnythingOfType("string")).Return(nil)

	err := usecase.DeleteQuestion("this-is-fake-question-id")
	assert.NoError(t, err)
}
