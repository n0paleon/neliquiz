package domain

type QuestionRepository interface {
	Create(q *Question) (*Question, error)
	FindById(id string) (*Question, error)
	DeleteById(id string) error
	GetRandom() (*Question, error)
	PaginateQuestions(page, limit int, sortBy, order string) ([]Question, int64, error)
	PaginateQuestionsByCategory(categoryID string, page, limit int, sortBy, order string) ([]Question, int64, error)
	Update(q *Question) (*Question, error)
	GetRandomByCategoryNames(names []string) (*Question, error)
}
