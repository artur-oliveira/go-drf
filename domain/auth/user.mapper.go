package auth

import (
	"fmt"
)

func MapUserToResponse(user *User) *UserResponseDTO {
	return &UserResponseDTO{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		IsActive:    user.IsActive,
		IsStaff:     user.IsStaff,
		IsSuperuser: user.IsSuperuser,
		LastLogin:   user.LastLogin,
		CreatedAt:   user.CreatedAt,
	}
}

func MapCreateToUser(dto *UserCreateDTO) *User {
	user := User{
		Username:  dto.Username,
		Email:     dto.Email,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		IsActive:  true,
		IsStaff:   false,
	}

	if err := user.SetPassword(dto.Password); err != nil {
		panic(fmt.Sprintf("Falha crítica ao gerar hash de senha: %v", err))
	}

	if dto.IsActive != nil {
		user.IsActive = *dto.IsActive
	}
	if dto.IsStaff != nil {
		user.IsStaff = *dto.IsStaff
	}

	return &user
}

func MapUpdateToUser(dto *UserUpdateDTO, user *User) *User {

	user.Username = dto.Username
	user.Email = dto.Email
	user.FirstName = dto.FirstName
	user.LastName = dto.LastName
	user.IsActive = dto.IsActive
	user.IsStaff = dto.IsStaff
	user.IsSuperuser = dto.IsSuperuser

	return user
}
