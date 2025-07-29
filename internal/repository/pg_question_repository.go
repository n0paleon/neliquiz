package repository

import (
	"NeliQuiz/internal/domain/entities"
	"NeliQuiz/internal/errorx"
	"NeliQuiz/internal/repository/schema"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
	if err := r.db.Preload("Options").Where("id = ?", id).First(&result).Error; err != nil {
		logrus.Error(err.Error())
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
	var result schema.Question
	if err := r.db.Preload("Options").Order("RANDOM()").First(&result).Error; err != nil {
		return nil, TranslateGormError(err)
	}
	return result.ToEntity(), nil
}

func (r *PGQuestionRepository) PaginateQuestions(page, limit int) ([]entities.Question, int64, error) {
	var questions []schema.Question
	var total int64

	r.db.Model(&schema.Question{}).Order("created_at ASC").Count(&total)

	offset := (page - 1) * limit
	if err := r.db.Limit(limit).Offset(offset).Find(&questions).Error; err != nil {
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
