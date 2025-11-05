package filter

import (
	"grf/core/filterset"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var _ filterset.IFilterSet = (*UserFilterSet)(nil)

type UserFilterSet struct {
	UsernameIContains string
	EmailIExact       string
	IsActive          *bool
	IsStaff           *bool
}

func (f *UserFilterSet) Bind(c *fiber.Ctx) error {
	f.UsernameIContains = c.Query("username__icontains")
	f.EmailIExact = c.Query("email__iexact")

	if val := c.Query("is_active"); val != "" {
		bVal, err := strconv.ParseBool(val)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Parâmetro 'is_active' inválido: use 'true' ou 'false'")
		}
		f.IsActive = &bVal
	}

	if val := c.Query("is_staff"); val != "" {
		bVal, err := strconv.ParseBool(val)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Parâmetro 'is_staff' inválido: use 'true' ou 'false'")
		}
		f.IsStaff = &bVal
	}

	return nil
}

func (f *UserFilterSet) Apply(db *gorm.DB) *gorm.DB {
	query := db
	if f.UsernameIContains != "" {
		query = query.Where("username ILIKE ?", "%"+f.UsernameIContains+"%")
	}
	if f.EmailIExact != "" {
		query = query.Where("LOWER(email) = LOWER(?)", f.EmailIExact)
	}
	if f.IsActive != nil {
		query = query.Where("is_active = ?", *f.IsActive)
	}
	if f.IsStaff != nil {
		query = query.Where("is_staff = ?", *f.IsStaff)
	}
	return query
}
