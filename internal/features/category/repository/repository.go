package repository

import (
	"NeliQuiz/internal/features/category/domain"
	"NeliQuiz/internal/shared/repoutil"
	"NeliQuiz/internal/shared/strutil"
	"gorm.io/gorm"
)

type PGCategoryRepository struct {
	db *gorm.DB
}

func (r *PGCategoryRepository) FindOrCreateBatch(categories []domain.Category) ([]domain.Category, error) {
	var results []domain.Category

	for _, inputCategory := range categories {
		var category CategorySchema

		normalizedName := strutil.NormalizeTitle(inputCategory.Name)
		err := r.db.Where("LOWER(name) = LOWER(?)", normalizedName).
			FirstOrCreate(&category, CategorySchema{Name: normalizedName}).Error

		if err != nil {
			return nil, repoutil.TranslateGormError(err)
		}

		results = append(results, *category.ToEntity())
	}

	return results, nil
}

func (r *PGCategoryRepository) FindOrCreateCategoryByName(name string) (*domain.Category, error) {
	var category CategorySchema

	normalizedName := strutil.NormalizeTitle(name)
	if err := r.db.
		Where("LOWER(name) = LOWER(?)", normalizedName).
		FirstOrCreate(&category, CategorySchema{Name: normalizedName}).Error; err != nil {
		return nil, repoutil.TranslateGormError(err)
	}

	return category.ToEntity(), nil
}

func (r *PGCategoryRepository) FindCategoryByName(name string) (*domain.Category, error) {
	var category CategorySchema

	if err := r.db.Where("LOWER(name) = LOWER(?)", name).First(&category).Error; err != nil {
		return nil, repoutil.TranslateGormError(err)
	}

	return category.ToEntity(), nil
}

func (r *PGCategoryRepository) FindAll() ([]domain.Category, error) {
	var categories []CategorySchema

	if err := r.db.Find(&categories).Error; err != nil {
		return nil, repoutil.TranslateGormError(err)
	}

	var results []domain.Category
	for _, category := range categories {
		results = append(results, *category.ToEntity())
	}

	return results, nil
}

func (r *PGCategoryRepository) SearchCategoriesByName(query string, limit int) ([]domain.Category, error) {
	if query == "" {
		return nil, nil
	}

	var categories []CategorySchema
	if err := r.db.
		Where("name ILIKE ?", "%"+query+"%").
		Limit(limit).
		Order("name ASC").
		Find(&categories).Error; err != nil {
		return nil, repoutil.TranslateGormError(err)
	}

	results := make([]domain.Category, len(categories))
	for i, cat := range categories {
		results[i] = *cat.ToEntity()
	}

	return results, nil
}

func NewCategoryRepository(db *gorm.DB) *PGCategoryRepository {
	return &PGCategoryRepository{
		db: db,
	}
}
