package errorx

import (
	"fmt"
	"net/http"
)

type AppError interface {
	error
	Code() int
	Message() string
	Cause() error
}

type appError struct {
	code    int
	message string
	cause   error
}

func (e *appError) Error() string {
	return fmt.Sprintf("%s: %v", e.message, e.cause)
}

func (e *appError) Code() int {
	return e.code
}

func (e *appError) Message() string {
	return e.message
}

func (e *appError) Cause() error {
	return e.cause
}

func NewAppError(code int, msg string, cause error) AppError {
	return &appError{
		code:    code,
		message: msg,
		cause:   cause,
	}
}

func NotFound(msg string) AppError {
	return NewAppError(http.StatusNotFound, msg, nil)
}

func Conflict(msg string) AppError {
	return NewAppError(http.StatusConflict, msg, nil)
}

func BadRequest(msg string) AppError {
	return NewAppError(http.StatusBadRequest, msg, nil)
}

func InternalError(cause error) AppError {
	return NewAppError(http.StatusInternalServerError, cause.Error(), cause)
}
