package auth

import (
	"grf/core/filterset"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var _ filterset.IFilterSet = (*PermissionFilterSet)(nil)

type PermissionFilterSet struct {
	Module string
	Action string
}

func (f *PermissionFilterSet) Bind(c *fiber.Ctx) error {
	f.Module = c.Query("module")
	f.Action = c.Query("action")
	return nil
}

func (f *PermissionFilterSet) Apply(db *gorm.DB) *gorm.DB {
	query := db
	if f.Module != "" {
		query = query.Where("module = ?", f.Module)
	}
	if f.Action != "" {
		query = query.Where("action = ?", f.Action)
	}
	return query
}
