package userapp

import (
	"context"
	"go-ca/internal/app/service"
	"go-ca/internal/domain"
)

type RefreshTokenCommand struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenCommandHandler struct {
	tokenRepo          domain.BaseRepository[*domain.Token, uint]
	currentUserService service.CurrentUserService
	jwtService         service.JWTService
}

func NewRefreshTokenCommandHandler(
	tokenRepo domain.BaseRepository[*domain.Token, uint],
	currentUserService service.CurrentUserService,
	jwtService service.JWTService,
) RefreshTokenCommandHandler {
	return RefreshTokenCommandHandler{
		tokenRepo:          tokenRepo,
		currentUserService: currentUserService,
		jwtService:         jwtService,
	}
}

func (r *RefreshTokenCommandHandler) Handle(ctx context.Context, request RefreshTokenCommand) (*TokenDTO, error) {
	token, err := r.currentUserService.GetToken(ctx)

	if err != nil {
		return nil, err
	}

	role, err := r.currentUserService.GetUserRole(ctx)

	if err != nil {
		return nil, err
	}

	userID, err := r.currentUserService.GetUserID(ctx)

	if err != nil {
		return nil, err
	}

	tokenEn, err := r.tokenRepo.Where(ctx, &domain.Token{Token: token}).First(ctx)

	if err != nil {
		return nil, err
	}

	token, err = r.jwtService.GenerateToken(userID, string(role))

	if err != nil {
		return nil, err
	}

	refreshToken, err := r.jwtService.GenerateRefreshToken(userID)

	if err != nil {
		return nil, err
	}

	tokenEn.SetToken(token)
	tokenEn.SetRefreshToken(refreshToken)

	r.tokenRepo.Update(ctx, tokenEn)

	dto := NewTokenDTO(token, refreshToken)

	return &dto, nil
}
