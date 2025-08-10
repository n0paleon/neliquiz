//go:build integration
// +build integration

package repository

import (
	"NeliQuiz/internal/config"
	"NeliQuiz/internal/domain/entities"
	"NeliQuiz/internal/infrastructures/database/postgres/testhelper"
	"NeliQuiz/internal/repository/schema"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPGQuestionRepository(t *testing.T) {
	cfg := config.New()
	db := testhelper.NewTestDBConnection(t, cfg)
	questionRepo := NewPGQuestionRepository(db)

	t.Cleanup(func() {
		testhelper.CleanupTable(t, db, "questions")
	})

	payload := entities.Question{
		Content: "hello world",
	}
	var questionID string

	t.Run("should create question correctly", func(t *testing.T) {
		result, err := questionRepo.Create(&payload)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, result.Hit, 0)
		assert.Equal(t, payload.Content, result.Content)

		questionID = result.ID
	})

	t.Run("should find question by id correctly", func(t *testing.T) {
		result, err := questionRepo.FindById(questionID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, questionID, result.ID)
	})

	t.Run("should return question records correctly", func(t *testing.T) {
		results, total, err := questionRepo.PaginateQuestions(1, 10, "", "")
		assert.NoError(t, err)
		assert.NotEqual(t, total, int64(0))
		assert.Equal(t, payload.Content, results[0].Content)
	})

	t.Run("should return error when request delete with empty id", func(t *testing.T) {
		err := questionRepo.DeleteById("")
		assert.Error(t, err)
	})

	t.Run("should return error when request delete with invalid id", func(t *testing.T) {
		err := questionRepo.DeleteById("abcd-invalid-id")
		assert.Error(t, err)
	})

	t.Run("should delete question correctly", func(t *testing.T) {
		err := questionRepo.DeleteById(questionID)
		assert.NoError(t, err)
	})

	t.Run("should return error when trying to get deleted question", func(t *testing.T) {
		result, err := questionRepo.FindById(questionID)
		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

func TestPGQuestionRepository_GetRandomAndValidate(t *testing.T) {
	cfg := config.New()
	db := testhelper.NewTestDBConnection(t, cfg)
	questionRepo := NewPGQuestionRepository(db)

	t.Cleanup(func() {
		testhelper.CleanupTable(t, db, "questions")
	})

	// Buat 4 pertanyaan
	for i := 1; i <= 4; i++ {
		options := []entities.Option{
			{Content: fmt.Sprintf("Pilihan A soal %d", i), IsCorrect: false},
			{Content: fmt.Sprintf("Pilihan B soal %d", i), IsCorrect: true},
			{Content: fmt.Sprintf("Pilihan C soal %d", i), IsCorrect: false},
			{Content: fmt.Sprintf("Pilihan D soal %d", i), IsCorrect: false},
		}

		question := &entities.Question{
			Content: fmt.Sprintf("Ini soal ke-%d", i),
			Options: options,
		}
		created, err := questionRepo.Create(question)
		assert.NoError(t, err)
		assert.NotNil(t, created)
		assert.Len(t, created.Options, 4)
	}

	// Ambil salah satu soal secara acak
	randomQuestion, err := questionRepo.GetRandom()
	assert.NoError(t, err)
	assert.NotNil(t, randomQuestion)
	assert.Len(t, randomQuestion.Options, 4)

	var correctOptionID string
	for _, opt := range randomQuestion.Options {
		if opt.IsCorrect {
			correctOptionID = opt.ID
			break
		}
	}
	assert.NotEmpty(t, correctOptionID)

	// Validasi jawaban
	result, selectedOption, err := randomQuestion.CheckAnswerWithOption(correctOptionID)
	assert.NoError(t, err)
	assert.True(t, result)
	assert.True(t, selectedOption.IsCorrect)
}

func TestPGQuestionRepository_PaginateQuestionsByCategory(t *testing.T) {
	cfg := config.New()
	db := testhelper.NewTestDBConnection(t, cfg)
	questionRepo := NewPGQuestionRepository(db)

	t.Cleanup(func() {
		testhelper.CleanupTable(t, db, "question_categories")
		testhelper.CleanupTable(t, db, "categories")
		testhelper.CleanupTable(t, db, "questions")
	})

	// 1. Buat kategori
	categorySchemaPayload := schema.Category{
		Name: "Matematika",
	}
	err := db.Create(&categorySchemaPayload).Error
	assert.NoError(t, err)
	assert.NotEmpty(t, categorySchemaPayload.ID)

	// 2. Buat pertanyaan dengan kategori tersebut
	question := entities.Question{
		Content:    "2 + 2 = ?",
		Categories: []entities.Category{{ID: categorySchemaPayload.ID, Name: categorySchemaPayload.Name}},
	}
	created, err := questionRepo.Create(&question)
	assert.NoError(t, err)
	assert.NotNil(t, created)

	// 3. Panggil PaginateQuestionsByCategory
	results, total, err := questionRepo.PaginateQuestionsByCategory(categorySchemaPayload.ID, 1, 10, "", "")
	assert.NoError(t, err)
	assert.Greater(t, total, int64(0))
	assert.Len(t, results, 1)
	assert.Equal(t, "2 + 2 = ?", results[0].Content)
	assert.NotEmpty(t, results[0].Categories)
	assert.Equal(t, "Matematika", results[0].Categories[0].Name)
}

func TestPGQuestionRepository_Update(t *testing.T) {
	cfg := config.New()
	db := testhelper.NewTestDBConnection(t, cfg)
	repo := NewPGQuestionRepository(db)

	t.Cleanup(func() {
		testhelper.CleanupTable(t, db, "question_categories")
		testhelper.CleanupTable(t, db, "categories")
		testhelper.CleanupTable(t, db, "questions")
	})

	// Seed kategori awal
	cat1 := schema.Category{Name: "Math"}
	cat2 := schema.Category{Name: "Science"}
	cat3 := schema.Category{Name: "Physics"}

	assert.NoError(t, db.Create(&cat1).Error)
	assert.NotEmpty(t, cat1.ID)
	assert.NoError(t, db.Create(&cat2).Error)
	assert.NotEmpty(t, cat2.ID)
	assert.NoError(t, db.Create(&cat3).Error)
	assert.NotEmpty(t, cat3.ID)

	t.Logf("Created categories - Math: %s, Science: %s, Physics: %s", cat1.ID, cat2.ID, cat3.ID)

	// Seed question awal dengan multiple categories
	initialQuestion := schema.Question{
		Content:        "Original Content",
		Hit:            0,
		Options:        []entities.Option{{Content: "A", IsCorrect: false}},
		ExplanationURL: "http://original.com",
		Categories:     []schema.Category{cat1, cat2}, // Start dengan 2 kategori
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	created, err := repo.Create(initialQuestion.ToEntity())
	assert.NoError(t, err)
	assert.NotNil(t, created)

	t.Logf("Created question with ID: %s", created.ID)

	// Verifikasi initial state
	var initialState schema.Question
	err = db.Preload("Categories").First(&initialState, "id = ?", created.ID).Error
	assert.NoError(t, err)
	assert.Len(t, initialState.Categories, 2, "Should start with 2 categories")

	categoryNames := make([]string, len(initialState.Categories))
	for i, cat := range initialState.Categories {
		categoryNames[i] = cat.Name
	}
	assert.Contains(t, categoryNames, "Math")
	assert.Contains(t, categoryNames, "Science")

	t.Log("✅ Initial state verified: question has Math and Science categories")

	// Simpan timestamp awal untuk perbandingan
	originalUpdatedAt := initialState.UpdatedAt

	// Tunggu sebentar agar UpdatedAt berbeda
	time.Sleep(10 * time.Millisecond)

	// Data baru untuk update - test semua 3 skenario sekaligus:
	// - Keep Science (skenario 2: sudah ada, tetap ada)
	// - Add Physics (skenario 1: baru ditambahkan)
	// - Remove Math (skenario 3: ada di lama, tidak ada di baru)
	updatedQuestion := entities.Question{
		ID:             created.ID,
		Content:        "Updated Content",
		Hit:            5,
		Options:        []entities.Option{{Content: "B", IsCorrect: true}},
		ExplanationURL: "http://updated.com",
		Categories: []entities.Category{
			*cat2.ToEntity(), // Keep Science (skenario 2)
			*cat3.ToEntity(), // Add Physics (skenario 1)
			// Math akan dihapus (skenario 3)
		},
	}

	t.Logf("Updating question with categories: Science (keep), Physics (add), Math (remove)")

	// Jalankan Update
	result, err := repo.Update(&updatedQuestion)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	t.Logf("Update completed, result ID: %s", result.ID)

	// Verifikasi hasil dari return value
	assert.Equal(t, created.ID, result.ID)
	assert.Equal(t, "Updated Content", result.Content)
	assert.Equal(t, 5, result.Hit)
	assert.Equal(t, "http://updated.com", result.ExplanationURL)
	assert.Len(t, result.Options, 1)
	assert.Equal(t, "B", result.Options[0].Content)
	assert.True(t, result.Options[0].IsCorrect)

	// Verifikasi categories dari return value
	assert.Len(t, result.Categories, 2, "Should have 2 categories after update")
	resultCategoryNames := make([]string, len(result.Categories))
	for i, cat := range result.Categories {
		resultCategoryNames[i] = cat.Name
	}
	assert.Contains(t, resultCategoryNames, "Science", "Should keep Science category")
	assert.Contains(t, resultCategoryNames, "Physics", "Should add Physics category")
	assert.NotContains(t, resultCategoryNames, "Math", "Should remove Math category")

	// Double check: Ambil dari DB lagi untuk memastikan persistence
	var dbQuestion schema.Question
	err = db.Preload("Categories").First(&dbQuestion, "id = ?", created.ID).Error
	assert.NoError(t, err)

	// Verifikasi data basic
	assert.Equal(t, "Updated Content", dbQuestion.Content)
	assert.Equal(t, 5, dbQuestion.Hit)
	assert.Equal(t, "http://updated.com", dbQuestion.ExplanationURL)
	assert.Len(t, dbQuestion.Options, 1)
	assert.Equal(t, "B", dbQuestion.Options[0].Content)
	assert.True(t, dbQuestion.Options[0].IsCorrect)
	assert.True(t, dbQuestion.UpdatedAt.After(originalUpdatedAt), "UpdatedAt should be newer")

	// Verifikasi categories di database
	assert.Len(t, dbQuestion.Categories, 2, "Database should have 2 categories")
	dbCategoryNames := make([]string, len(dbQuestion.Categories))
	dbCategoryIDs := make([]string, len(dbQuestion.Categories))
	for i, cat := range dbQuestion.Categories {
		dbCategoryNames[i] = cat.Name
		dbCategoryIDs[i] = cat.ID
	}

	// Test skenario category update:
	assert.Contains(t, dbCategoryNames, "Science", "✅ Skenario 2: Should keep existing Science category")
	assert.Contains(t, dbCategoryNames, "Physics", "✅ Skenario 1: Should add new Physics category")
	assert.NotContains(t, dbCategoryNames, "Math", "✅ Skenario 3: Should remove old Math category")

	assert.Contains(t, dbCategoryIDs, cat2.ID, "Science category ID should match")
	assert.Contains(t, dbCategoryIDs, cat3.ID, "Physics category ID should match")
	assert.NotContains(t, dbCategoryIDs, cat1.ID, "Math category ID should not be present")

	t.Log("✅ All update scenarios verified successfully:")
	t.Log("  - Skenario 1: Added Physics category")
	t.Log("  - Skenario 2: Kept Science category")
	t.Log("  - Skenario 3: Removed Math category")

	// Verifikasi di join table secara langsung (optional, untuk debugging)
	var joinCount int64
	db.Table("question_categories").Where("question_id = ?", created.ID).Count(&joinCount)
	assert.Equal(t, int64(2), joinCount, "Join table should have exactly 2 relationships")

	t.Logf("✅ Join table verification: %d relationships found", joinCount)
}
