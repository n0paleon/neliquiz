package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewOption(t *testing.T) {
	t.Run("should return error when content is empty", func(t *testing.T) {
		opt, err := NewOption("", false)
		assert.Error(t, err)
		assert.Empty(t, opt.Content)
	})

	t.Run("should create option correctly", func(t *testing.T) {
		content := "1945"

		opt, err := NewOption(content, true)
		assert.NoError(t, err)
		assert.NotNil(t, opt)

		assert.Equal(t, content, opt.Content)
		assert.Equal(t, true, opt.IsCorrect)
	})
}
