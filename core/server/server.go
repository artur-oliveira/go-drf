package server

import (
	"grf/core/config"
	"grf/core/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type App struct {
	FiberApp  *fiber.App
	Config    *config.Config
	DB        *gorm.DB
	Validator *validator.Validate

	I18nMw *middleware.I18NMiddleware
	AuthMw *middleware.AuthenticationMiddleware
	PermMw *middleware.PermissionMiddleware
}
