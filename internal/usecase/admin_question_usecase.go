package usecase

import (
	"NeliQuiz/internal/domain"
	"NeliQuiz/internal/domain/entities"
	"NeliQuiz/internal/errorx"
	"errors"
)

type AdminQuestionUseCase struct {
	questionRepo domain.QuestionRepository
	categoryRepo domain.CategoryRepository
}

func (u *AdminQuestionUseCase) CreateQuestion(q *entities.Question) error {
	for i, cat := range q.Categories {
		if err := cat.Validate(); err != nil {
			return err
		}
		newCat, err := u.categoryRepo.FindOrCreateCategoryByName(cat.Name)
		if err != nil {
			return err
		}
		q.Categories[i] = *newCat
	}

	if err := q.Validate(); err != nil {
		return err
	}
	_, err := u.questionRepo.Create(q)
	if err != nil {
		return errorx.InternalError(err)
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

func (u *AdminQuestionUseCase) GetQuestionDetail(id string) (*entities.Question, error) {
	if id == "" {
		return nil, errors.New("invalid question id")
	}

	return u.questionRepo.FindById(id)
}

func NewAdminQuestionUseCase(questionRepository domain.QuestionRepository, categoryRepository domain.CategoryRepository) *AdminQuestionUseCase {
	return &AdminQuestionUseCase{
		questionRepo: questionRepository,
		categoryRepo: categoryRepository,
	}
}
