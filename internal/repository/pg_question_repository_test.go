//go:build integration
// +build integration

package repository

import (
	"NeliQuiz/internal/config"
	"NeliQuiz/internal/domain/entities"
	"NeliQuiz/internal/infrastructures/database/postgres/testhelper"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
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
	questionID := ""

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
		results, total, err := questionRepo.PaginateQuestions(1, 10)
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
	questions := []*entities.Question{}
	for i := 1; i <= 4; i++ {
		options := []entities.Option{
			{
				Content:   fmt.Sprintf("Pilihan A soal %d", i),
				IsCorrect: false,
			},
			{
				Content:   fmt.Sprintf("Pilihan B soal %d", i),
				IsCorrect: true, // hanya satu yang benar
			},
			{
				Content:   fmt.Sprintf("Pilihan C soal %d", i),
				IsCorrect: false,
			},
			{
				Content:   fmt.Sprintf("Pilihan D soal %d", i),
				IsCorrect: false,
			},
		}

		question := &entities.Question{
			Content: fmt.Sprintf("Ini soal ke-%d", i),
			Options: options,
		}
		created, err := questionRepo.Create(question)
		assert.NoError(t, err)
		assert.NotNil(t, created)
		assert.Len(t, created.Options, 4)

		questions = append(questions, created)
	}

	// Ambil salah satu soal secara acak
	randomQuestion, err := questionRepo.GetRandom()
	assert.NoError(t, err)
	assert.NotNil(t, randomQuestion)
	assert.Len(t, randomQuestion.Options, 4, "randomQuestion.Options harus terisi dengan 4 data")

	// Cari ID jawaban benar
	var correctOptionID string
	for _, opt := range randomQuestion.Options {
		if opt.IsCorrect {
			correctOptionID = opt.ID
			break
		}
	}
	assert.NotEmpty(t, correctOptionID, "ID dari jawaban benar tidak boleh kosong")

	// Validasi jawaban
	result, selectedOption, err := randomQuestion.CheckAnswerWithOption(correctOptionID)
	assert.NoError(t, err)
	assert.True(t, result, "Hasil validasi seharusnya true untuk jawaban yang benar")
	assert.True(t, selectedOption.IsCorrect, "Opsi yang dipilih seharusnya benar")
}
