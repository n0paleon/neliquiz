package http

import (
	"NeliQuiz/internal/delivery/http/response"
	"NeliQuiz/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type AdminCategoryHandler struct {
	categoryUseCase domain.AdminCategoryUseCase
}

func (h *AdminCategoryHandler) GetListCategories(c *fiber.Ctx) error {
	query := c.Query("q", "")
	limit := 10

	categories, err := h.categoryUseCase.GetListCategories(query, limit)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, categories)
}

func NewAdminCategoryHandler(categoryUseCase domain.AdminCategoryUseCase) *AdminCategoryHandler {
	return &AdminCategoryHandler{
		categoryUseCase: categoryUseCase,
	}
}
