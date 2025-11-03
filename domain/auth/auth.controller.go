package auth

import (
	"errors"
	"grf/config"
	"grf/core/exceptions"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Controller struct {
	DB           *gorm.DB
	Validator    *validator.Validate
	TokenService *TokenService
}

func NewAuthController(
	db *gorm.DB,
	config *config.Config,
	validate *validator.Validate,
) *Controller {
	return &Controller{
		DB:           db,
		Validator:    validate,
		TokenService: NewTokenService(db, config),
	}
}

func (ac *Controller) ObtainToken(c *fiber.Ctx) error {
	var input ObtainTokenDTO
	if err := c.BodyParser(&input); err != nil {
		return exceptions.NewBadRequest("Payload inválido", err)
	}
	if err := ac.Validator.Struct(input); err != nil {
		return err
	}

	var user User
	if err := ac.DB.Where("username = ? OR email = ?", input.Login, input.Login).First(&user).Error; err != nil {
		return exceptions.NewError(401, "Credenciais inválidas", err)
	}

	if !user.IsActive {
		return exceptions.NewError(401, "Usuário inativo", nil)
	}

	if !user.CheckPassword(input.Password) {
		return exceptions.NewError(401, "Credenciais inválidas", nil)
	}

	access, refresh, err := ac.TokenService.GenerateTokenPair(&user)
	if err != nil {
		return exceptions.NewInternal(err)
	}

	return c.JSON(TokenResponseDTO{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}
func (ac *Controller) ObtainTokenRefresh(c *fiber.Ctx) error {
	var input RefreshTokenDTO
	if err := c.BodyParser(&input); err != nil {
		return exceptions.NewBadRequest("Payload inválido", err)
	}
	if err := ac.Validator.Struct(input); err != nil {
		return err // Handler global
	}

	// 1. Validar o Refresh Token
	user, err := ac.TokenService.ValidateToken(input.Refresh, "refresh")
	if err != nil {
		// Erro 401 se o refresh token for inválido ou expirado
		return exceptions.NewError(fiber.StatusUnauthorized, err.Error(), err)
	}

	// 2. Gerar um novo par de tokens
	access, refresh, err := ac.TokenService.GenerateTokenPair(user)
	if err != nil {
		return exceptions.NewInternal(err)
	}

	return c.JSON(TokenResponseDTO{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

// ChangePassword (POST /auth/change-password)
// Requer autenticação (AuthMiddleware)
func (ac *Controller) ChangePassword(c *fiber.Ctx) error {
	var input ChangePasswordDTO
	if err := c.BodyParser(&input); err != nil {
		return exceptions.NewBadRequest("Payload inválido", err)
	}
	if err := ac.Validator.Struct(input); err != nil {
		return err // Handler global
	}

	// 1. Obter usuário do middleware
	user, ok := c.Locals("user").(*User)
	if !ok {
		return exceptions.NewInternal(errors.New("c.Locals(\"user\") não encontrado"))
	}

	// 2. Verificar senha antiga
	if !user.CheckPassword(input.OldPassword) {
		return exceptions.NewBadRequest("Senha antiga incorreta", nil)
	}

	// 3. Definir nova senha
	if err := user.SetPassword(input.NewPassword); err != nil {
		return exceptions.NewInternal(err)
	}

	// 4. Salvar no DB
	if err := ac.DB.Save(user).Error; err != nil {
		return exceptions.NewInternal(err)
	}

	return c.SendStatus(fiber.StatusNoContent) // 204
}

func (ac *Controller) GetMe(c *fiber.Ctx) error {

	user, ok := c.Locals("user").(*User)
	if !ok {
		return exceptions.NewInternal(errors.New("c.Locals(\"user\") não encontrado"))
	}

	response := MapUserToResponse(user)
	return c.JSON(response)
}
