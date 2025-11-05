package exceptions

import (
	"github.com/gofiber/fiber/v2"
)

type AppError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewError(code int, message string, err error) *AppError {
	return &AppError{
		StatusCode: code,
		Message:    message,
		Err:        err,
	}
}

func NewUnauthorized(message string, err error) *AppError {
	return NewError(fiber.StatusUnauthorized, message, err)
}

func NewForbidden(message string, err error) *AppError {
	return NewError(fiber.StatusForbidden, message, err)
}

func NewBadRequest(message string, err error) *AppError {
	return NewError(fiber.StatusBadRequest, message, err)
}

func NewInternal(err error) *AppError {
	return NewError(fiber.StatusInternalServerError, "unexpected_server_error", err)
}
