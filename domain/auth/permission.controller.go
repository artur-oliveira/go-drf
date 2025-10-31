package auth

import (
	controllers "grf/core/controller"
	"grf/core/pagination"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type PermissionController struct {
	*controllers.GenericController[
		*Permission,
		*PermissionCreateDTO,
		*PermissionUpdateDTO,
		*PermissionPatchDTO,
		*PermissionResponseDTO,
		*PermissionFilterSet,
	]
}

func NewDefaultPermissionController(
	db *gorm.DB,
	validator *validator.Validate,
) *PermissionController {
	PermissionPaginator := pagination.NewLimitOffsetPagination[*Permission](
		10,
		100,
	)

	PermissionConfig := &controllers.ControllerConfig[
		*Permission,
		*PermissionCreateDTO,
		*PermissionUpdateDTO,
		*PermissionPatchDTO,
		*PermissionResponseDTO,
		*PermissionFilterSet,
	]{
		DB:        db,
		Validator: validator,
		Paginator: PermissionPaginator,

		NewFilterSet: func() *PermissionFilterSet {
			return new(PermissionFilterSet)
		},
		NewPatchDTO: func() *PermissionPatchDTO {
			return new(PermissionPatchDTO)
		},

		MapToResponse:    MapPermissionToResponse,
		MapCreateToModel: MapCreateToPermission,
		MapUpdateToModel: MapUpdateToPermission,
	}
	return NewPermissionController(PermissionConfig)
}

func NewPermissionController(
	config *controllers.ControllerConfig[
		*Permission,
		*PermissionCreateDTO,
		*PermissionUpdateDTO,
		*PermissionPatchDTO,
		*PermissionResponseDTO,
		*PermissionFilterSet,
	],
) *PermissionController {
	return &PermissionController{
		GenericController: controllers.NewGenericController(config),
	}
}
