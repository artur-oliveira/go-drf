package middleware

import (
	"fmt"
	"grf/core/exceptions"
	"grf/domain/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PermissionMiddleware struct {
	DB *gorm.DB
}

func NewPermissionMiddleware(db *gorm.DB) *PermissionMiddleware {
	return &PermissionMiddleware{DB: db}
}

func (m *PermissionMiddleware) RequirePerm(module, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userLocal := c.Locals("user")
		if userLocal == nil {
			return exceptions.NewError(401, "Usuário não autenticado", nil)
		}

		user, ok := userLocal.(*auth.User)
		if !ok {
			return exceptions.NewInternal(fmt.Errorf("c.Locals(\"user\") não é do tipo *auth.User"))
		}

		if !user.HasPerm(m.DB, module, action) {
			msg := fmt.Sprintf("Permissão necessária: %s.%s", module, action)
			return exceptions.NewError(403, msg, nil)
		}

		return c.Next()
	}
}
