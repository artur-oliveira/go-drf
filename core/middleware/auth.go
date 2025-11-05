package middleware

import (
	"errors"
	"grf/core/auth"
	"grf/core/exceptions"

	"github.com/gofiber/fiber/v2"
)

type AuthenticationMiddleware struct {
	Backends []auth.IAuthBackend
}

func NewAuthenticationMiddleware(backends ...auth.IAuthBackend) *AuthenticationMiddleware {
	if len(backends) == 0 {
		panic("AuthenticationMiddleware: pelo menos um backend é necessário")
	}
	return &AuthenticationMiddleware{
		Backends: backends,
	}
}

func (m *AuthenticationMiddleware) RequireAuth(c *fiber.Ctx) error {

	for _, backend := range m.Backends {
		user, err := backend.Authenticate(c)

		if err != nil {
			if errors.Is(err, auth.ErrCannotAuthenticate) {
				continue
			}
			return err
		}

		if user != nil {
			c.Locals("user", user)
			return c.Next()
		}
	}

	return exceptions.NewUnauthorized("auth_invalid_or_not_provided", nil)
}
