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
	categoryRepo := new(mocks.CategoryRepositoryMock)

	usecase := NewAdminQuestionUseCase(questionRepo, categoryRepo)

	questionRepo.On("Create", mock.AnythingOfType("*entities.Question")).Return(&entities.Question{}, nil)

	questionPayload, err := entities.NewQuestion("siapa presiden pertama indonesia??", "")
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
		additionalOptions := []entities.Option{
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

		err = questionPayload.AddOption(additionalOptions...)
		assert.Error(t, err)
		assert.Len(t, options, 3)

		assert.True(t, len(questionPayload.Options) < 5) // 5 is options maximum number
	})

	t.Run("should return error if options exceed the maximum when calling CreateQuestion", func(t *testing.T) {
		// reset state
		questionPayload, err := entities.NewQuestion("siapa penemu listrik?", "")
		assert.NoError(t, err)

		options := []entities.Option{
			{Content: "Einstein", IsCorrect: false},
			{Content: "Tesla", IsCorrect: false},
			{Content: "Newton", IsCorrect: false},
			{Content: "Faraday", IsCorrect: false},
			{Content: "Galvani", IsCorrect: false},
			{Content: "Franklin", IsCorrect: true}, // kelebihan
		}

		err = questionPayload.AddOption(options...)
		assert.Error(t, err, "AddOption seharusnya error karena melebihi 5")

		// tetap paksa kirim ke usecase
		err = usecase.CreateQuestion(questionPayload)
		assert.Error(t, err, "CreateQuestion seharusnya error karena jumlah opsi melebihi batas")
	})
}

func TestAdminQuestionUseCase_ListQuestions(t *testing.T) {
	questionRepo := new(mocks.QuestionRepositoryMock)
	categoryRepo := new(mocks.CategoryRepositoryMock)
	usecase := NewAdminQuestionUseCase(questionRepo, categoryRepo)

	questionRepo.On("PaginateQuestions", mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return([]entities.Question{}, int64(10), nil)

	_, total, err := usecase.GetListQuestions(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), total)
}

func TestAdminQuestionUseCase_DeleteQuestion(t *testing.T) {
	questionRepo := new(mocks.QuestionRepositoryMock)
	categoryRepo := new(mocks.CategoryRepositoryMock)
	usecase := NewAdminQuestionUseCase(questionRepo, categoryRepo)

	questionRepo.On("DeleteById", mock.AnythingOfType("string")).Return(nil)

	err := usecase.DeleteQuestion("this-is-fake-question-id")
	assert.NoError(t, err)
}

func TestAdminQuestionUseCase_GetQuestionDetail(t *testing.T) {
	questionRepo := new(mocks.QuestionRepositoryMock)
	categoryRepo := new(mocks.CategoryRepositoryMock)
	usecase := NewAdminQuestionUseCase(questionRepo, categoryRepo)

	questionRepo.On("FindById", mock.AnythingOfType("string")).Return(&entities.Question{}, nil)

	question, err := usecase.GetQuestionDetail("this-is-fake-question-id")
	assert.NoError(t, err)
	assert.NotNil(t, question)
}
