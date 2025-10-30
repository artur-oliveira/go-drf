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

// IPagination Interface genérica de paginador
type IPagination[T any] interface {
	Paginate(c *fiber.Ctx, db *gorm.DB) (*Response[T], error)
}
