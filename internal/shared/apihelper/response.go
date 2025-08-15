package apihelper

import "github.com/gofiber/fiber/v2"

type APIResponse[T any] struct {
	Status  string `json:"status"`
	Data    *T     `json:"data"`
	Message string `json:"error"`
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
		Status:  "error",
		Message: errMsg,
	}
	return c.Status(statusCode).JSON(resp)
}
