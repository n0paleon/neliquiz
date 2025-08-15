package domain

import (
	"NeliQuiz/internal/features/category/domain"
	"errors"
	"strings"
	"time"
)

type Question struct {
	ID             string            `json:"id"`
	Content        string            `json:"content"`
	Hit            int               `json:"hit"`
	Options        []Option          `json:"options"`
	Categories     []domain.Category `json:"categories"`
	ExplanationURL string            `json:"explanation_url"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

type Option struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	IsCorrect bool   `json:"is_correct"`
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

// Validate checks if the question satisfies all rules:
// 1. Max 5 options allowed
// 2. Option contents must be unique (case-insensitive)
// 3. Exactly one option must be marked as correct
func (q *Question) Validate() error {
	if len(q.Options) == 0 {
		return errors.New("no options provided")
	}

	if len(q.Options) > 5 {
		return errors.New("the maximum number of answer choices is 5")
	}

	// Check for duplicate option content (case-insensitive)
	seen := make(map[string]bool)
	for _, o := range q.Options {
		content := normalizeString(o.Content)
		if seen[content] {
			return errors.New("duplicate option content found: " + o.Content)
		}
		seen[content] = true
	}

	// Check correct answer count
	correctCount := 0
	for _, o := range q.Options {
		if o.IsCorrect {
			correctCount++
		}
	}

	if correctCount == 0 {
		return errors.New("must have one correct answer")
	}
	if correctCount > 1 {
		return errors.New("only one option can be marked as correct")
	}

	return nil
}

// normalizeString converts a string to lowercase and trims spaces.
// You can extend this later with unicode normalization if needed.
func normalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
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
