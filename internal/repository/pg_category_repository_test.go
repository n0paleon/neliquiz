package repository

import (
	"NeliQuiz/internal/config"
	"NeliQuiz/internal/infrastructures/database/postgres/testhelper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPGCategoryRepository_FindOrCreateCategoryByName(t *testing.T) {
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
}
