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
	services "grf/domain/auth/service"
	"strconv"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func NewDefaultGroupController(
	db *gorm.DB,
	validate *validator.Validate,
) *controllers.GenericController[
	*model.Group, *dto.GroupCreateDTO, *dto.GroupUpdateDTO, *dto.GroupPatchDTO, *dto.GroupResponseDTO, *filter.GroupFilterSet, uint64,
] {

	groupRepo := repository.NewGenericRepository[*model.Group, uint64](
		&repository.Config[*model.Group, uint64]{
			DB:       db,
			NewModel: func() *model.Group { return new(model.Group) },
		},
	)

	groupService := services.NewGroupService(
		&service.Config[*model.Group, *dto.GroupCreateDTO, *dto.GroupUpdateDTO, *dto.GroupPatchDTO, *dto.GroupResponseDTO, *filter.GroupFilterSet, uint64]{
			Repo:             groupRepo,
			MapCreateToModel: mapper.MapCreateToGroup,
			MapUpdateToModel: mapper.MapUpdateToGroup,
		},
		db,
	)

	groupPaginator := pagination.NewLimitOffsetPagination[*model.Group](10, 100)
	groupConfig := &controllers.Config[
		*model.Group, *dto.GroupCreateDTO, *dto.GroupUpdateDTO, *dto.GroupPatchDTO, *dto.GroupResponseDTO, *filter.GroupFilterSet, uint64,
	]{
		Service:       groupService,
		Validator:     validate,
		Paginator:     groupPaginator,
		MapToResponse: mapper.MapGroupToResponse,
		NewFilterSet:  func() *filter.GroupFilterSet { return new(filter.GroupFilterSet) },
		NewPatchDTO:   func() *dto.GroupPatchDTO { return new(dto.GroupPatchDTO) },
		ParseID: func(s string) (uint64, error) {
			id, err := strconv.ParseUint(s, 10, 64)
			return id, err
		},
	}

	return controllers.NewGenericController(groupConfig)
}
