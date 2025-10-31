package exceptions

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// AppError é nosso erro customizado que os handlers podem retornar
// quando um erro padrão (como GORM) não é suficiente.
type AppError struct {
	StatusCode int    `json:"-"`       // O código de status HTTP (ex: 400, 404)
	Message    string `json:"message"` // A mensagem pública para o JSON
	Err        error  `json:"-"`       // O erro interno (para logging)
}

// Error implementa a interface 'error'
func (e *AppError) Error() string {
	return e.Message
}

// NewError cria um novo AppError (para erros 400 ou 500)
func NewError(code int, message string, err error) *AppError {
	return &AppError{
		StatusCode: code,
		Message:    message,
		Err:        err,
	}
}

func NewBadRequest(message string, err error) *AppError {
	return NewError(400, message, err)
}
func NewNotFound(message string, err error) *AppError {
	return NewError(404, message, err)
}
func NewInternal(err error) *AppError {
	return NewError(500, "Unexpected server error", err)
}

func formatValidationErrors(errs validator.ValidationErrors) string {
	if len(errs) == 0 {
		return "Validation error"
	}

	firstErr := errs[0]

	switch firstErr.Tag() {
	case "required":
		return fmt.Sprintf("The field '%s' is required", firstErr.Field())
	case "email":
		return fmt.Sprintf("The field '%s' must be a valid email", firstErr.Field())
	case "min":
		return fmt.Sprintf("The field '%s' must have at least %s characters", firstErr.Field(), firstErr.Param())
	case "max":
		return fmt.Sprintf("The field '%s' must have at most %s caracteres", firstErr.Field(), firstErr.Param())
	default:
		return fmt.Sprintf("Validation failed for field '%s' (rule: %s)", firstErr.Field(), firstErr.Tag())
	}
}
