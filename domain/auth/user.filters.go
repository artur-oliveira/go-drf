package auth

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserFilterSet struct {
	Username string `query:"username"`
	Email    string `query:"email"`
	IsAdmin  *bool  `query:"is_admin"` // Ponteiro para diferenciar 'false' de 'não fornecido'
}

func (f *UserFilterSet) Bind(c *fiber.Ctx) error {
	if err := c.QueryParser(f); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Parâmetros de filtro inválidos")
	}
	return nil
}

func (f *UserFilterSet) Apply(db *gorm.DB) *gorm.DB {
	query := db

	if f.Username != "" {
		// Use ILIKE para case-insensitive, ou ajuste conforme o dialeto do DB
		query = query.Where("username ILIKE ?", "%"+f.Username+"%")
	}
	if f.Email != "" {
		query = query.Where("email = ?", f.Email)
	}
	if f.IsAdmin != nil {
		query = query.Where("is_admin = ?", *f.IsAdmin)
	}

	return query
}
