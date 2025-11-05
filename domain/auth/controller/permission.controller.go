package controller

import (
	controllers "grf/core/controller"
	"grf/core/pagination"
	"grf/core/repository"
	"grf/core/service"
	"grf/domain/auth/dto"
	"grf/domain/auth/filter"
	"grf/domain/auth/mapper"
	"grf/domain/auth/model"
	"strconv"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func NewDefaultPermissionController(
	db *gorm.DB,
	validate *validator.Validate,
) *controllers.GenericController[
	*model.Permission, *dto.PermissionCreateDTO, *dto.PermissionUpdateDTO, *dto.PermissionPatchDTO, *dto.PermissionResponseDTO, *filter.PermissionFilterSet, uint64,
] {

	repo := repository.NewGenericRepository[*model.Permission, uint64](
		&repository.Config[*model.Permission, uint64]{
			DB:       db,
			NewModel: func() *model.Permission { return new(model.Permission) },
		},
	)

	svc := service.NewGenericService(
		&service.Config[*model.Permission, *dto.PermissionCreateDTO, *dto.PermissionUpdateDTO, *dto.PermissionPatchDTO, *dto.PermissionResponseDTO, *filter.PermissionFilterSet, uint64]{
			Repo:             repo,
			MapCreateToModel: mapper.MapCreateToPermission,
			MapUpdateToModel: mapper.MapUpdateToPermission,
		},
	)

	paginator := pagination.NewLimitOffsetPagination[*model.Permission](10, 100)
	config := &controllers.Config[
		*model.Permission, *dto.PermissionCreateDTO, *dto.PermissionUpdateDTO, *dto.PermissionPatchDTO, *dto.PermissionResponseDTO, *filter.PermissionFilterSet, uint64,
	]{
		Service:       svc,
		Validator:     validate,
		Paginator:     paginator,
		MapToResponse: mapper.MapPermissionToResponse,
		NewFilterSet:  func() *filter.PermissionFilterSet { return new(filter.PermissionFilterSet) },
		NewPatchDTO:   func() *dto.PermissionPatchDTO { return new(dto.PermissionPatchDTO) },
		ParseID: func(s string) (uint64, error) {
			id, err := strconv.ParseUint(s, 10, 64)
			return id, err
		},
	}

	return controllers.NewGenericController(config)
}
