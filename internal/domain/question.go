package domain

import "NeliQuiz/internal/domain/entities"

type AdminQuestionUseCase interface {
	CreateQuestion(q *entities.Question) error
	GetListQuestions(page, limit int) ([]entities.Question, int64, error)
	DeleteQuestion(id string) error
}

type UserQuestionUseCase interface {
	GetRandomQuestion() (*entities.Question, error)
	CheckAnswer(questionID, selectedQuestionID string) (bool, *entities.Option, error)
}

type QuestionRepository interface {
	Create(q *entities.Question) (*entities.Question, error)
	FindById(id string) (*entities.Question, error)
	DeleteById(id string) error
	GetRandom() (*entities.Question, error)
	PaginateQuestions(page, limit int) ([]entities.Question, int64, error)
}
