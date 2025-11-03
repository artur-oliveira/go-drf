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

func NewDefaultGroupController(
	db *gorm.DB,
	validate *validator.Validate,
) *controllers.GenericController[
	*Group, *GroupCreateDTO, *GroupUpdateDTO, *GroupPatchDTO, *GroupResponseDTO, *GroupFilterSet, IDType,
] {

	groupRepo := repository.NewGenericRepository[*Group, IDType](
		&repository.Config[*Group, IDType]{
			DB:       db,
			NewModel: func() *Group { return new(Group) },
		},
	)

	groupService := NewGroupService(
		&service.ServiceConfig[*Group, *GroupCreateDTO, *GroupUpdateDTO, *GroupPatchDTO, *GroupResponseDTO, *GroupFilterSet, IDType]{
			Repo:             groupRepo,
			MapCreateToModel: MapCreateToGroup,
			MapUpdateToModel: MapUpdateToGroup,
		},
		db,
	)

	groupPaginator := pagination.NewLimitOffsetPagination[*Group](10, 100)
	groupConfig := &controllers.ControllerConfig[
		*Group, *GroupCreateDTO, *GroupUpdateDTO, *GroupPatchDTO, *GroupResponseDTO, *GroupFilterSet, IDType,
	]{
		Service:       groupService,
		Validator:     validate,
		Paginator:     groupPaginator,
		MapToResponse: MapGroupToResponse,
		NewFilterSet:  func() *GroupFilterSet { return new(GroupFilterSet) },
		NewPatchDTO:   func() *GroupPatchDTO { return new(GroupPatchDTO) },
		ParseID: func(s string) (IDType, error) {
			id, err := strconv.ParseUint(s, 10, 64)
			return id, err
		},
	}

	return controllers.NewGenericController(groupConfig)
}
