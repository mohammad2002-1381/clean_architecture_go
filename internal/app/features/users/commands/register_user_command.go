// commands/register_user_command.go
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

type RegisterUserCommand struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type RegisterUserCommandHandler struct {
	Mapper            mappers.UserMapper
	UserRepo          domain.UserRepository
	TokenRepo         domain.BaseRepository[domain.Token, int32]
	TokenService      services.TokenService
	PasswordService   services.PasswordService
	BackgroundService services.BackgroundJobService
}

func NewRegisterUserCommandHandler(mapper mappers.UserMapper,
	userRepo domain.UserRepository,
	tokenRepo domain.BaseRepository[domain.Token, int32]) RegisterUserCommandHandler {
	return RegisterUserCommandHandler{
		Mapper:            mapper,
		UserRepo:          userRepo,
		TokenRepo:         tokenRepo,
		TokenService:      app.GetTokenService(),
		PasswordService:   app.GetPwdService(),
		BackgroundService: app.GetBKService(),
	}
}

func (h *RegisterUserCommandHandler) Handle(ctx context.Context, req RegisterUserCommand) (*dto.AuthDto, *errs.RequestError) {
	if req.Email == "" || req.Password == "" {
		return nil, errs.NewRequestError(http.StatusBadRequest, "email and password are required")
	}

	// Check if user already exists
	existing, _ := h.UserRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errs.NewRequestError(http.StatusConflict, "email already registered")
	}

	// Hash password
	hashedPassword, err := h.PasswordService.Hash(req.Password)
	if err != nil {
		return nil, errs.NewRequestError(
			http.StatusInternalServerError,
			"failed to hash password",
		)
	}

	// Create user
	user := domain.NewUser(
		req.FirstName,
		req.LastName,
		req.Email,
		domain.UserRoleType("user"),
		hashedPassword,
		true,
	)

	uow := h.UserRepo.UnitOfWork()
	h.UserRepo.Add(ctx, user)
	_, err = uow.SaveEntitiesAsync(ctx)
	if err != nil {
		return nil, errs.NewRequestError(
			http.StatusInternalServerError,
			fmt.Sprintf("failed to save user: %v", err),
		)
	}

	// Generate tokens
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

	h.BackgroundService.Enqueue(h.InsertToken(ctx, token, refreshToken, user.ID))

	userDto := h.Mapper.MapUser(user)
	return &dto.AuthDto{
		User:         userDto,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (h *RegisterUserCommandHandler) InsertToken(ctx context.Context, token string, refreshToken string, userID int32) services.Job {
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
