package userapp

import (
	"context"
	"go-ca/internal/app/service"
	"go-ca/internal/domain"
)

type GetUserQuery struct{}

type GetUserQueryHandler struct {
	userRepo           domain.BaseRepository[*domain.User, uint]
	currentUserService service.CurrentUserService
}

func NewGetUserQueryHandler(
	userRepo	domain.BaseRepository[*domain.User, uint],
	currentUserService service.CurrentUserService,
) GetUserQueryHandler {
	return GetUserQueryHandler{
		userRepo: userRepo,
		currentUserService: currentUserService,
	}
}

func (h *GetUserQueryHandler) Handle(ctx context.Context, request GetUserQuery) (*UserDTO, error) {
	userID, err := h.currentUserService.GetUserID(ctx)
	
	if err != nil {
		return nil, err
	}

	user, err := h.userRepo.Find(ctx, userID)

	if err != nil {
		return nil, err
	}

	dto := NewUserDTO(user)

	return &dto, nil
}
