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

func NewDefaultPermissionController(
	db *gorm.DB,
	validate *validator.Validate,
) *controllers.GenericController[
	*Permission, *PermissionCreateDTO, *PermissionUpdateDTO, *PermissionPatchDTO, *PermissionResponseDTO, *PermissionFilterSet, IDType,
] {

	repo := repository.NewGenericRepository[*Permission, IDType](
		&repository.Config[*Permission, IDType]{
			DB:       db,
			NewModel: func() *Permission { return new(Permission) },
		},
	)

	svc := service.NewGenericService(
		&service.ServiceConfig[*Permission, *PermissionCreateDTO, *PermissionUpdateDTO, *PermissionPatchDTO, *PermissionResponseDTO, *PermissionFilterSet, IDType]{
			Repo:             repo,
			MapCreateToModel: MapCreateToPermission,
			MapUpdateToModel: MapUpdateToPermission,
		},
	)

	paginator := pagination.NewLimitOffsetPagination[*Permission](10, 100)
	config := &controllers.ControllerConfig[
		*Permission, *PermissionCreateDTO, *PermissionUpdateDTO, *PermissionPatchDTO, *PermissionResponseDTO, *PermissionFilterSet, IDType,
	]{
		Service:       svc,
		Validator:     validate,
		Paginator:     paginator,
		MapToResponse: MapPermissionToResponse,
		NewFilterSet:  func() *PermissionFilterSet { return new(PermissionFilterSet) },
		NewPatchDTO:   func() *PermissionPatchDTO { return new(PermissionPatchDTO) },
		ParseID: func(s string) (IDType, error) {
			id, err := strconv.ParseUint(s, 10, 64)
			return id, err
		},
	}

	return controllers.NewGenericController(config)
}
