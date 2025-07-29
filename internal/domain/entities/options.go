package entities

import (
	"errors"
)

type Option struct {
	ID         string
	QuestionID string
	Content    string
	IsCorrect  bool
}

func NewOption(content string, isCorrect bool) (Option, error) {
	if content == "" {
		return Option{}, errors.New("content cannot be empty")
	}
	return Option{
		Content:   content,
		IsCorrect: isCorrect,
	}, nil
}
