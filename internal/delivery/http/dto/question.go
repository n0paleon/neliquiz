package dto

import (
	"NeliQuiz/internal/domain/entities"
	"time"
)

type CreateQuestionRequest struct {
	Content        string   `json:"content" validate:"required,min=1,max=1000"`
	Options        []Option `json:"options" validate:"required,dive"`
	Categories     []string `json:"categories"`
	ExplanationURL string   `json:"explanation_url" validate:"url"`
}

type Option struct {
	Content   string `json:"content" validate:"required,min=1,max=500"`
	IsCorrect bool   `json:"is_correct"`
}

func (r *CreateQuestionRequest) ToEntity() (*entities.Question, error) {
	question, err := entities.NewQuestion(r.Content, r.ExplanationURL)
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

	categories := make([]entities.Category, len(r.Categories))
	for i, c := range r.Categories {
		cat, err := entities.NewCategory(c)
		if err != nil {
			return nil, err
		}
		categories[i] = *cat
	}

	question.Options = options
	question.Categories = categories
	return question, nil
}

type GetListQuestionsResponse struct {
	ID         string     `json:"id"`
	Content    string     `json:"content"`
	Hit        int        `json:"hit"`
	Categories []Category `json:"categories"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type GetQuestionDetailResponse struct {
	ID             string                            `json:"id"`
	Content        string                            `json:"content"`
	Hit            int                               `json:"hit"`
	Options        []GetQuestionDetailOptionResponse `json:"options"`
	Categories     []Category                        `json:"categories"`
	ExplanationURL string                            `json:"explanation_url"`
	CreatedAt      time.Time                         `json:"created_at"`
	UpdatedAt      time.Time                         `json:"updated_at"`
}

type GetQuestionDetailOptionResponse struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	IsCorrect bool   `json:"is_correct"`
}

func EntityToGetQuestionDetailResponse(e *entities.Question) *GetQuestionDetailResponse {
	results := make([]GetQuestionDetailOptionResponse, len(e.Options))
	for i, o := range e.Options {
		results[i] = GetQuestionDetailOptionResponse{
			ID:        o.ID,
			Content:   o.Content,
			IsCorrect: o.IsCorrect,
		}
	}

	categories := make([]Category, len(e.Categories))
	for i, c := range e.Categories {
		categories[i] = Category{
			ID:   c.ID,
			Name: c.Name,
		}
	}

	return &GetQuestionDetailResponse{
		ID:             e.ID,
		Content:        e.Content,
		Hit:            e.Hit,
		Options:        results,
		Categories:     categories,
		ExplanationURL: e.ExplanationURL,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}
