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

func (ac *Controller) Login(c *fiber.Ctx) error {
	var input LoginDTO
	if err := c.BodyParser(&input); err != nil {
		return exceptions.NewBadRequest("Payload inválido", err)
	}
	if err := ac.Validator.Struct(input); err != nil {
		return err
	}

	var user User
	if err := ac.DB.Where("username = ? OR email = ?", input.Username, input.Username).First(&user).Error; err != nil {
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

func (ac *Controller) GetMe(c *fiber.Ctx) error {

	user, ok := c.Locals("user").(*User)
	if !ok {
		return exceptions.NewInternal(errors.New("c.Locals(\"user\") não encontrado"))
	}

	response := MapUserToResponse(user)
	return c.JSON(response)
}
