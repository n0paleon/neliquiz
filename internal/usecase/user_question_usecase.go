package usecase

import (
	"NeliQuiz/internal/domain"
	"NeliQuiz/internal/domain/entities"
	"errors"
)

type UserQuestionUseCase struct {
	questionRepo domain.QuestionRepository
}

func (u *UserQuestionUseCase) GetRandomQuestion() (*entities.Question, error) {
	result, err := u.questionRepo.GetRandom()
	if err != nil {
		return nil, errors.New("failed to get random question")
	}

	return result, nil
}

func (u *UserQuestionUseCase) CheckAnswer(questionID, selectedOptionID string) (bool, *entities.Option, error) {
	question, err := u.questionRepo.FindById(questionID)
	if err != nil {
		return false, nil, errors.New("invalid question id")
	}

	return question.CheckAnswerWithOption(selectedOptionID)
}

func NewUserQuestionUseCase(questionRepo domain.QuestionRepository) *UserQuestionUseCase {
	return &UserQuestionUseCase{
		questionRepo: questionRepo,
	}
}
