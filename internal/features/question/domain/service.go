package domain

type QuestionUseCase interface {
	CreateQuestion(q *Question) (*Question, error)
	GetListQuestions(category string, page, limit int, sortBy, order string) ([]Question, int64, error)
	DeleteQuestion(id string) error
	UpdateQuestion(q *Question) (*Question, error)
	GetQuestionDetail(id string) (*Question, error)
	GetRandomQuestion(categories ...string) (*Question, error)
	CheckAnswer(questionID, selectedQuestionID string) (isCorrect bool, option *Option, explanationURL string, err error)
}
