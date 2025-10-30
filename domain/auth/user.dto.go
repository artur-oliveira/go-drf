package auth

import "time"

type UserCreateDTO struct {
	Username string `json:"username" validate:"required,gte=2,lte=50"`
	Email    string `json:"email" validate:"required,gte=3,lte=255,email"`
	Password string `json:"password" validate:"required,gte=8,lte=128"`
}

type UserUpdateDTO struct {
	Username string `json:"username" validate:"required,gte=2,lte=50"`
	Email    string `json:"email" validate:"required,gte=3,lte=255,email"`
	Password string `json:"password" validate:"required,gte=8,lte=128"`
	IsAdmin  bool   `json:"is_admin" validate:"required"`
}

type UserResponseDTO struct {
	ID        uint64    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"ative"`
	IsAdmin   bool      `json:"admin"`
}

type UserPatchDTO struct {
}

func (d *UserPatchDTO) ToPatchMap() map[string]interface{} {
	return nil
}
