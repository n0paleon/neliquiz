package dto

type GetRandomQuestionRequest struct {
	Categories []string `json:"categories"`
}

type VerifyAnswerRequest struct {
	QuestionID       string `json:"question_id"`
	SelectedOptionID string `json:"selected_option_id"`
}
