package auth

func MapPermissionToResponse(perm *Permission) *PermissionResponseDTO {
	return &PermissionResponseDTO{
		ID:          perm.ID,
		Module:      perm.Module,
		Action:      perm.Action,
		Description: perm.Description,
	}
}

func MapCreateToPermission(dto *PermissionCreateDTO) *Permission {
	return &Permission{
		Module:      dto.Module,
		Action:      dto.Action,
		Description: dto.Description,
	}
}

func MapUpdateToPermission(dto *PermissionUpdateDTO, perm *Permission) *Permission {
	perm.Description = dto.Description
	return perm
}
