package userapp

import "go-ca/internal/domain"

type AuthDTO struct {
	User         UserDTO `json:"user"`
	Token        string  `json:"token"`
	RefreshToken string  `json:"refresh_token"`
}

func NewAuthDTO(user *domain.User, token string, refresh_token string) AuthDTO {
	userDto := NewUserDTO(user)
	return AuthDTO{
		User:         userDto,
		Token:        token,
		RefreshToken: refresh_token,
	}
}
