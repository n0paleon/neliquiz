package dto

import "NeliQuiz/internal/domain/entities"

type GetRandomQuestionResponse struct {
	QuestionID string       `json:"question_id"`
	Content    string       `json:"content"`
	Options    []QuizOption `json:"options"`
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

	return &GetRandomQuestionResponse{
		QuestionID: q.ID,
		Content:    q.Content,
		Options:    options,
	}
}

type PostVerifyAnswerRequest struct {
	QuestionID string `json:"question_id" validate:"required,min=1,max=1000,uuid"`
	OptionID   string `json:"option_id" validate:"required,min=1,max=1000,uuid"`
}

type PostVerifyAnswerResponse struct {
	Correct       bool       `json:"correct"`
	CorrectOption QuizOption `json:"correct_option,omitempty"`
	Explanation   string     `json:"explanation"`
}
