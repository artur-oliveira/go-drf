package pagination

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Response[T any] struct {
	Results []T   `json:"results"`
	HasNext bool  `json:"has_next"`
	Count   *uint `json:"count"`
}

type IPagination[T any] interface {
	Bind(c *fiber.Ctx) error

	Paginate(db *gorm.DB) (*Response[T], error)
}
