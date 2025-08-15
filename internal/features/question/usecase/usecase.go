package usecase

import (
	categoryDomain "NeliQuiz/internal/features/category/domain"
	questionDomain "NeliQuiz/internal/features/question/domain"
	"NeliQuiz/internal/shared/errorx"
	"errors"
	"github.com/sirupsen/logrus"
)

type UseCase struct {
	questionRepo    questionDomain.QuestionRepository
	categoryUseCase categoryDomain.CategoryUseCase
}

func (u *UseCase) CreateQuestion(q *questionDomain.Question) (*questionDomain.Question, error) {
	categories, err := u.categoryUseCase.FindOrCreateBatch(q.Categories)
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

func (u *UseCase) GetListQuestions(categoryName string, page, limit int, sortBy, order string) (results []questionDomain.Question, total int64, err error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	if categoryName != "" {
		var category *categoryDomain.Category
		category, err = u.categoryUseCase.FindCategoryByName(categoryName)
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

func (u *UseCase) DeleteQuestion(id string) error {
	if id == "" {
		return errors.New("invalid question id")
	}

	return u.questionRepo.DeleteById(id)
}

func (u *UseCase) GetQuestionDetail(id string) (*questionDomain.Question, error) {
	if id == "" {
		return nil, errorx.NotFound("question not found")
	}

	return u.questionRepo.FindById(id)
}

func (u *UseCase) UpdateQuestion(q *questionDomain.Question) (*questionDomain.Question, error) {
	if q.ID == "" {
		return nil, errorx.NotFound("question not found")
	}

	categories, err := u.categoryUseCase.FindOrCreateBatch(q.Categories)
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

func (u *UseCase) CheckAnswer(questionID, selectedOptionID string) (isCorrect bool, option *questionDomain.Option, explanationURL string, err error) {
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

func (u *UseCase) GetRandomQuestion(categories ...string) (*questionDomain.Question, error) {
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

func NewQuestionUseCase(questionRepository questionDomain.QuestionRepository, categoryUseCase categoryDomain.CategoryUseCase) *UseCase {
	return &UseCase{
		questionRepo:    questionRepository,
		categoryUseCase: categoryUseCase,
	}
}
