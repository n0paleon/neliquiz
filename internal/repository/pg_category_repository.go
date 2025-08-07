package repository

import (
	"NeliQuiz/internal/domain/entities"
	"NeliQuiz/internal/repository/schema"
	"errors"
	"gorm.io/gorm"
)

type PGCategoryRepository struct {
	db *gorm.DB
}

func (r *PGCategoryRepository) FindOrCreateCategoryByName(name string) (*entities.Category, error) {
	var category schema.Category

	// Cari berdasarkan nama
	if err := r.db.Where("LOWER(name) = LOWER(?)", name).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Jika tidak ditemukan, buat baru
			newCategory := schema.Category{
				Name: name,
			}

			if err := r.db.Create(&newCategory).Error; err != nil {
				return nil, TranslateGormError(err)
			}

			return newCategory.ToEntity(), nil
		}
		return nil, TranslateGormError(err) // error lain (misal: DB down)
	}

	return category.ToEntity(), nil
}

func NewPGCategoryRepository(db *gorm.DB) *PGCategoryRepository {
	return &PGCategoryRepository{
		db: db,
	}
}
