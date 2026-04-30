package mappers

import (
	"clean_architecture_go/internal/app/features/users/dto"
	"clean_architecture_go/internal/domain"
)

type UserMapper struct{}

func (*UserMapper) MapUser(user *domain.User) dto.UserDTO {
	return dto.NewUserDto(
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Role,
		user.IsActive,
		user.CreatedAt.GoString(),
		user.UpdatedAt.GoString(),
	)
}

func (u *UserMapper) MapAuth(user *domain.User, token string, refreshToken string) dto.AuthDto {
	userDto := u.MapUser(user)
	return dto.NewAuthDto(
		userDto,
		token,
		refreshToken,
	)
}

func NewUserMapper() UserMapper {
	return UserMapper{}
}
