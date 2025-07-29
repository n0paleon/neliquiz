package schema

import (
	"NeliQuiz/internal/domain/entities"
	"time"
)

type Question struct {
	ID        string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Content   string
	Options   []Option `gorm:"foreignKey:QuestionID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Question) TableName() string {
	return "questions"
}

func (s *Question) ToEntity() *entities.Question {
	options := make([]entities.Option, len(s.Options))
	for i, option := range s.Options {
		options[i] = entities.Option{
			ID:         option.ID,
			QuestionID: option.QuestionID,
			Content:    option.Content,
			IsCorrect:  option.IsCorrect,
		}
	}

	return &entities.Question{
		ID:        s.ID,
		Content:   s.Content,
		Options:   options,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func ToQuestionSchema(q *entities.Question) *Question {
	options := make([]Option, len(q.Options))
	for i, option := range q.Options {
		options[i] = Option{
			ID:         option.ID,
			QuestionID: option.QuestionID,
			Content:    option.Content,
			IsCorrect:  option.IsCorrect,
		}
	}

	return &Question{
		ID:        q.ID,
		Content:   q.Content,
		Options:   options,
		CreatedAt: q.CreatedAt,
		UpdatedAt: q.UpdatedAt,
	}
}
