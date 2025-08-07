package dto

import "NeliQuiz/internal/domain/entities"

type GetRandomQuestionResponse struct {
	QuestionID string       `json:"question_id"`
	Content    string       `json:"content"`
	Options    []QuizOption `json:"options"`
	Categories []Category   `json:"categories"`
}

type QuizOption struct {
	OptionID string `json:"option_id"`
	Content  string `json:"content"`
}

func EntityToGetRandomQuestionResponse(q *entities.Question) *GetRandomQuestionResponse {
	options := make([]QuizOption, len(q.Options))
	for i, option := range q.Options {
		options[i] = QuizOption{
			OptionID: option.ID,
			Content:  option.Content,
		}
	}

	categories := make([]Category, len(q.Categories))
	for i, category := range q.Categories {
		categories[i] = Category{
			ID:   category.ID,
			Name: category.Name,
		}
	}

	return &GetRandomQuestionResponse{
		QuestionID: q.ID,
		Content:    q.Content,
		Options:    options,
		Categories: categories,
	}
}

type PostVerifyAnswerRequest struct {
	SelectedOptionID string `json:"selected_option_id" validate:"required,min=1,max=1000,uuid"`
}

type PostVerifyAnswerResponse struct {
	Correct        bool       `json:"correct"`
	CorrectOption  QuizOption `json:"correct_option,omitempty"`
	ExplanationURL string     `json:"explanation_url"`
}
