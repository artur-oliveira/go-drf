package auth

import (
	"grf/core/config"
	"grf/core/exceptions"
	"grf/domain/auth/model"
	"grf/domain/auth/service"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type JWTAuthBackend struct {
	TokenService *service.TokenService
}

func NewJWTAuthBackend(db *gorm.DB, config *config.Config) *JWTAuthBackend {
	return &JWTAuthBackend{
		TokenService: service.NewTokenService(db, config),
	}
}

func (b *JWTAuthBackend) Authenticate(c *fiber.Ctx) (*model.User, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return nil, ErrCannotAuthenticate
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {

		if parts[0] == "Bearer" {
			return nil, exceptions.NewUnauthorized("invalid_jwt_credentials", nil)
		}
		return nil, ErrCannotAuthenticate
	}

	tokenString := parts[1]
	user, err := b.TokenService.ValidateToken(tokenString, "access")
	if err != nil {
		return nil, exceptions.NewUnauthorized(err.Error(), err)
	}

	return user, nil
}
