package schema

import (
	"NeliQuiz/internal/domain/entities"
	"time"
)

type Question struct {
	ID             string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Content        string
	Hit            int
	Options        []Option   `gorm:"foreignKey:QuestionID"`
	Categories     []Category `gorm:"many2many:question_categories"`
	ExplanationURL string
	CreatedAt      time.Time
	UpdatedAt      time.Time
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

	categories := make([]entities.Category, len(s.Categories))
	for i, category := range s.Categories {
		categories[i] = *category.ToEntity()
	}

	return &entities.Question{
		ID:             s.ID,
		Content:        s.Content,
		Hit:            s.Hit,
		Options:        options,
		Categories:     categories,
		ExplanationURL: s.ExplanationURL,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
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

	categories := make([]Category, len(q.Categories))
	for i, category := range q.Categories {
		categories[i] = *ToCategorySchema(&category)
	}

	return &Question{
		ID:             q.ID,
		Content:        q.Content,
		Hit:            q.Hit,
		Options:        options,
		Categories:     categories,
		ExplanationURL: q.ExplanationURL,
		CreatedAt:      q.CreatedAt,
		UpdatedAt:      q.UpdatedAt,
	}
}
