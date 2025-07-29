package usecase

import (
	"NeliQuiz/internal/domain"
	"NeliQuiz/internal/domain/entities"
	"errors"
)

type AdminQuestionUseCase struct {
	questionRepo domain.QuestionRepository
}

func (u *AdminQuestionUseCase) CreateQuestion(q *entities.Question) error {
	if err := q.ValidateAnswerKey(); err != nil {
		return err
	}
	_, err := u.questionRepo.Create(q)
	if err != nil {
		return errors.New("failed to create question")
	}
	return nil
}

func (u *AdminQuestionUseCase) GetListQuestions(page, limit int) ([]entities.Question, int64, error) {
	var results []entities.Question
	var total int64

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	results, total, err := u.questionRepo.PaginateQuestions(page, limit)
	if err != nil {
		return results, total, err
	}

	return results, total, nil
}

func (u *AdminQuestionUseCase) DeleteQuestion(id string) error {
	if id == "" {
		return errors.New("invalid question id")
	}

	return u.questionRepo.DeleteById(id)
}

func NewAdminQuestionUseCase(questionRepository domain.QuestionRepository) *AdminQuestionUseCase {
	return &AdminQuestionUseCase{
		questionRepo: questionRepository,
	}
}
