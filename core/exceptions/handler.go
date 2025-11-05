package exceptions

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gorm.io/gorm"
)

func getLocalizer(c *fiber.Ctx) *i18n.Localizer {
	loc := c.Locals("localizer")
	if loc == nil {
		return i18n.NewLocalizer(nil, "en-US")
	}
	return loc.(*i18n.Localizer)
}

func GlobalErrorHandler(c *fiber.Ctx, err error) error {

	localizer := getLocalizer(c)

	code := fiber.StatusInternalServerError
	messageKey := "error_internal"

	var templateData map[string]interface{}
	var appErr *AppError
	var validationErrs validator.ValidationErrors

	if errors.As(err, &appErr) {
		code = appErr.StatusCode
		messageKey = appErr.Message

		if appErr.Err != nil && code >= 500 {
			log.Printf("Internal Error (AppError): %v", appErr.Err)
		}

	} else if errors.As(err, &validationErrs) {
		code = fiber.StatusUnprocessableEntity
		messageKey = "error_validation"

		fieldErrors := formatValidationErrors(validationErrs, localizer)

		translatedMessage, e := localizer.Localize(&i18n.LocalizeConfig{
			MessageID: messageKey,
		})
		if e != nil {
			log.Printf("Failed to translate key '%s': %v", messageKey, e)
			translatedMessage = messageKey
		}

		return c.Status(code).JSON(fiber.Map{
			"error":  translatedMessage,
			"fields": fieldErrors,
		})
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		code = fiber.StatusNotFound
		messageKey = "error_not_found"

	} else {
		log.Printf("Unexpected error: %v", err)
	}

	translatedMessage, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageKey,
		TemplateData: templateData,
	})

	if err != nil {
		log.Printf("Failed to translate key '%s': %v", messageKey, err)
		translatedMessage = messageKey
	}

	return c.Status(code).JSON(fiber.Map{
		"error": translatedMessage,
	})
}

func formatValidationErrors(errs validator.ValidationErrors, localizer *i18n.Localizer) map[string]string {
	fields := make(map[string]string)

	for _, err := range errs {
		key := "validation_" + err.Tag()

		msg, e := localizer.Localize(&i18n.LocalizeConfig{
			MessageID: key,
			TemplateData: map[string]interface{}{
				"Field": err.Field(),
				"Param": err.Param(),
			},
		})

		if e != nil {
			msg = fmt.Sprintf("Falha na validação da regra: '%s'", err.Tag())
			log.Printf("Aviso I18N: Chave de validação não encontrada: '%s'", key)
		}

		fields[err.Field()] = msg
	}
	return fields
}
