package permission

import (
	"errors"
	"grf/core/exceptions"
	"grf/core/models"

	"github.com/gofiber/fiber/v2"
)

func GetUser(c *fiber.Ctx) (models.IUser, error) {
	userLocal := c.Locals("user")
	if userLocal == nil {
		return nil, exceptions.NewUnauthorized("auth_invalid_or_not_provided", nil)
	}

	user, ok := userLocal.(models.IUser)
	if !ok {
		return nil, exceptions.NewInternal(errors.New("c.Locals(\"user\") não implementa models.IUser"))
	}

	if !user.Active() {
		return nil, exceptions.NewUnauthorized("inactive_user", nil)
	}

	return user, nil
}

func IsReadOnlyMethod(method string) bool {
	switch method {
	case fiber.MethodGet, fiber.MethodHead, fiber.MethodOptions:
		return true
	default:
		return false
	}
}

func getActionForMethod(method string) (string, error) {
	switch method {
	case fiber.MethodGet, fiber.MethodHead, fiber.MethodOptions:
		return "view", nil
	case fiber.MethodPost:
		return "add", nil
	case fiber.MethodPut, fiber.MethodPatch:
		return "change", nil
	case fiber.MethodDelete:
		return "delete", nil
	default:
		return "", errors.New("método HTTP não suportado: " + method)
	}
}
