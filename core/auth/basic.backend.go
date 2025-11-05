package auth

import (
	"encoding/base64"
	"grf/core/exceptions"
	"grf/domain/auth/model"
	"grf/domain/auth/repository"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BasicAuthBackend struct {
	UserRepo *repository.UserRepository
}

func NewBasicAuthBackend(db *gorm.DB) *BasicAuthBackend {
	return &BasicAuthBackend{UserRepo: repository.NewUserRepository(db)}
}

func (b *BasicAuthBackend) Authenticate(c *fiber.Ctx) (*model.User, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
		return nil, ErrCannotAuthenticate
	}

	base64Credentials := strings.TrimSpace(authHeader[6:])
	creds, err := base64.StdEncoding.DecodeString(base64Credentials)
	if err != nil {
		return nil, exceptions.NewUnauthorized("invalid_basic_credentials", err)
	}

	parts := strings.SplitN(string(creds), ":", 2)
	if len(parts) != 2 {
		return nil, exceptions.NewUnauthorized("invalid_basic_credentials", nil)
	}
	username := parts[0]
	password := parts[1]

	user, err := b.UserRepo.FindUserByEmailOrUsername(username)
	if err != nil || !user.CheckPassword(password) {
		return nil, exceptions.NewUnauthorized("invalid_credentials", err)
	}
	if !user.IsActive {
		return nil, exceptions.NewUnauthorized("inactive_user", nil)
	}

	return &user, nil
}
