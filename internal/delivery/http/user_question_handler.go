package http

import (
	"NeliQuiz/internal/delivery/http/dto"
	"NeliQuiz/internal/delivery/http/response"
	"NeliQuiz/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type UserQuestionHandler struct {
	questionUseCase domain.UserQuestionUseCase
}

func (h *UserQuestionHandler) GetRandomQuestion(c *fiber.Ctx) error {
	question, err := h.questionUseCase.GetRandomQuestion()
	if err != nil {
		return response.ErrorResponse(c, 500, err.Error())
	}

	result := dto.EntityToGetRandomQuestionResponse(question)
	return response.SuccessResponse(c, result)
}

func (h *UserQuestionHandler) CheckAnswer(c *fiber.Ctx) error {
	questionID := c.Params("id", "")
	if questionID == "" {
		return response.ErrorResponse(c, 400, "questionID is required!")
	}

	input, err := response.ParseAndValidate[dto.PostVerifyAnswerRequest](c)
	if err != nil {
		return nil
	}

	correct, optionData, explanationURL, err := h.questionUseCase.CheckAnswer(questionID, input.SelectedOptionID)
	if err != nil {
		return err
	}

	responseData := dto.PostVerifyAnswerResponse{
		Correct:        correct,
		ExplanationURL: explanationURL,
	}
	if correct {
		responseData.CorrectOption = dto.QuizOption{
			OptionID: optionData.ID,
			Content:  optionData.Content,
		}
	}

	return response.SuccessResponse(c, &responseData)
}

func NewUserQuestionHandler(questionUseCase domain.UserQuestionUseCase) *UserQuestionHandler {
	return &UserQuestionHandler{
		questionUseCase: questionUseCase,
	}
}
