package auth

import (
	controllers "grf/core/controller"
	"grf/core/pagination"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type UserController struct {
	*controllers.GenericController[
		*User,
		*UserCreateDTO,
		*UserUpdateDTO,
		*UserPatchDTO,
		*UserResponseDTO,
		*UserFilterSet,
	]
}

func NewDefaultUserController(
	db *gorm.DB,
	validator *validator.Validate,
) *UserController {
	userPaginator := pagination.NewLimitOffsetPagination[*User](
		10,
		100,
	)

	userConfig := &controllers.ControllerConfig[
		*User,
		*UserCreateDTO,
		*UserUpdateDTO,
		*UserPatchDTO,
		*UserResponseDTO,
		*UserFilterSet,
	]{
		DB:        db,
		Validator: validator,
		Paginator: userPaginator,

		NewFilterSet: func() *UserFilterSet {
			return new(UserFilterSet)
		},
		NewPatchDTO: func() *UserPatchDTO {
			return new(UserPatchDTO)
		},

		MapToResponse:    MapUserToResponse,
		MapCreateToModel: MapCreateToUser,
		MapUpdateToModel: MapUpdateToUser,
	}
	return NewUserController(userConfig)
}

func NewUserController(
	config *controllers.ControllerConfig[
		*User,
		*UserCreateDTO,
		*UserUpdateDTO,
		*UserPatchDTO,
		*UserResponseDTO,
		*UserFilterSet,
	],
) *UserController {
	return &UserController{
		GenericController: controllers.NewGenericController(config),
	}
}
