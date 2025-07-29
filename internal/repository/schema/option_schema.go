package schema

import (
	"NeliQuiz/internal/domain/entities"
)

type Option struct {
	ID         string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;"`
	QuestionID string
	Content    string
	IsCorrect  bool
}

func (s *Option) TableName() string {
	return "question_options"
}

func (s *Option) ToEntity() *entities.Option {
	return &entities.Option{
		ID:         s.ID,
		QuestionID: s.QuestionID,
		Content:    s.Content,
		IsCorrect:  s.IsCorrect,
	}
}

func ToOptionSchema(e *entities.Option) *Option {
	return &Option{
		ID:         e.ID,
		QuestionID: e.QuestionID,
		Content:    e.Content,
		IsCorrect:  e.IsCorrect,
	}
}
