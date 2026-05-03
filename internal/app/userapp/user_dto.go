package userapp

import (
	"go-ca/internal/domain"
	"time"
)

type UserDTO struct {
	ID        uint                `json:"id"`
	FirstName string              `json:"first_name"`
	LastName  string              `json:"last_name"`
	Email     string              `json:"email"`
	Role      domain.UserRoleType `json:"role"`
	IsActive  bool                `json:"is_active"`
	CreatedAt string              `json:"created_at"`
	UpdatedAt string              `json:"updated_at"`
}

func NewUserDTO(u *domain.User) UserDTO {
	return UserDTO{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}
