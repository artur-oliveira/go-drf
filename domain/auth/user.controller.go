package auth

import (
	controllers "grf/core/controller"
	"grf/core/pagination"
	"grf/core/repository"
	"grf/core/service"
	"strconv"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func NewDefaultUserController(
	db *gorm.DB,
	validate *validator.Validate,
) *controllers.GenericController[
	*User, *UserCreateDTO, *UserUpdateDTO, *UserPatchDTO, *UserResponseDTO, *UserFilterSet, IDType,
] {

	userRepo := repository.NewGenericRepository[*User, IDType](
		&repository.Config[*User, IDType]{
			DB:       db,
			NewModel: func() *User { return new(User) },
		},
	)

	userService := service.NewGenericService(
		&service.ServiceConfig[*User, *UserCreateDTO, *UserUpdateDTO, *UserPatchDTO, *UserResponseDTO, *UserFilterSet, IDType]{
			Repo:             userRepo,
			MapCreateToModel: MapCreateToUser,
			MapUpdateToModel: MapUpdateToUser,
		},
	)

	userPaginator := pagination.NewLimitOffsetPagination[*User](10, 100)
	userConfig := &controllers.ControllerConfig[
		*User, *UserCreateDTO, *UserUpdateDTO, *UserPatchDTO, *UserResponseDTO, *UserFilterSet, IDType,
	]{
		Service:   userService,
		Validator: validate,
		Paginator: userPaginator,

		MapToResponse: MapUserToResponse,

		NewFilterSet: func() *UserFilterSet { return new(UserFilterSet) },
		NewPatchDTO:  func() *UserPatchDTO { return new(UserPatchDTO) },

		ParseID: func(s string) (IDType, error) {
			id, err := strconv.ParseUint(s, 10, 64)
			return id, err
		},
	}

	return controllers.NewGenericController(userConfig)
}
