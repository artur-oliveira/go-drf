package middleware

import (
	"grf/core/i18n"

	"github.com/gofiber/fiber/v2"
)

type I18NMiddleware struct {
	service *i18n.Service
}

func NewI18NMiddleware(i18nSvc *i18n.Service) *I18NMiddleware {
	return &I18NMiddleware{service: i18nSvc}
}

func (m *I18NMiddleware) UseMiddleWare(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		lang := c.Get("Accept-Language", "en-US")

		localizer := m.service.GetLocalizer(lang)

		c.Locals("localizer", localizer)
		return c.Next()
	})
}
