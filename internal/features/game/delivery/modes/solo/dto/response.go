package dto

import (
	categoryDomain "NeliQuiz/internal/features/category/domain"
	questionDomain "NeliQuiz/internal/features/question/domain"
)

type VerifyAnswerResponse struct {
	IsCorrect      bool                `json:"is_correct"`
	CorrectOption  VerifyAnswer_Option `json:"correct_option"`
	ExplanationURL string              `json:"explanation_url"`
}

type VerifyAnswer_Option struct {
	OptionID string `json:"option_id"`
	Content  string `json:"content"`
}

type GetRandomQuestionResponse struct {
	QuestionID string                    `json:"question_id"`
	Content    string                    `json:"content"`
	Options    []GetRandomQuestionOption `json:"options"`
	Categories []categoryDomain.Category `json:"categories"`
}

type GetRandomQuestionOption struct {
	OptionID string `json:"option_id"`
	Content  string `json:"content"`
}

func ToGetRandomQuestionResponse(q *questionDomain.Question) *GetRandomQuestionResponse {
	options := make([]GetRandomQuestionOption, len(q.Options))
	for i, option := range q.Options {
		options[i] = GetRandomQuestionOption{
			OptionID: option.ID,
			Content:  option.Content,
		}
	}

	question := &GetRandomQuestionResponse{
		QuestionID: q.ID,
		Content:    q.Content,
		Options:    options,
		Categories: q.Categories,
	}

	return question
}
