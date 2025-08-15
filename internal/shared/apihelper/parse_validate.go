package apihelper

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

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
