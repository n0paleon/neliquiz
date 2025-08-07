package entities

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNewCategory_ValidName(t *testing.T) {
	name := "Science 101"
	cat, err := NewCategory(name)

	assert.NoError(t, err)
	assert.NotNil(t, cat)
	assert.Equal(t, name, cat.Name)
}

func TestNewCategory_EmptyName(t *testing.T) {
	name := ""
	cat, err := NewCategory(name)

	assert.Error(t, err)
	assert.Nil(t, cat)
	assert.EqualError(t, err, "category name cannot be empty")
}

func TestNewCategory_InvalidCharacters(t *testing.T) {
	names := []string{
		"Science!",       // symbol
		"Math@2025",      // symbol
		"Hello-World",    // symbol
		"Fiction/Nonfic", // symbol
	}

	for _, name := range names {
		cat, err := NewCategory(name)
		assert.Error(t, err, "expected error for name: %s", name)
		assert.Nil(t, cat)
		assert.EqualError(t, err, "category name can only contain letters, numbers, and spaces")
	}
}

func generateString(length int) string {
	return strings.Repeat("a", length)
}

func TestNewCategory_LongCharacters(t *testing.T) {
	name := generateString(51)
	cat, err := NewCategory(name)
	assert.Error(t, err)
	assert.Nil(t, cat)

	name = generateString(50)
	cat, err = NewCategory(name)
	assert.NoError(t, err)
	assert.NotNil(t, cat)
}

func TestNewCategory_ValidCharacters(t *testing.T) {
	names := []string{
		"Science",
		"Math 2025",
		"Kategori 1",
		"Belajar Bareng",
	}

	for _, name := range names {
		cat, err := NewCategory(name)
		assert.NoError(t, err, "should be valid for name: %s", name)
		assert.NotNil(t, cat)
	}
}
