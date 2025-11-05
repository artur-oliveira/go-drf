package filter

import (
	"grf/core/filterset"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var _ filterset.IFilterSet = (*GroupFilterSet)(nil)

type GroupFilterSet struct {
	NameIContains string
}

func (f *GroupFilterSet) Bind(c *fiber.Ctx) error {
	f.NameIContains = c.Query("name__icontains")
	return nil
}

func (f *GroupFilterSet) Apply(db *gorm.DB) *gorm.DB {
	query := db
	if f.NameIContains != "" {
		query = query.Where("name ILIKE ?", "%"+f.NameIContains+"%")
	}
	return query
}
