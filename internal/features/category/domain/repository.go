package domain

type CategoryRepository interface {
	FindOrCreateBatch(categories []Category) ([]Category, error)
	FindOrCreateCategoryByName(name string) (*Category, error)
	FindCategoryByName(name string) (*Category, error)
	FindAll() ([]Category, error)
	SearchCategoriesByName(query string, limit int) ([]Category, error)
}
