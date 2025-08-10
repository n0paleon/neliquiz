package domain

import "NeliQuiz/internal/domain/entities"

type AdminQuestionUseCase interface {
	CreateQuestion(q *entities.Question) (*entities.Question, error)
	GetListQuestions(category string, page, limit int, sortBy, order string) ([]entities.Question, int64, error)
	DeleteQuestion(id string) error
	UpdateQuestion(q *entities.Question) (*entities.Question, error)
	GetQuestionDetail(id string) (*entities.Question, error)
}

type UserQuestionUseCase interface {
	GetRandomQuestion() (*entities.Question, error)
	CheckAnswer(questionID, selectedQuestionID string) (isCorrect bool, option *entities.Option, explanationURL string, err error)
}

type QuestionRepository interface {
	Create(q *entities.Question) (*entities.Question, error)
	FindById(id string) (*entities.Question, error)
	DeleteById(id string) error
	GetRandom() (*entities.Question, error)
	PaginateQuestions(page, limit int, sortBy, order string) ([]entities.Question, int64, error)
	PaginateQuestionsByCategory(categoryID string, page, limit int, sortBy, order string) ([]entities.Question, int64, error)
	Update(q *entities.Question) (*entities.Question, error)
}
