package repository

import (
	"NeliQuiz/internal/domain/entities"
	"NeliQuiz/internal/errorx"
	"NeliQuiz/internal/repository/schema"
	"crypto/rand"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math/big"
)

type PGQuestionRepository struct {
	db *gorm.DB
}

func (r *PGQuestionRepository) Create(q *entities.Question) (*entities.Question, error) {
	model := schema.ToQuestionSchema(q)
	if err := r.db.Create(model).Error; err != nil {
		logrus.Error(err.Error())
		return nil, TranslateGormError(err)
	}
	return model.ToEntity(), nil
}

func (r *PGQuestionRepository) FindById(id string) (*entities.Question, error) {
	var result schema.Question
	if err := r.db.
		Preload("Options").
		Preload("Categories").
		Where("id = ?", id).
		First(&result).Error; err != nil {
		return nil, TranslateGormError(err)
	}
	return result.ToEntity(), nil
}

func (r *PGQuestionRepository) DeleteById(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	err := r.db.Where("id = ?", id).Delete(&schema.Question{})
	if err.Error != nil {
		return TranslateGormError(err.Error)
	}
	if err.RowsAffected < 1 {
		return errorx.NotFound("question not found")
	}

	return nil
}

func (r *PGQuestionRepository) GetRandom() (*entities.Question, error) {
	var questions []schema.Question
	if err := r.db.
		Preload("Options").
		Preload("Categories").
		Order("updated_at ASC").
		Limit(20).
		Find(&questions).Error; err != nil {
		return nil, TranslateGormError(err)
	}

	n := big.NewInt(int64(len(questions)))
	i, _ := rand.Int(rand.Reader, n)
	selected := questions[i.Int64()]

	go func() {
		_ = r.updateHit(selected.ID)
	}()

	return selected.ToEntity(), nil
}

func (r *PGQuestionRepository) updateHit(id string) error {
	err := r.db.Model(schema.Question{}).
		Where("id = ?", id).
		UpdateColumn("hit", gorm.Expr("hit + 1")).
		Error

	return TranslateGormError(err)
}

func (r *PGQuestionRepository) PaginateQuestions(page, limit int) ([]entities.Question, int64, error) {
	var questions []schema.Question
	var total int64

	r.db.Model(&schema.Question{}).Count(&total)

	offset := (page - 1) * limit
	if err := r.db.
		Preload("Categories").
		Preload("Options").
		Limit(limit).
		Offset(offset).
		Order("created_at ASC").
		Find(&questions).Error; err != nil {
		return nil, 0, TranslateGormError(err)
	}

	results := make([]entities.Question, len(questions))
	for i, question := range questions {
		results[i] = *question.ToEntity()
	}

	return results, total, nil
}

func NewPGQuestionRepository(db *gorm.DB) *PGQuestionRepository {
	return &PGQuestionRepository{db: db}
}
