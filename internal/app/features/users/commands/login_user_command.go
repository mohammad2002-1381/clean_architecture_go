// commands/login_user_command.go
package commands

import (
	"clean_architecture_go/internal/app"
	"clean_architecture_go/internal/app/features/users/dto"
	"clean_architecture_go/internal/app/mappers"
	"clean_architecture_go/internal/app/services"
	"clean_architecture_go/internal/domain"
	errs "clean_architecture_go/internal/pkg/error"
	"context"
	"fmt"
	"net/http"
)

type LoginUserCommand struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserCommandHandler struct {
	Mapper               mappers.UserMapper
	UserRepo             domain.UserRepository
	TokenRepo            domain.BaseRepository[domain.Token, int32]
	TokenService         services.TokenService
	PasswordService      services.PasswordService
	BackgroundJobService services.BackgroundJobService
}

func NewLoginUserCommandHandler(mapper mappers.UserMapper,
	userRepo domain.UserRepository,
	tokenRepo domain.BaseRepository[domain.Token, int32]) LoginUserCommandHandler {
	return LoginUserCommandHandler{
		Mapper:               mapper,
		UserRepo:             userRepo,
		TokenRepo:            tokenRepo,
		TokenService:         app.GetTokenService(),
		PasswordService:      app.GetPwdService(),
		BackgroundJobService: app.GetBKService(),
	}
}

func (h *LoginUserCommandHandler) Handle(ctx context.Context, req LoginUserCommand) (*dto.AuthDto, *errs.RequestError) {
	if req.Email == "" || req.Password == "" {
		return nil, errs.NewRequestError(http.StatusBadRequest, "email and password are required")
	}

	user, err := h.UserRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errs.NewRequestError(
			http.StatusInternalServerError,
			fmt.Sprintf("failed to find user: %v", err),
		)
	}

	if user == nil {
		return nil, errs.NewRequestError(http.StatusUnauthorized, "invalid email or password")
	}

	if !h.PasswordService.Verify(req.Password, user.PasswordHash) {
		return nil, errs.NewRequestError(http.StatusUnauthorized, "invalid email or password")
	}

	if !user.IsActive {
		return nil, errs.NewRequestError(http.StatusForbidden, "account is deactivated")
	}

	token, err := h.TokenService.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return nil, errs.NewRequestError(
			http.StatusInternalServerError,
			fmt.Sprintf("failed to generate token: %v", err),
		)
	}

	refreshToken, err := h.TokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errs.NewRequestError(
			http.StatusInternalServerError,
			fmt.Sprintf("failed to generate refresh token: %v", err),
		)
	}

	h.BackgroundJobService.Enqueue(h.InsertToken(ctx, token, refreshToken, user.ID))

	userDto := h.Mapper.MapUser(user)
	return &dto.AuthDto{
		User:         userDto,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (h *LoginUserCommandHandler) InsertToken(ctx context.Context, token string, refreshToken string, userID int32) services.Job {
	return func(ctx context.Context) error {
		t := domain.NewToken(token, refreshToken, userID)
		h.TokenRepo.Add(ctx, t)
		_, err := h.TokenRepo.UnitOfWork().SaveEntitiesAsync(ctx)
		if err != nil {
			return err
		}
		return nil
	}
}
