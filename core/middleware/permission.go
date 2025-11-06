package middleware

import (
	"grf/core/permission"

	"github.com/gofiber/fiber/v2"
)

func Check(perm permission.IPermission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := perm.Check(c); err != nil {
			return err
		}
		return c.Next()
	}
}
