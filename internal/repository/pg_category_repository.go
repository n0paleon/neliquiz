package repository

import (
	"NeliQuiz/internal/domain/entities"
	"NeliQuiz/internal/repository/schema"
	"NeliQuiz/pkg/utils"
	"gorm.io/gorm"
)

type PGCategoryRepository struct {
	db *gorm.DB
}

func (r *PGCategoryRepository) FindOrCreateBatch(categories []entities.Category) ([]entities.Category, error) {
	var results []entities.Category

	for _, inputCategory := range categories {
		var category schema.Category

		normalizedName := utils.NormalizeTitle(inputCategory.Name)
		err := r.db.Where("LOWER(name) = LOWER(?)", normalizedName).
			FirstOrCreate(&category, schema.Category{Name: normalizedName}).Error

		if err != nil {
			return nil, TranslateGormError(err)
		}

		results = append(results, *category.ToEntity())
	}

	return results, nil
}

func (r *PGCategoryRepository) FindOrCreateCategoryByName(name string) (*entities.Category, error) {
	var category schema.Category

	normalizedName := utils.NormalizeTitle(name)
	if err := r.db.
		Where("LOWER(name) = LOWER(?)", normalizedName).
		FirstOrCreate(&category, schema.Category{Name: normalizedName}).Error; err != nil {
		return nil, TranslateGormError(err)
	}

	return category.ToEntity(), nil
}

func (r *PGCategoryRepository) FindCategoryByName(name string) (*entities.Category, error) {
	var category schema.Category

	if err := r.db.Where("LOWER(name) = LOWER(?)", name).First(&category).Error; err != nil {
		return nil, TranslateGormError(err)
	}

	return category.ToEntity(), nil
}

func (r *PGCategoryRepository) FindAll() ([]entities.Category, error) {
	var categories []schema.Category

	if err := r.db.Find(&categories).Error; err != nil {
		return nil, TranslateGormError(err)
	}

	var results []entities.Category
	for _, category := range categories {
		results = append(results, *category.ToEntity())
	}

	return results, nil
}

func (r *PGCategoryRepository) SearchCategoriesByName(query string, limit int) ([]entities.Category, error) {
	if query == "" {
		return nil, nil
	}

	var categories []schema.Category
	if err := r.db.
		Where("name ILIKE ?", "%"+query+"%").
		Limit(limit).
		Order("name ASC").
		Find(&categories).Error; err != nil {
		return nil, TranslateGormError(err)
	}

	results := make([]entities.Category, len(categories))
	for i, cat := range categories {
		results[i] = *cat.ToEntity()
	}

	return results, nil
}

func NewPGCategoryRepository(db *gorm.DB) *PGCategoryRepository {
	return &PGCategoryRepository{
		db: db,
	}
}
