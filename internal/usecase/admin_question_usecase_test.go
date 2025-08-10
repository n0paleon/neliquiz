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
	categoryRepo.On("FindOrCreateBatch", mock.AnythingOfType("[]entities.Category")).Return([]entities.Category{}, nil)

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

		_, err := usecase.CreateQuestion(questionPayload)
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
		_, err = usecase.CreateQuestion(questionPayload)
		assert.Error(t, err, "CreateQuestion seharusnya error karena jumlah opsi melebihi batas")
	})
}

func TestAdminQuestionUseCase_ListQuestions(t *testing.T) {
	questionRepo := new(mocks.QuestionRepositoryMock)
	categoryRepo := new(mocks.CategoryRepositoryMock)
	usecase := NewAdminQuestionUseCase(questionRepo, categoryRepo)

	questionRepo.On("PaginateQuestions", mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return([]entities.Question{}, int64(10), nil)

	_, total, err := usecase.GetListQuestions("", 1, 10, "", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(10), total)
}

func TestAdminQuestionUseCase_ListQuestionsWithCategory(t *testing.T) {
	categoryRepo := new(mocks.CategoryRepositoryMock)
	questionRepo := new(mocks.QuestionRepositoryMock)
	usecase := NewAdminQuestionUseCase(questionRepo, categoryRepo)

	categoryRepo.On("FindCategoryByName", mock.AnythingOfType("string")).Return(&entities.Category{}, nil)
	questionRepo.On("PaginateQuestionsByCategory", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return([]entities.Question{}, int64(10), nil)
	questionRepo.On("PaginateQuestions", mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return([]entities.Question{}, int64(10), nil)

	_, total, err := usecase.GetListQuestions("", 1, 10, "", "")
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

func TestAdminQuestionUseCase_UpdateQuestion(t *testing.T) {
	questionRepo := new(mocks.QuestionRepositoryMock)
	categoryRepo := new(mocks.CategoryRepositoryMock)
	usecase := NewAdminQuestionUseCase(questionRepo, categoryRepo)

	categoryRepo.On("FindOrCreateBatch", mock.AnythingOfType("[]entities.Category")).Return([]entities.Category{}, nil)

	// Data awal yang mau diupdate
	originalQuestion, err := entities.NewQuestion("Siapa presiden pertama Indonesia?", "")
	assert.NoError(t, err)
	originalQuestion.ID = "question-123"
	originalQuestion.AddOption(
		entities.Option{Content: "Jokowi", IsCorrect: false},
		entities.Option{Content: "Soekarno", IsCorrect: true},
	)

	// Data hasil update
	updatedQuestion := *originalQuestion
	updatedQuestion.Content = "Siapa presiden pertama RI?"
	updatedQuestion.Options[0].Content = "Prabowo"

	t.Run("should update question successfully", func(t *testing.T) {
		// Mock pemanggilan repo
		questionRepo.On("Update", mock.AnythingOfType("*entities.Question")).
			Return(&updatedQuestion, nil).
			Once()

		result, err := usecase.UpdateQuestion(originalQuestion)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Siapa presiden pertama RI?", result.Content)
		assert.Equal(t, "Prabowo", result.Options[0].Content)

		questionRepo.AssertCalled(t, "Update", originalQuestion)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Mock error dari repo
		questionRepo.On("Update", mock.AnythingOfType("*entities.Question")).
			Return(nil, assert.AnError).
			Once()

		result, err := usecase.UpdateQuestion(originalQuestion)

		assert.Error(t, err)
		assert.Nil(t, result)

		questionRepo.AssertCalled(t, "Update", originalQuestion)
	})
}
