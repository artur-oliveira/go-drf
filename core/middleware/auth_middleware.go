package middleware

import (
	"grf/config"
	"grf/core/exceptions"
	"grf/domain/auth"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// AuthMiddleware fornece o middleware de autenticação
type AuthMiddleware struct {
	TokenService *auth.TokenService
}

// NewAuthMiddleware cria a factory do middleware
func NewAuthMiddleware(db *gorm.DB, config *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		TokenService: auth.NewTokenService(db, config),
	}
}

// RequireAuth é o handler de middleware que protege uma rota
func (m *AuthMiddleware) RequireAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return exceptions.NewError(401, "Token de autenticação não fornecido", nil)
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return exceptions.NewError(401, "Formato de token inválido (esperado: Bearer <token>)", nil)
	}

	tokenString := parts[1]

	// Validar o Access Token
	user, err := m.TokenService.ValidateToken(tokenString, "access")
	if err != nil {
		// Retorna 401 com a mensagem de erro do serviço (ex: "token expirado")
		return exceptions.NewError(401, err.Error(), err)
	}

	// Anexar o usuário ao contexto para uso posterior (ex: middleware de permissão)
	c.Locals("user", user)
	return c.Next()
}
