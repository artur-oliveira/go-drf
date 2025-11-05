package auth

import (
	"errors"
	"grf/domain/auth/model"

	"github.com/gofiber/fiber/v2"
)

var ErrCannotAuthenticate = errors.New("não é possível autenticar com este backend")

type IAuthBackend interface {
	Authenticate(c *fiber.Ctx) (*model.User, error)
}
