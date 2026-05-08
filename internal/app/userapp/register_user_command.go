package userapp

import (
	"context"

	"go-ca/internal/app/service"
	"go-ca/internal/domain"
)

type RegisterUserCommand struct {
	FirstName string              `json:"first_name"`
	LastName  string              `json:"last_name"`
	Email     string              `json:"email"`
	Password  string              `json:"password"`
	Role      domain.UserRoleType `json:"role"`
}

type RegisterUserCommandHandler struct {
	userRepo        domain.BaseRepository[*domain.User, uint]
	passwordService service.PasswordService
	jwtService      service.JWTService
}

func NewRegisterUserCommandHandler(
	userRepo domain.BaseRepository[*domain.User, uint],
	passwordService service.PasswordService,
	jwtService service.JWTService,
) RegisterUserCommandHandler {
	return RegisterUserCommandHandler{
		userRepo:        userRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

func (c *RegisterUserCommandHandler) Handle(ctx context.Context, request RegisterUserCommand) (*AuthDTO, error) {
	hashedPassword, err := c.passwordService.Hash(request.Password)
	if err != nil {
		return nil, err
	}

	u := domain.NewUser(
		request.FirstName,
		request.LastName,
		request.Email,
		hashedPassword,
		request.Role,
		true,
	)

	err = c.userRepo.Add(ctx, u)
	if err != nil {
		return nil, err
	}

	token, err := c.jwtService.GenerateToken(u.ID, string(request.Role))

	if err != nil {
		return nil, err
	}

	refreshToken, err := c.jwtService.GenerateRefreshToken(u.ID)

	if err != nil {
		return nil, err
	}

	dto := NewAuthDTO(u, token, refreshToken)

	return &dto, nil
}
