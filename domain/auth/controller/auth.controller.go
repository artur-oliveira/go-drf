package controller

import (
	"errors"
	"grf/core/config"
	"grf/core/exceptions"
	"grf/domain/auth/dto"
	"grf/domain/auth/mapper"
	"grf/domain/auth/model"
	"grf/domain/auth/repository"
	"grf/domain/auth/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Controller struct {
	UserRepo     *repository.UserRepository
	Validator    *validator.Validate
	TokenService *service.TokenService
}

func NewAuthController(
	db *gorm.DB,
	config *config.Config,
	validate *validator.Validate,
) *Controller {
	return &Controller{
		UserRepo:     repository.NewUserRepository(db),
		Validator:    validate,
		TokenService: service.NewTokenService(db, config),
	}
}

func (ac *Controller) ObtainToken(c *fiber.Ctx) error {
	var input dto.ObtainTokenDTO
	if err := c.BodyParser(&input); err != nil {
		return exceptions.NewBadRequest("invalid_payload", err)
	}
	if err := ac.Validator.Struct(input); err != nil {
		return err
	}

	user, err := ac.UserRepo.FindUserByEmailOrUsername(input.Login)
	if err != nil {
		return exceptions.NewUnauthorized("invalid_credentials", err)
	}

	if !user.IsActive {
		return exceptions.NewUnauthorized("inactive_user", nil)
	}

	if !user.CheckPassword(input.Password) {
		return exceptions.NewUnauthorized("invalid_credentials", nil)
	}

	access, refresh, err := ac.TokenService.GenerateTokenPair(&user)
	if err != nil {
		return exceptions.NewInternal(err)
	}

	return c.JSON(dto.TokenResponseDTO{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

func (ac *Controller) ObtainTokenRefresh(c *fiber.Ctx) error {
	var input dto.RefreshTokenDTO
	if err := c.BodyParser(&input); err != nil {
		return exceptions.NewBadRequest("invalid_payload", err)
	}
	if err := ac.Validator.Struct(input); err != nil {
		return err
	}

	user, err := ac.TokenService.ValidateToken(input.Refresh, "refresh")
	if err != nil {
		return exceptions.NewError(fiber.StatusUnauthorized, err.Error(), err)
	}

	access, refresh, err := ac.TokenService.GenerateTokenPair(user)
	if err != nil {
		return exceptions.NewInternal(err)
	}

	return c.JSON(dto.TokenResponseDTO{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

func (ac *Controller) ChangePassword(c *fiber.Ctx) error {
	var input dto.ChangePasswordDTO
	if err := c.BodyParser(&input); err != nil {
		return exceptions.NewBadRequest("invalid_payload", err)
	}
	if err := ac.Validator.Struct(input); err != nil {
		return err
	}

	user, ok := c.Locals("user").(*model.User)
	if !ok {
		return exceptions.NewInternal(errors.New("c.Locals(\"user\") não encontrado"))
	}

	if !user.CheckPassword(input.OldPassword) {
		return exceptions.NewBadRequest("incorrect_old_password", nil)
	}

	if input.NewPassword != input.NewPassword {
		return exceptions.NewBadRequest("incorrect_new_password", nil)
	}

	if err := user.SetPassword(input.NewPassword); err != nil {
		return exceptions.NewInternal(err)
	}

	if err := ac.UserRepo.Update(user); err != nil {
		return exceptions.NewInternal(err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (ac *Controller) GetMe(c *fiber.Ctx) error {

	user, ok := c.Locals("user").(*model.User)
	if !ok {
		return exceptions.NewInternal(errors.New("c.Locals(\"user\") não encontrado"))
	}

	response := mapper.MapUserToResponse(user)
	return c.JSON(response)
}
