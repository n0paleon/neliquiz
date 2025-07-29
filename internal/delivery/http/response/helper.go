package response

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type APIResponse[T any] struct {
	Status string `json:"status"`
	Data   *T     `json:"data"`
	Error  string `json:"error"`
}

// Validate instance (global, singleton, thread-safe)
var Validate = validator.New()

// ParseAndValidate parse and validate user input in one func
// should return nil in handler if error returned
func ParseAndValidate[T any](c *fiber.Ctx) (*T, error) {
	var req T

	if err := c.BodyParser(&req); err != nil {
		_ = ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
		return nil, err
	}

	if err := Validate.Struct(&req); err != nil {
		_ = ErrorResponse(c, fiber.StatusBadRequest, err.Error())
		return nil, err
	}

	return &req, nil
}

// SuccessResponse returns JSON with status: success.
func SuccessResponse[T any](c *fiber.Ctx, data T) error {
	resp := APIResponse[T]{
		Status: "success",
		Data:   &data,
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

// ErrorResponse returns JSON with status: error.
func ErrorResponse(c *fiber.Ctx, statusCode int, errMsg string) error {
	resp := APIResponse[any]{
		Status: "error",
		Error:  errMsg,
	}
	return c.Status(statusCode).JSON(resp)
}
