package delivery

import (
	"NeliQuiz/internal/features/category/domain"
	"NeliQuiz/internal/shared/apihelper"
	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	categoryUseCase domain.CategoryUseCase
}

func (h *CategoryHandler) GetListCategories(c *fiber.Ctx) error {
	query := c.Query("q", "")
	limit := 10

	categories, err := h.categoryUseCase.GetListCategories(query, limit)
	if err != nil {
		return err
	}

	return apihelper.SuccessResponse(c, categories)
}

func NewCategoryHandler(categoryUseCase domain.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		categoryUseCase: categoryUseCase,
	}
}
