package usecase

import (
	"NeliQuiz/internal/domain"
	"NeliQuiz/internal/domain/entities"
	"NeliQuiz/internal/errorx"
	"errors"
	"github.com/sirupsen/logrus"
)

type AdminQuestionUseCase struct {
	questionRepo domain.QuestionRepository
	categoryRepo domain.CategoryRepository
}

func (u *AdminQuestionUseCase) CreateQuestion(q *entities.Question) (*entities.Question, error) {
	categories, err := u.categoryRepo.FindOrCreateBatch(q.Categories)
	if err != nil {
		return nil, err
	}
	q.Categories = categories

	if err = q.Validate(); err != nil {
		return nil, err
	}

	result, err := u.questionRepo.Create(q)
	if err != nil {
		return nil, errorx.InternalError(err)
	}

	return result, nil
}

func (u *AdminQuestionUseCase) GetListQuestions(categoryName string, page, limit int, sortBy, order string) (results []entities.Question, total int64, err error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	if categoryName != "" {
		var category *entities.Category
		category, err = u.categoryRepo.FindCategoryByName(categoryName)
		if err != nil {
			err = errorx.NotFound("category not found")
			return
		}
		results, total, err = u.questionRepo.PaginateQuestionsByCategory(category.ID, page, limit, sortBy, order)
		return
	}

	results, total, err = u.questionRepo.PaginateQuestions(page, limit, sortBy, order)
	return
}

func (u *AdminQuestionUseCase) DeleteQuestion(id string) error {
	if id == "" {
		return errors.New("invalid question id")
	}

	return u.questionRepo.DeleteById(id)
}

func (u *AdminQuestionUseCase) GetQuestionDetail(id string) (*entities.Question, error) {
	if id == "" {
		return nil, errorx.NotFound("question not found")
	}

	return u.questionRepo.FindById(id)
}

func (u *AdminQuestionUseCase) UpdateQuestion(q *entities.Question) (*entities.Question, error) {
	if q.ID == "" {
		return nil, errorx.NotFound("question not found")
	}

	categories, err := u.categoryRepo.FindOrCreateBatch(q.Categories)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	q.Categories = categories

	if err = q.Validate(); err != nil {
		return nil, errorx.BadRequest(err.Error())
	}

	newQuestionData, err := u.questionRepo.Update(q)
	if err != nil {
		return nil, err
	}

	return newQuestionData, nil
}

func NewAdminQuestionUseCase(questionRepository domain.QuestionRepository, categoryRepository domain.CategoryRepository) *AdminQuestionUseCase {
	return &AdminQuestionUseCase{
		questionRepo: questionRepository,
		categoryRepo: categoryRepository,
	}
}
