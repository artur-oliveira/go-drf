package grf

import (
	"grf/config"
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

	AuthMw *middleware.AuthMiddleware
	PermMw *middleware.PermissionMiddleware
}
