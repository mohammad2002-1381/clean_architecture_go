package userapp

import (
	"go-ca/internal/domain"
	"time"
)

type UserDTO struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewUserDTO(user *domain.User) UserDTO {
	if user == nil {
		return UserDTO{}
	}

	return UserDTO{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}

func NewUserDTOList(users []domain.User) []UserDTO {
	userDtos := make([]UserDTO, 0, len(users))
	for _, u := range users {
		userDtos = append(userDtos, NewUserDTO(&u))
	}
	return userDtos
}
