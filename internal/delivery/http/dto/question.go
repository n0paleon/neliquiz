package dto

import (
	"NeliQuiz/internal/domain/entities"
	"time"
)

type CreateQuestionRequest struct {
	Content string   `json:"content" validate:"required,min=1,max=1000"`
	Options []Option `json:"options" validate:"required,dive"`
}

type Option struct {
	Content   string `json:"content" validate:"required,min=1,max=500"`
	IsCorrect bool   `json:"is_correct"`
}

func (r *CreateQuestionRequest) ToEntity() (*entities.Question, error) {
	question, err := entities.NewQuestion(r.Content)
	if err != nil {
		return nil, err
	}

	var options []entities.Option
	for _, o := range r.Options {
		options = append(options, entities.Option{
			Content:   o.Content,
			IsCorrect: o.IsCorrect,
		})
	}

	question.Options = options
	return question, nil
}

type GetListQuestionsResponse struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostDeleteQuestionRequest struct {
	QuestionID string `json:"question_id" validate:"required,uuid"`
}
