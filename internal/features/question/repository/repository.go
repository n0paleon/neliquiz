package repository

import (
	categoryRepo "NeliQuiz/internal/features/category/repository"
	questionDomain "NeliQuiz/internal/features/question/domain"
	"NeliQuiz/internal/shared/errorx"
	"NeliQuiz/internal/shared/repoutil"
	"NeliQuiz/internal/shared/strutil"
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

func (r *PGQuestionRepository) Create(q *questionDomain.Question) (*questionDomain.Question, error) {
	model := ToQuestionSchema(q)

	// load existing categories untuk menghindari data terduplikat
	categories := make([]categoryRepo.CategorySchema, len(q.Categories))
	for i, c := range q.Categories {
		categories[i] = categoryRepo.CategorySchema{
			ID:   c.ID,
			Name: c.Name,
		}
	}
	model.Categories = categories

	if err := r.db.Omit("Categories.*").Create(model).Error; err != nil {
		logrus.Error(err.Error())
		return nil, repoutil.TranslateGormError(err)
	}
	return model.ToEntity(), nil
}

func (r *PGQuestionRepository) FindById(id string) (*questionDomain.Question, error) {
	var result QuestionSchema
	if err := r.db.
		Preload("Categories").
		Where("id = ?", id).
		First(&result).Error; err != nil {
		return nil, repoutil.TranslateGormError(err)
	}
	return result.ToEntity(), nil
}

func (r *PGQuestionRepository) DeleteById(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	err := r.db.Where("id = ?", id).Delete(&QuestionSchema{})
	if err.Error != nil {
		return repoutil.TranslateGormError(err.Error)
	}
	if err.RowsAffected < 1 {
		return errorx.NotFound("question not found")
	}

	return nil
}

func (r *PGQuestionRepository) GetRandom() (*questionDomain.Question, error) {
	var questions []QuestionSchema
	if err := r.db.
		Preload("Categories").
		Order("updated_at ASC").
		Limit(getRandomLimitQuery).
		Find(&questions).Error; err != nil {
		return nil, repoutil.TranslateGormError(err)
	}

	if len(questions) < 1 {
		return nil, errorx.NotFound("no question found")
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
	err := r.db.Model(QuestionSchema{}).
		Where("id = ?", id).
		UpdateColumn("hit", gorm.Expr("hit + 1")).
		UpdateColumn("updated_at", time.Now()).
		Error

	return repoutil.TranslateGormError(err)
}

func (r *PGQuestionRepository) PaginateQuestions(page, limit int, sortBy, order string) ([]questionDomain.Question, int64, error) {
	var questions []QuestionSchema
	var total int64

	sortBy, order = repoutil.SanitizeSort(sortBy, order)

	r.db.Model(&QuestionSchema{}).Count(&total)

	offset := (page - 1) * limit
	if err := r.db.
		Preload("Categories").
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sortBy, order)).
		Find(&questions).Error; err != nil {
		return nil, 0, repoutil.TranslateGormError(err)
	}

	results := make([]questionDomain.Question, len(questions))
	for i, question := range questions {
		results[i] = *question.ToEntity()
	}

	return results, total, nil
}

func (r *PGQuestionRepository) PaginateQuestionsByCategory(categoryID string, page, limit int, sortBy, order string) ([]questionDomain.Question, int64, error) {
	var questions []QuestionSchema
	var total int64

	sortBy, order = repoutil.SanitizeSort(sortBy, order)

	query := r.db.Model(&QuestionSchema{}).
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
		return nil, 0, repoutil.TranslateGormError(err)
	}

	results := make([]questionDomain.Question, len(questions))
	for i, question := range questions {
		results[i] = *question.ToEntity()
	}
	return results, total, nil
}

func (r *PGQuestionRepository) Update(q *questionDomain.Question) (*questionDomain.Question, error) {
	var existing QuestionSchema
	tx := r.db.Begin()

	if err := tx.First(&existing, "id = ?", q.ID).Error; err != nil {
		return nil, repoutil.TranslateGormError(err)
	}

	// update kolom biasa
	existing.Content = q.Content
	existing.Hit = q.Hit
	existing.Options = q.Options
	existing.ExplanationURL = q.ExplanationURL
	existing.UpdatedAt = time.Now()
	if err := tx.Save(&existing).Error; err != nil {
		return nil, repoutil.TranslateGormError(err)
	}

	newCategories := make([]categoryRepo.CategorySchema, len(q.Categories))
	for i, cat := range q.Categories {
		newCategories[i] = categoryRepo.CategorySchema{
			ID:   cat.ID,
			Name: cat.Name,
		}
	}

	if err := tx.Model(&QuestionSchema{ID: q.ID}).
		Association("Categories").
		Replace(newCategories); err != nil {
		return nil, repoutil.TranslateGormError(err)
	}

	var finalResult QuestionSchema
	if err := tx.Preload("Categories").First(&finalResult, "id = ?", q.ID).Error; err != nil {
		return nil, repoutil.TranslateGormError(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, repoutil.TranslateGormError(err)
	}

	return finalResult.ToEntity(), nil
}

func (r *PGQuestionRepository) GetRandomByCategoryNames(names []string) (*questionDomain.Question, error) {
	if len(names) == 0 {
		return nil, errorx.BadRequest("category names are empty")
	}

	for i, name := range names {
		names[i] = strutil.NormalizeTitle(name)
	}

	var questions []QuestionSchema
	if err := r.db.
		Joins("JOIN question_categories qc ON qc.question_id = questions.id").
		Joins("JOIN categories c ON c.id = qc.category_id").
		Where("c.name IN ?", names).
		Preload("Categories").
		Order("updated_at ASC").
		Limit(getRandomLimitQuery).
		Find(&questions).Error; err != nil {
		return nil, repoutil.TranslateGormError(err)
	}

	if len(questions) < 1 {
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

func NewQuestionRepository(db *gorm.DB) *PGQuestionRepository {
	return &PGQuestionRepository{db: db}
}
