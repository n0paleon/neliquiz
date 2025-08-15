package usecase

import (
	categoryDomain "NeliQuiz/internal/features/category/domain"
	categoryService "NeliQuiz/internal/features/category/usecase/mocks"
	questionDomain "NeliQuiz/internal/features/question/domain"
	"NeliQuiz/internal/features/question/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUsecase_CreateQuestion(t *testing.T) {
	questionRepo := new(mocks.QuestionRepositoryMock)
	categoryUseCase := new(categoryService.CategoryUseCaseMock)

	usecase := NewQuestionUseCase(questionRepo, categoryUseCase)

	questionRepo.On("Create", mock.Anything).Return(&questionDomain.Question{}, nil)
	categoryUseCase.On("FindOrCreateBatch", mock.Anything).Return([]categoryDomain.Category{}, nil)

	questionPayload := questionDomain.Question{
		Content: "Siapa presiden pertama Indonesia?",
	}

	options := []questionDomain.Option{
		questionDomain.Option{
			Content:   "jokowi",
			IsCorrect: false,
		},
		questionDomain.Option{
			Content:   "prabowo",
			IsCorrect: false,
		},
		questionDomain.Option{
			Content:   "soekarno",
			IsCorrect: true,
		},
	}

	t.Run("should create question correctly", func(t *testing.T) {
		err := questionPayload.AddOption(options...)
		assert.NoError(t, err)

		_, err = usecase.CreateQuestion(&questionPayload)
		assert.NoError(t, err)
	})

	t.Run("should error when options is reach the maximum", func(t *testing.T) {
		additionalOptions := []questionDomain.Option{
			questionDomain.Option{
				Content:   "mohammad hatta",
				IsCorrect: false,
			},
			questionDomain.Option{
				Content:   "tan malaka",
				IsCorrect: false,
			},
			questionDomain.Option{
				Content:   "sby",
				IsCorrect: true,
			},
		}

		err := questionPayload.AddOption(additionalOptions...)
		assert.Error(t, err)
		assert.Len(t, options, 3)

		assert.True(t, len(questionPayload.Options) < 5) // 5 is options maximum number
	})

	t.Run("should return error if options exceed the maximum when calling CreateQuestion", func(t *testing.T) {
		// reset state
		questionPayload := questionDomain.Question{
			Content: "Siapakah penemu listrik?",
		}

		options := []questionDomain.Option{
			{Content: "Einstein", IsCorrect: false},
			{Content: "Tesla", IsCorrect: false},
			{Content: "Newton", IsCorrect: false},
			{Content: "Faraday", IsCorrect: false},
			{Content: "Galvani", IsCorrect: false},
			{Content: "Franklin", IsCorrect: true}, // kelebihan, harus error
		}

		err := questionPayload.AddOption(options...)
		assert.Error(t, err, "AddOption seharusnya error karena melebihi 5")

		// tetap paksa kirim ke usecase
		_, err = usecase.CreateQuestion(&questionPayload)
		assert.Error(t, err, "CreateQuestion seharusnya error karena jumlah opsi melebihi batas")
	})
}

