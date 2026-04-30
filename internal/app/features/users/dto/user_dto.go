package dto

import (
	"clean_architecture_go/internal/domain"
)

type UserDTO struct {
	ID        int32               `json:"id"`
	FirstName string              `json:"first_name"`
	LastName  string              `json:"last_name"`
	Email     string              `json:"email"`
	Role      domain.UserRoleType `json:"role"`
	IsActive  bool                `json:"is_active"`
	CreatedAt string              `json:"created_at"`
	UpdatedAt string              `json:"updated_at"`
}

func NewUserDto(
	id int32,
	firstName string,
	lastName string,
	email string,
	role domain.UserRoleType,
	isActive bool,
	createdAt string,
	updatedAt string) UserDTO {

	return UserDTO{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Role:      role,
		IsActive:  isActive,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
