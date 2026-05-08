package userapp

import (
	"context"
	"go-ca/internal/app"
	"go-ca/internal/domain"
)

type GetUsersListQuery struct {
	domain.PaginationParams
}

type GetUsersListQueryHandler struct {
	userRepo domain.BaseRepository[*domain.User, uint]
}

func NewGetUsersListQueryHandler(
	userRepo domain.BaseRepository[*domain.User, uint],
) GetUsersListQueryHandler {
	return GetUsersListQueryHandler{
		userRepo: userRepo,
	}
}

func (g *GetUsersListQueryHandler) Handle(ctx context.Context, request GetUsersListQuery) (*domain.PaginatedResult[UserDTO], error) {
	userList, err := g.userRepo.GetPaged(ctx, request.PaginationParams)
	if err != nil {
		return nil, err
	}

	return app.MapPaginatedResult(&userList, NewUserDTO), nil
}