func TestUseCase_ListQuestions(t *testing.T) {
	questionRepo := new(mocks.QuestionRepositoryMock)
	categoryUseCase := new(categoryService.CategoryUseCaseMock)
	usecase := NewQuestionUseCase(questionRepo, categoryUseCase)

	questionRepo.On("PaginateQuestions", mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return([]questionDomain.Question{}, int64(10), nil)

	_, total, err := usecase.GetListQuestions("", 1, 10, "", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(10), total)
}

func TestUseCase_ListQuestionsWithCategory(t *testing.T) {
	questionRepo := new(mocks.QuestionRepositoryMock)
	categoryUseCase := new(categoryService.CategoryUseCaseMock)
	usecase := NewQuestionUseCase(questionRepo, categoryUseCase)

	categoryUseCase.On("FindCategoryByName", mock.AnythingOfType("string")).Return(&categoryDomain.Category{}, nil)
	questionRepo.On("PaginateQuestionsByCategory", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return([]questionDomain.Question{}, int64(10), nil)
	questionRepo.On("PaginateQuestions", mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return([]questionDomain.Question{}, int64(10), nil)

	_, total, err := usecase.GetListQuestions("", 1, 10, "", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(10), total)
}

func TestUseCase_DeleteQuestion(t *testing.T) {
	questionRepo := new(mocks.QuestionRepositoryMock)
	categoryUseCase := new(categoryService.CategoryUseCaseMock)
	usecase := NewQuestionUseCase(questionRepo, categoryUseCase)

	questionRepo.On("DeleteById", mock.AnythingOfType("string")).Return(nil)

	err := usecase.DeleteQuestion("this-is-fake-question-id")
	assert.NoError(t, err)
}

func TestUseCase_GetQuestionDetail(t *testing.T) {
	questionRepo := new(mocks.QuestionRepositoryMock)
	categoryUseCase := new(categoryService.CategoryUseCaseMock)
	usecase := NewQuestionUseCase(questionRepo, categoryUseCase)

	questionRepo.On("FindById", mock.AnythingOfType("string")).Return(&questionDomain.Question{}, nil)

	question, err := usecase.GetQuestionDetail("this-is-fake-question-id")
	assert.NoError(t, err)
	assert.NotNil(t, question)
}

func TestUseCase_UpdateQuestion(t *testing.T) {
	questionRepo := new(mocks.QuestionRepositoryMock)
	categoryUseCase := new(categoryService.CategoryUseCaseMock)
	usecase := NewQuestionUseCase(questionRepo, categoryUseCase)

	categoryUseCase.On("FindOrCreateBatch", mock.Anything).Return([]categoryDomain.Category{}, nil)

	// Data awal yang mau diupdate
	originalQuestion := questionDomain.Question{
		Content: "Siapakah presiden pertama Indonesia?",
	}
	originalQuestion.ID = "question-123"
	_ = originalQuestion.AddOption(
		questionDomain.Option{Content: "Jokowi", IsCorrect: false},
		questionDomain.Option{Content: "Soekarno", IsCorrect: true},
	)

	// Data hasil update
	updatedQuestion := originalQuestion
	updatedQuestion.Content = "Siapa presiden pertama RI?"
	updatedQuestion.Options[0].Content = "Prabowo"

	t.Run("should update question successfully", func(t *testing.T) {
		// Mock pemanggilan repo
		questionRepo.On("Update", mock.Anything).
			Return(&updatedQuestion, nil).
			Once()

		result, err := usecase.UpdateQuestion(&originalQuestion)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Siapa presiden pertama RI?", result.Content)
		assert.Equal(t, "Prabowo", result.Options[0].Content)

		questionRepo.AssertCalled(t, "Update", &originalQuestion)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Mock error dari repo
		questionRepo.On("Update", mock.Anything).
			Return(nil, assert.AnError).
			Once()

		result, err := usecase.UpdateQuestion(&originalQuestion)
		assert.Error(t, err)
		assert.Nil(t, result)

		questionRepo.AssertCalled(t, "Update", &originalQuestion)
	})
}

func TestUseCase_GetRandomQuestion(t *testing.T) {
	questionRepo := new(mocks.QuestionRepositoryMock)
	categoryUseCase := new(categoryService.CategoryUseCaseMock)
	usecase := NewQuestionUseCase(questionRepo, categoryUseCase)

	expectedQuestion := &questionDomain.Question{
		ID:      "123",
		Content: "Apa ibu kota Indonesia?",
		Options: []questionDomain.Option{
			{
				ID:        "1",
				Content:   "Jakarta",
				IsCorrect: true,
			},
			{
				ID:        "2",
				Content:   "Bandung",
				IsCorrect: false,
			},
			{
				ID:        "3",
				Content:   "Surabaya",
				IsCorrect: false,
			},
			{
				ID:        "4",
				Content:   "Yogyakarta",
				IsCorrect: false,
			},
		},
	}

	questionRepo.On("GetRandom").Return(expectedQuestion, nil)

	result, err := usecase.GetRandomQuestion()

	assert.NoError(t, err)
	assert.Equal(t, expectedQuestion, result)
	questionRepo.AssertExpectations(t)
}

func TestGetRandomQuestion_Error(t *testing.T) {
	questionRepo := new(mocks.QuestionRepositoryMock)
	categoryUseCase := new(categoryService.CategoryUseCaseMock)
	usecase := NewQuestionUseCase(questionRepo, categoryUseCase)

	questionRepo.On("GetRandom").Return(nil, assert.AnError)

	result, err := usecase.GetRandomQuestion()

	assert.Error(t, err)
	assert.Nil(t, result)
	questionRepo.AssertExpectations(t)
}

func TestCheckAnswer_Correct(t *testing.T) {
	questionRepo := new(mocks.QuestionRepositoryMock)
	categoryUseCase := new(categoryService.CategoryUseCaseMock)
	usecase := NewQuestionUseCase(questionRepo, categoryUseCase)

	option := questionDomain.Option{
		ID:        "option1",
		Content:   "Jakarta",
		IsCorrect: true,
	}

	question := &questionDomain.Question{
		ID:      "question1",
		Content: "Ibu kota Indonesia?",
		Options: []questionDomain.Option{option},
	}

	questionRepo.On("FindById", "question1").Return(question, nil)

	// Call usecase
	isCorrect, selected, _, err := usecase.CheckAnswer("question1", "option1")

	assert.NoError(t, err)
	assert.True(t, isCorrect)
	assert.Equal(t, &option, selected)
	questionRepo.AssertExpectations(t)
}

func TestCheckAnswer_InvalidQuestion(t *testing.T) {
	questionRepo := new(mocks.QuestionRepositoryMock)
	categoryUseCase := new(categoryService.CategoryUseCaseMock)
	usecase := NewQuestionUseCase(questionRepo, categoryUseCase)

	questionRepo.On("FindById", "invalid-id").Return(nil, assert.AnError)

	isCorrect, selected, _, err := usecase.CheckAnswer("invalid-id", "option1")

	assert.Error(t, err)
	assert.False(t, isCorrect)
	assert.Nil(t, selected)
	questionRepo.AssertExpectations(t)
}

func TestCheckAnswer_WrongOption(t *testing.T) {
	questionRepo := new(mocks.QuestionRepositoryMock)
	categoryUseCase := new(categoryService.CategoryUseCaseMock)
	usecase := NewQuestionUseCase(questionRepo, categoryUseCase)

	option := questionDomain.Option{
		ID:        "option1",
		Content:   "Surabaya",
		IsCorrect: false,
	}

	question := &questionDomain.Question{
		ID:      "question1",
		Content: "Ibu kota Indonesia?",
		Options: []questionDomain.Option{option},
	}

	questionRepo.On("FindById", "question1").Return(question, nil)

	isCorrect, selected, _, err := usecase.CheckAnswer("question1", "option1")

	assert.NoError(t, err)
	assert.False(t, isCorrect)
	assert.Equal(t, &option, selected)
	questionRepo.AssertExpectations(t)
}
