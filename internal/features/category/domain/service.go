package domain

type CategoryUseCase interface {
	GetListCategories(query string, limit int) ([]Category, error)
	FindOrCreateBatch(categories []Category) ([]Category, error)
	FindCategoryByName(name string) (*Category, error)
}
