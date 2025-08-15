//go:build integration
// +build integration

package repository

import (
	"NeliQuiz/internal/infrastructures/database/postgres/testhelper"
	"NeliQuiz/internal/shared/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPGCategoryRepository(t *testing.T) {
	cfg := config.New(".env.test")
	db := testhelper.NewTestDBConnection(t, cfg)
	categoryRepo := NewCategoryRepository(db)

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

func TestPGCategoryRepository_SearchCategoriesByName(t *testing.T) {
	cfg := config.New(".env.test")
	db := testhelper.NewTestDBConnection(t, cfg)
	categoryRepo := NewCategoryRepository(db)

	t.Cleanup(func() {
		testhelper.CleanupTable(t, db, "categories")
	})

	// Seed beberapa kategori
	names := []string{"Matematika", "Bahasa Indonesia", "Bahasa Inggris", "Fisika", "Kimia"}
	for _, n := range names {
		_, err := categoryRepo.FindOrCreateCategoryByName(n)
		assert.NoError(t, err)
	}

	t.Run("should return categories matching partial name", func(t *testing.T) {
		query := "bahasa"
		results, err := categoryRepo.SearchCategoriesByName(query, 10)
		assert.NoError(t, err)
		assert.Len(t, results, 2)
		expectedNames := []string{"Bahasa Indonesia", "Bahasa Inggris"}
		for _, res := range results {
			assert.Contains(t, expectedNames, res.Name)
		}
	})

	t.Run("should be case-insensitive", func(t *testing.T) {
		query := "MATEMATIKA"
		results, err := categoryRepo.SearchCategoriesByName(query, 10)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Matematika", results[0].Name)
	})

	t.Run("should respect limit", func(t *testing.T) {
		query := "a" // banyak kategori mengandung "a"
		results, err := categoryRepo.SearchCategoriesByName(query, 3)
		assert.NoError(t, err)
		assert.LessOrEqual(t, len(results), 3)
	})

	t.Run("should return nil for empty query", func(t *testing.T) {
		results, err := categoryRepo.SearchCategoriesByName("", 10)
		assert.NoError(t, err)
		assert.Nil(t, results)
	})
}
