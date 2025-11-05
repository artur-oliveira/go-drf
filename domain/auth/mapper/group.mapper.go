package mapper

import (
	"grf/domain/auth/dto"
	"grf/domain/auth/model"
)

func MapGroupToResponse(group *model.Group) *dto.GroupResponseDTO {
	resp := dto.GroupResponseDTO{
		ID:   group.ID,
		Name: group.Name,
	}

	if group.Permissions != nil && len(group.Permissions) > 0 {
		resp.Permissions = make([]dto.PermissionResponseDTO, len(group.Permissions))
		for i, perm := range group.Permissions {
			resp.Permissions[i] = *MapPermissionToResponse(perm)
		}
	} else {
		resp.Permissions = make([]dto.PermissionResponseDTO, 0)
	}

	return &resp
}

func MapCreateToGroup(dto *dto.GroupCreateDTO) *model.Group {
	return &model.Group{
		Name: dto.Name,
	}
}

func MapUpdateToGroup(dto *dto.GroupUpdateDTO, group *model.Group) *model.Group {
	group.Name = dto.Name
	return group
}
