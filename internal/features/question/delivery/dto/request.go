package dto

import (
	categoryDomain "NeliQuiz/internal/features/category/domain"
	questionDomain "NeliQuiz/internal/features/question/domain"
)

type CreateQuestionRequest struct {
	Content        string         `json:"content" validate:"required,max=1000"`
	Options        []CreateOption `json:"options" validate:"required,dive"`
	Categories     []string       `json:"categories"`
	ExplanationURL string         `json:"explanation_url" validate:"omitempty,url"`
}

type CreateOption struct {
	Content   string `json:"content" validate:"required,max=1000"`
	IsCorrect bool   `json:"is_correct"`
}

func (d *CreateQuestionRequest) ToDomain() *questionDomain.Question {
	options := make([]questionDomain.Option, len(d.Options))
	for i, option := range d.Options {
		options[i] = questionDomain.Option{
			Content:   option.Content,
			IsCorrect: option.IsCorrect,
		}
	}

	categories := make([]categoryDomain.Category, len(d.Categories))
	for i, category := range d.Categories {
		categories[i] = categoryDomain.Category{
			Name: category,
		}
	}

	question := questionDomain.Question{
		Content:        d.Content,
		Options:        options,
		Categories:     categories,
		ExplanationURL: d.ExplanationURL,
	}

	return &question
}

type UpdateQuestionDetailRequest struct {
	Content        string                               `json:"content" validate:"max=1000"`
	Hit            int                                  `json:"hit" validate:"omitempty,numeric"`
	Options        []UpdateQuestionDetailOptionsRequest `json:"options" validate:"required,dive"`
	Categories     []string                             `json:"categories"`
	ExplanationURL string                               `json:"explanation_url" validate:"omitempty,url"`
}

type UpdateQuestionDetailOptionsRequest struct {
	ID        string `json:"id" validate:"max=100"`
	Content   string `json:"content" validate:"max=1000"`
	IsCorrect bool   `json:"is_correct"`
}

func (d *UpdateQuestionDetailRequest) ToDomain() *questionDomain.Question {
	options := make([]questionDomain.Option, len(d.Options))
	for i, option := range d.Options {
		options[i] = questionDomain.Option{
			ID:        option.ID,
			Content:   option.Content,
			IsCorrect: option.IsCorrect,
		}
	}

	categories := make([]categoryDomain.Category, len(d.Categories))
	for i, catName := range d.Categories {
		categories[i] = categoryDomain.Category{
			Name: catName,
		}
	}

	question := questionDomain.Question{
		Content:        d.Content,
		Options:        options,
		Categories:     categories,
		ExplanationURL: d.ExplanationURL,
	}

	return &question
}

type PostVerifyAnswerRequest struct {
	SelectedOptionID string `json:"selected_option_id" validate:"required"`
}
