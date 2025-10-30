package auth

import "golang.org/x/crypto/bcrypt"

func MapUserToResponse(user User) UserResponseDTO {
	return UserResponseDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		IsAdmin:   user.IsAdmin,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}
}

func MapCreateToUser(dto UserCreateDTO) User {
	// (Aqui é o local para lógica de negócios, como HASH de senha)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(dto.Password), 12)

	return User{
		Username: dto.Username,
		Email:    dto.Email,
		Password: string(hashedPassword), // Em produção: string(hashedPassword)
	}
}

func MapUpdateToUser(dto UserUpdateDTO, user User) User {
	user.Username = dto.Username
	user.Email = dto.Email
	user.IsAdmin = *dto.IsAdmin

	return user
}
