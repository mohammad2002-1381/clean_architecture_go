package queries

import (
	"clean_architecture_go/internal/app/features/users/dto"
	"clean_architecture_go/internal/app/mappers"
	currentuserservice "clean_architecture_go/internal/app/services/current_user_service"
	"clean_architecture_go/internal/domain"
	"clean_architecture_go/internal/pkg/error"
	"context"
	"fmt"
	"net/http"
)

type GetUserQuery struct{}

type GetUserQueryHandler struct {
	Mapper   mappers.UserMapper
	UserRepo domain.UserRepository
}

func NewGetUserQueryHandler(mapper mappers.UserMapper, userRepo domain.UserRepository) GetUserQueryHandler {
	return GetUserQueryHandler{
		Mapper:   mapper,
		UserRepo: userRepo,
	}
}

func (h *GetUserQueryHandler) Handle(ctx context.Context, req GetUserQuery) (*dto.UserDTO, *error.RequestError) {
	currentUser := currentuserservice.NewCurrentUserService(ctx)
	userID, err := currentUser.GetUserID()

	if err != nil {
		return nil, error.NewRequestError(http.StatusUnauthorized, err.Error())
	}

	user, err := h.UserRepo.FindAsync(ctx, userID)
	if err != nil {
		return nil, error.NewRequestError(
			http.StatusInternalServerError,
			fmt.Sprintf("failed to find user: %v", err),
		)
	}

	if user == nil {
		return nil, error.NewRequestError(http.StatusNotFound, "user not found")
	}

	dto := h.Mapper.MapUser(user)
	return &dto, nil
}