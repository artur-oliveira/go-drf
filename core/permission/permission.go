package permission

import (
	"github.com/gofiber/fiber/v2"
)

type IPermission interface {
	Check(c *fiber.Ctx) error
}
