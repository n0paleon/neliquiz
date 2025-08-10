//go:build integration
// +build integration

package repository

import (
	"NeliQuiz/internal/config"
	"NeliQuiz/internal/infrastructures/database/postgres/testhelper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPGCategoryRepository(t *testing.T) {
	cfg := config.New()
	db := testhelper.NewTestDBConnection(t, cfg)
	categoryRepo := NewPGCategoryRepository(db)

	t.Cleanup(func() {
		testhelper.CleanupTable(t, db, "categories")
	})

	categoryName := "Bahasa Indonesia"

	t.Run("should create category correctly", func(t *testing.T) {
		category, err := categoryRepo.FindOrCreateCategoryByName(categoryName)
		assert.NoError(t, err)
		assert.NotNil(t, category)
		assert.Equal(t, categoryName, category.Name)
	})

	t.Run("should find category by name correctly", func(t *testing.T) {
		category, err := categoryRepo.FindCategoryByName(categoryName)
		assert.NoError(t, err)
		assert.NotNil(t, category)
		assert.Equal(t, categoryName, category.Name)
	})

	t.Run("should FindAll categories correctly", func(t *testing.T) {
		categories, err := categoryRepo.FindAll()
		assert.NoError(t, err)
		assert.Greater(t, len(categories), 0)
	})
}
