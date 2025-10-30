package filterset

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// IFilterSet Base Interface for all filters
type IFilterSet interface {
	// Bind preenche a struct com os query params
	Bind(c *fiber.Ctx) error
	// Apply aplica os filtros à query GORM
	Apply(db *gorm.DB) *gorm.DB
}
