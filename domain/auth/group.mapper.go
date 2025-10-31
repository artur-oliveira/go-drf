package auth

func MapGroupToResponse(group *Group) *GroupResponseDTO {
	resp := GroupResponseDTO{
		ID:   group.ID,
		Name: group.Name,
	}

	if group.Permissions != nil && len(group.Permissions) > 0 {
		resp.Permissions = make([]PermissionResponseDTO, len(group.Permissions))
		for i, perm := range group.Permissions {
			resp.Permissions[i] = *MapPermissionToResponse(perm)
		}
	} else {
		resp.Permissions = make([]PermissionResponseDTO, 0)
	}

	return &resp
}

func MapCreateToGroup(dto *GroupCreateDTO) *Group {
	return &Group{
		Name: dto.Name,
	}
}

func MapUpdateToGroup(dto *GroupUpdateDTO, group *Group) *Group {
	group.Name = dto.Name
	return group
}
