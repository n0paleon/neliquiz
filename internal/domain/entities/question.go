package entities

import (
	"errors"
	"time"
)

type Question struct {
	ID        string
	Content   string
	Options   []Option
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewQuestion(content string) (*Question, error) {
	if content == "" {
		return nil, errors.New("content cannot be empty")
	}

	return &Question{
		Content:   content,
		Options:   make([]Option, 0),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// CheckAnswerWithOption checks if the selected option ID matches the correct answer,
// and also returns the full Option object. It returns:
// - a boolean indicating whether the answer is correct,
// - the corresponding Option (if found),
// - and an error if the option ID does not exist.
func (q *Question) CheckAnswerWithOption(selectedOptionID string) (bool, *Option, error) {
	for _, o := range q.Options {
		if o.ID == selectedOptionID {
			return o.IsCorrect, &o, nil
		}
	}
	return false, nil, errors.New("selected option not found in this question")
}

// ValidateAnswerKey ensures that the Question has exactly one correct answer.
// It returns an error if there are no correct answers or if more than one option is marked as correct.
func (q *Question) ValidateAnswerKey() error {
	var correctCount int

	if len(q.Options) <= 1 {
		return errors.New("the number of options should be greater than 1")
	}

	for _, o := range q.Options {
		if o.IsCorrect {
			correctCount++
		}
	}

	if correctCount == 0 {
		return errors.New("must have at least one correct answer")
	}

	if correctCount > 1 {
		return errors.New("only one option can be marked as correct")
	}

	return nil
}

func (q *Question) AddOption(option ...Option) error {
	for _, o := range option {
		if o.Content == "" {
			return errors.New("option content cannot be empty")
		}
	}

	if (len(q.Options) + len(option)) > 5 {
		return errors.New("the maximum number of answer choices is 5")
	}

	for _, o := range option {
		q.Options = append(q.Options, o)
	}
	q.UpdatedAt = time.Now()

	return nil
}
