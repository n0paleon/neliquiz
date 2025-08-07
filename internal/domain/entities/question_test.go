package entities

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestQuestion(t *testing.T) {
	var question *Question
	t.Run("should return error when content is empty", func(t *testing.T) {
		question, err := NewQuestion("", "")

		assert.Nil(t, question)
		assert.Error(t, err)
	})

	t.Run("should return question correctly", func(t *testing.T) {
		content := "kapan indonesia merdeka?"
		question, err := NewQuestion(content, "")

		assert.NoError(t, err)
		assert.NotNil(t, question)
		assert.Equal(t, content, question.Content)
	})

	t.Run("should return error when given option is empty", func(t *testing.T) {
		option := Option{}
		err := question.AddOption(option)
		assert.Error(t, err)
	})

	t.Run("should return error when question is having options greater than 5", func(t *testing.T) {
		var options []Option
		question, err := NewQuestion("ini adalah pertanyaan?", "")
		assert.NoError(t, err)

		for i := 0; i < 6; i++ {
			options = append(options, Option{Content: "ini adalah jawabannya"})
		}

		err = question.AddOption(options...)
		assert.Error(t, err)
	})

	t.Run("should return question correctly", func(t *testing.T) {
		content := "kapan indonesia merdeka?"
		question, err := NewQuestion(content, "")
		assert.NoError(t, err)
		assert.NotNil(t, question)

		contentOpt := "1945"
		option, err := NewOption(contentOpt, true)
		assert.NoError(t, err)
		assert.NotNil(t, option)

		err = question.AddOption(option)
		assert.NoError(t, err)

		assert.Equal(t, len(question.Options), 1)
		assert.Equal(t, question.Options[0].Content, contentOpt)
	})
}

func TestQuestionValidateRejectsDuplicateOptions(t *testing.T) {
	q, err := NewQuestion("Apa warna bendera Indonesia?", "")
	if err != nil {
		t.Fatalf("failed to create question: %v", err)
	}

	options := []Option{
		{ID: "1", Content: "Merah Putih", IsCorrect: true},
		{ID: "2", Content: "   merah putih", IsCorrect: false}, // duplikat (trimmed & lowercase)
		{ID: "3", Content: "Biru", IsCorrect: false},
		{ID: "4", Content: "MERAH PUTIH", IsCorrect: false}, // duplikat
	}

	err = q.AddOption(options...)
	if err != nil {
		t.Fatalf("failed to add options: %v", err)
	}

	validateErr := q.Validate()
	if validateErr == nil {
		t.Fatalf("expected validation error due to duplicate option content, got nil")
	}

	expected := "duplicate option content found"
	if !strings.Contains(validateErr.Error(), expected) {
		t.Fatalf("expected error containing '%s', got '%v'", expected, validateErr)
	}
}
