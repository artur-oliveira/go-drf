package controller

import (
	controllers "grf/core/controller"
	"grf/core/pagination"
	"grf/core/service"
	"grf/domain/auth/dto"
	"grf/domain/auth/filter"
	"grf/domain/auth/mapper"
	"grf/domain/auth/model"
	"grf/domain/auth/repository"
	"strconv"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func NewDefaultUserController(
	db *gorm.DB,
	validate *validator.Validate,
) *controllers.GenericController[
	*model.User, *dto.UserCreateDTO, *dto.UserUpdateDTO, *dto.UserPatchDTO, *dto.UserResponseDTO, *filter.UserFilterSet, uint64,
] {
	userRepo := repository.NewUserRepository(db)

	userService := service.NewGenericService(
		&service.Config[*model.User, *dto.UserCreateDTO, *dto.UserUpdateDTO, *dto.UserPatchDTO, *dto.UserResponseDTO, *filter.UserFilterSet, uint64]{
			Repo:             userRepo,
			MapCreateToModel: mapper.MapCreateToUser,
			MapUpdateToModel: mapper.MapUpdateToUser,
		},
	)

	userPaginator := pagination.NewLimitOffsetPagination[*model.User](10, 100)
	userConfig := &controllers.Config[
		*model.User, *dto.UserCreateDTO, *dto.UserUpdateDTO, *dto.UserPatchDTO, *dto.UserResponseDTO, *filter.UserFilterSet, uint64,
	]{
		Service:   userService,
		Validator: validate,
		Paginator: userPaginator,

		MapToResponse: mapper.MapUserToResponse,

		NewFilterSet: func() *filter.UserFilterSet { return new(filter.UserFilterSet) },
		NewPatchDTO:  func() *dto.UserPatchDTO { return new(dto.UserPatchDTO) },

		ParseID: func(s string) (uint64, error) {
			id, err := strconv.ParseUint(s, 10, 64)
			return id, err
		},
	}

	return controllers.NewGenericController(userConfig)
}
