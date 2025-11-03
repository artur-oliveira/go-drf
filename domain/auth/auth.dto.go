package auth

type ObtainTokenDTO struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type RefreshTokenDTO struct {
	Refresh string `json:"refresh" validate:"required"`
}
type ChangePasswordDTO struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}
