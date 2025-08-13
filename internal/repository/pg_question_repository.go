package repository

import (
	"NeliQuiz/internal/domain/entities"
	"NeliQuiz/internal/errorx"
	"NeliQuiz/internal/repository/schema"
	"NeliQuiz/pkg/utils"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math/big"
	"time"
)

type PGQuestionRepository struct {
	db *gorm.DB
}

var getRandomLimitQuery = 10

func (r *PGQuestionRepository) Create(q *entities.Question) (*entities.Question, error) {
	model := schema.ToQuestionSchema(q)

	// load existing categories untuk menghindari data terduplikat
	categories := make([]schema.Category, len(q.Categories))
	for i, category := range q.Categories {
		categories[i] = schema.Category{
			Schema: schema.Schema{ID: category.ID},
			Name:   category.Name,
		}
	}
	model.Categories = categories

	if err := r.db.Omit("Categories.*").Create(model).Error; err != nil {
		logrus.Error(err.Error())
		return nil, TranslateGormError(err)
	}
	return model.ToEntity(), nil
}

func (r *PGQuestionRepository) FindById(id string) (*entities.Question, error) {
	var result schema.Question
	if err := r.db.
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
		Preload("Categories").
		Order("updated_at ASC").
		Limit(getRandomLimitQuery).
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
		UpdateColumn("updated_at", time.Now()).
		Error

	return TranslateGormError(err)
}

func (r *PGQuestionRepository) PaginateQuestions(page, limit int, sortBy, order string) ([]entities.Question, int64, error) {
	var questions []schema.Question
	var total int64

	sortBy, order = sanitizeSort(sortBy, order)

	r.db.Model(&schema.Question{}).Count(&total)

	offset := (page - 1) * limit
	if err := r.db.
		Preload("Categories").
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sortBy, order)).
		Find(&questions).Error; err != nil {
		return nil, 0, TranslateGormError(err)
	}

	results := make([]entities.Question, len(questions))
	for i, question := range questions {
		results[i] = *question.ToEntity()
	}

	return results, total, nil
}

func (r *PGQuestionRepository) PaginateQuestionsByCategory(categoryID string, page, limit int, sortBy, order string) ([]entities.Question, int64, error) {
	var questions []schema.Question
	var total int64

	sortBy, order = sanitizeSort(sortBy, order)

	query := r.db.Model(&schema.Question{}).
		Joins("JOIN question_categories qc ON qc.question_id = questions.id").
		Where("qc.category_id = ?", categoryID).
		Preload("Categories").
		Group("questions.id")

	query.Count(&total)

	offset := (page - 1) * limit
	if err := query.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sortBy, order)).
		Find(&questions).Error; err != nil {
		return nil, 0, TranslateGormError(err)
	}

	results := make([]entities.Question, len(questions))
	for i, question := range questions {
		results[i] = *question.ToEntity()
	}
	return results, total, nil
}

func (r *PGQuestionRepository) Update(q *entities.Question) (*entities.Question, error) {
	var existing schema.Question
	tx := r.db.Begin()

	if err := tx.First(&existing, "id = ?", q.ID).Error; err != nil {
		return nil, TranslateGormError(err)
	}

	// update kolom biasa
	existing.Content = q.Content
	existing.Hit = q.Hit
	existing.Options = q.Options
	existing.ExplanationURL = q.ExplanationURL
	existing.UpdatedAt = time.Now()
	if err := tx.Model(&schema.Question{ID: q.ID}).Updates(existing).Error; err != nil {
		return nil, TranslateGormError(err)
	}

	newCategories := make([]schema.Category, len(q.Categories))
	for i, cat := range q.Categories {
		newCategories[i] = schema.Category{
			Schema: schema.Schema{ID: cat.ID},
			Name:   cat.Name,
		}
	}

	if err := tx.Model(&schema.Question{ID: q.ID}).
		Association("Categories").
		Replace(newCategories); err != nil {
		return nil, TranslateGormError(err)
	}

	var finalResult schema.Question
	if err := tx.Preload("Categories").First(&finalResult, "id = ?", q.ID).Error; err != nil {
		return nil, TranslateGormError(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, TranslateGormError(err)
	}

	return finalResult.ToEntity(), nil
}

func (r *PGQuestionRepository) GetRandomByCategoryNames(names []string) (*entities.Question, error) {
	if len(names) == 0 {
		return nil, errorx.BadRequest("category names are empty")
	}

	for i, name := range names {
		names[i] = utils.NormalizeTitle(name)
	}

	var questions []schema.Question
	if err := r.db.
		Joins("JOIN question_categories qc ON qc.question_id = questions.id").
		Joins("JOIN categories c ON c.id = qc.category_id").
		Where("c.name IN ?", names).
		Preload("Categories").
		Order("updated_at ASC").
		Limit(getRandomLimitQuery).
		Find(&questions).Error; err != nil {
		return nil, TranslateGormError(err)
	}

	if len(questions) == 0 {
		return nil, errorx.NotFound("no questions found for given category names")
	}

	n := big.NewInt(int64(len(questions)))
	i, _ := rand.Int(rand.Reader, n)
	selected := questions[i.Int64()]

	go func() {
		_ = r.updateHit(selected.ID)
	}()

	return selected.ToEntity(), nil
}

func NewPGQuestionRepository(db *gorm.DB) *PGQuestionRepository {
	return &PGQuestionRepository{db: db}
}
