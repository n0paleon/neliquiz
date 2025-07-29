package http

import (
	"NeliQuiz/internal/delivery/http/dto"
	"NeliQuiz/internal/delivery/http/response"
	"NeliQuiz/internal/domain"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type AdminQuestionHandler struct {
	questionUseCase domain.AdminQuestionUseCase
}

func (h *AdminQuestionHandler) PostCreateQuestion(c *fiber.Ctx) error {
	input, err := response.ParseAndValidate[dto.CreateQuestionRequest](c)
	if err != nil {
		return nil
	}

	question, err := input.ToEntity()
	if err != nil {
		return response.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.questionUseCase.CreateQuestion(question); err != nil {
		return response.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, "question saved!")
}

func (h *AdminQuestionHandler) GetListQuestions(c *fiber.Ctx) error {
	pageQuery := c.Query("page", "1")
	limitQuery := c.Query("limit", "10")

	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		return response.ErrorResponse(c, fiber.StatusBadRequest, "invalid page number")
	}
	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		return response.ErrorResponse(c, fiber.StatusBadRequest, "invalid page limit number")
	}

	results, total, err := h.questionUseCase.GetListQuestions(page, limit)
	if err != nil {
		return response.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	questions := make([]dto.GetListQuestionsResponse, len(results))
	for i, result := range results {
		questions[i] = dto.GetListQuestionsResponse{
			ID:        result.ID,
			Content:   result.Content,
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
		}
	}

	return response.SuccessResponse(c, fiber.Map{
		"questions":   questions,
		"total_count": total,
	})
}

func (h *AdminQuestionHandler) DeleteQuestion(c *fiber.Ctx) error {
	input, err := response.ParseAndValidate[dto.PostDeleteQuestionRequest](c)
	if err != nil {
		return nil
	}

	if err := h.questionUseCase.DeleteQuestion(input.QuestionID); err != nil {
		return err
	}

	return response.SuccessResponse(c, "question deleted successfully!")
}

func NewAdminQuestionHandler(questionUseCase domain.AdminQuestionUseCase) *AdminQuestionHandler {
	return &AdminQuestionHandler{
		questionUseCase: questionUseCase,
	}
}
