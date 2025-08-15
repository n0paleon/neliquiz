package delivery

import (
	"NeliQuiz/internal/features/question/delivery/dto"
	questionDomain "NeliQuiz/internal/features/question/domain"
	"NeliQuiz/internal/shared/apihelper"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

type QuestionHandler struct {
	questionUseCase questionDomain.QuestionUseCase
}

func (h *QuestionHandler) CreateQuestion(c *fiber.Ctx) error {
	request, err := apihelper.ParseAndValidate[dto.CreateQuestionRequest](c)
	if err != nil {
		return nil
	}

	payload := request.ToDomain()
	result, err := h.questionUseCase.CreateQuestion(payload)
	if err != nil {
		return err
	}

	return apihelper.SuccessResponse(c, fiber.Map{
		"question_id": result.ID,
	})
}

func (h *QuestionHandler) GetListQuestion(c *fiber.Ctx) error {
	pageQuery := c.Query("page", "1")
	limitQuery := c.Query("limit", "10")
	category := c.Query("category", "")
	sortBy := c.Query("sort", "created_at")
	order := c.Query("order", "desc")

	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		return apihelper.ErrorResponse(c, fiber.StatusBadRequest, "invalid page number")
	}
	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		return apihelper.ErrorResponse(c, fiber.StatusBadRequest, "invalid page limit number")
	}

	results, total, err := h.questionUseCase.GetListQuestions(category, page, limit, order, sortBy)
	if err != nil {
		return err
	}

	return apihelper.SuccessResponse(c, fiber.Map{
		"questions":   results,
		"total_count": total,
	})
}

func (h *QuestionHandler) DeleteQuestion(c *fiber.Ctx) error {
	questionID := c.Params("id", "")
	if questionID == "" {
		return apihelper.ErrorResponse(c, fiber.StatusBadRequest, "invalid question id")
	}

	if err := h.questionUseCase.DeleteQuestion(questionID); err != nil {
		return err
	}

	return apihelper.SuccessResponse(c, "question deleted successfully!")
}

func (h *QuestionHandler) GetQuestionDetail(c *fiber.Ctx) error {
	questionID := c.Params("id", "")
	if questionID == "" {
		return apihelper.ErrorResponse(c, fiber.StatusBadRequest, "invalid question id")
	}

	question, err := h.questionUseCase.GetQuestionDetail(questionID)
	if err != nil {
		return err
	}

	return apihelper.SuccessResponse(c, question)
}

func (h *QuestionHandler) UpdateQuestionDetail(c *fiber.Ctx) error {
	questionID := c.Params("id", "")
	if questionID == "" {
		return apihelper.ErrorResponse(c, fiber.StatusBadRequest, "invalid question id")
	}

	request, err := apihelper.ParseAndValidate[dto.UpdateQuestionDetailRequest](c)
	if err != nil {
		return nil
	}

	data := request.ToDomain()
	data.ID = questionID

	newData, err := h.questionUseCase.UpdateQuestion(data)
	if err != nil {
		return err
	}

	return apihelper.SuccessResponse(c, newData)
}

func (h *QuestionHandler) GetRandomQuestion(c *fiber.Ctx) error {
	category := c.Query("category", "")
	var categories []string

	if category != "" {
		categories = strings.Split(strings.TrimSpace(category), ",")
	}

	question, err := h.questionUseCase.GetRandomQuestion(categories...)
	if err != nil {
		return err
	}

	result := dto.ToGetRandomQuestionResponse(question)
	return apihelper.SuccessResponse(c, result)
}

func (h *QuestionHandler) PostVerifyAnswer(c *fiber.Ctx) error {
	questionID := c.Params("id", "")
	if questionID == "" {
		return apihelper.ErrorResponse(c, 400, "questionID is required!")
	}

	request, err := apihelper.ParseAndValidate[dto.PostVerifyAnswerRequest](c)
	if err != nil {
		return nil
	}

	IsCorrect, CorrectOption, explanationURL, err := h.questionUseCase.CheckAnswer(questionID, request.SelectedOptionID)
	if err != nil {
		return err
	}

	responseData := dto.PostVerifyAnswerResponse{
		IsCorrect: IsCorrect,
		CorrectOption: dto.PostVerifyAnswerOption{
			OptionID: CorrectOption.ID,
			Content:  CorrectOption.Content,
		},
		ExplanationURL: explanationURL,
	}

	return apihelper.SuccessResponse(c, &responseData)
}

func NewQuestionHandler(questionUseCase questionDomain.QuestionUseCase) *QuestionHandler {
	return &QuestionHandler{
		questionUseCase: questionUseCase,
	}
}
