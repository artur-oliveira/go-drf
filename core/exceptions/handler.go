package exceptions

import (
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GlobalErrorHandler é o manipulador centralizado.
func GlobalErrorHandler(c *fiber.Ctx, err error) error {

	code := fiber.StatusInternalServerError
	message := "Unexpected server error"

	var appErr *AppError
	var fiberErr *fiber.Error
	var validationErrs validator.ValidationErrors

	if errors.As(err, &appErr) {
		code = appErr.StatusCode
		message = appErr.Message
		if appErr.Err != nil && code >= 500 {
			log.Printf("Internal Error: %v", appErr.Err)
		}
	} else if errors.As(err, &fiberErr) {
		code = fiberErr.Code
		message = fiberErr.Message
	} else if errors.As(err, &validationErrs) {
		code = fiber.StatusUnprocessableEntity // 422
		message = formatValidationErrors(validationErrs)
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		code = fiber.StatusNotFound // 404
		message = "Not found"
	} else {
		log.Printf("Unexpected server error: %v", err)
	}

	type StandardErrorResponse struct {
		Error string `json:"error"`
	}

	return c.Status(code).JSON(StandardErrorResponse{
		Error: message,
	})
}
