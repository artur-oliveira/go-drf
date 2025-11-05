package filterset

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type IFilterSet interface {
	Bind(c *fiber.Ctx) error
	Apply(db *gorm.DB) *gorm.DB
}
