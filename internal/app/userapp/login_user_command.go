package userapp

import (
	"context"
	"errors"

	"go-ca/internal/app/service"
	"go-ca/internal/domain"
)

type LoginUserCommand struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserCommandHandler struct {
	userRepo        domain.BaseRepository[domain.User, uint]
	passwordService service.PasswordService
	jwtService      service.JWTService
}

func NewLoginUserCommandHandler(
	userRepo domain.BaseRepository[domain.User, uint],
	passwordService service.PasswordService,
	jwtService service.JWTService,
) LoginUserCommandHandler {
	return LoginUserCommandHandler{
		userRepo:        userRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

func (c *LoginUserCommandHandler) Handle(ctx context.Context, request LoginUserCommand) (*AuthDTO, error) {
	users, err := c.userRepo.Where(ctx, &domain.User{Email: request.Email}).Get(ctx)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	user := users[0]

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	isValid := c.passwordService.Verify(request.Password, user.PasswordHash)
	if !isValid {
		return nil, errors.New("invalid email or password")
	}

	token, err := c.jwtService.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := c.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	dto := NewAuthDTO(user, token, refreshToken)

	return &dto, nil
}
