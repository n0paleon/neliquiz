package usecase

import (
	"NeliQuiz/internal/domain"
	"NeliQuiz/internal/domain/entities"
	"NeliQuiz/internal/errorx"
)

type UserQuestionUseCase struct {
	questionRepo domain.QuestionRepository
}

func (u *UserQuestionUseCase) GetRandomQuestion(categories ...string) (*entities.Question, error) {
	if len(categories) > 0 {
		result, err := u.questionRepo.GetRandomByCategoryNames(categories)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	result, err := u.questionRepo.GetRandom()
	if err != nil {
		return nil, errorx.InternalError(err)
	}

	return result, nil
}

func (u *UserQuestionUseCase) CheckAnswer(questionID, selectedOptionID string) (isCorrect bool, option *entities.Option, explanationURL string, err error) {
	question, err := u.questionRepo.FindById(questionID)
	if err != nil {
		isCorrect = false
		explanationURL = ""
		return
	}

	isCorrect, option, err = question.CheckAnswerWithOption(selectedOptionID)
	if err != nil {
		err = errorx.NotFound(err.Error())
	}
	explanationURL = question.ExplanationURL
	return
}

func NewUserQuestionUseCase(questionRepo domain.QuestionRepository) *UserQuestionUseCase {
	return &UserQuestionUseCase{
		questionRepo: questionRepo,
	}
}
