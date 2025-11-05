package mapper

import (
	"grf/domain/auth/dto"
	"grf/domain/auth/model"
)

func MapPermissionToResponse(perm *model.Permission) *dto.PermissionResponseDTO {
	return &dto.PermissionResponseDTO{
		ID:          perm.ID,
		Module:      perm.Module,
		Action:      perm.Action,
		Description: perm.Description,
	}
}

func MapCreateToPermission(dto *dto.PermissionCreateDTO) *model.Permission {
	return &model.Permission{
		Module:      dto.Module,
		Action:      dto.Action,
		Description: dto.Description,
	}
}

func MapUpdateToPermission(dto *dto.PermissionUpdateDTO, perm *model.Permission) *model.Permission {
	perm.Description = dto.Description
	return perm
}
