// internal/app/services/current_user_service/current_user.go
package currentuserservice

import (
	"clean_architecture_go/internal/app/services"
	"context"
	"errors"
)

func getUserID(ctx context.Context) (int32, error) {
	userID, ok := ctx.Value("user_id").(int32)
	if !ok {
		return 0, errors.New("user_id not found in context")
	}
	return userID, nil
}

func getUserRole(ctx context.Context) (string, error) {
	role, ok := ctx.Value("role").(string)
	if !ok {
		return "", errors.New("role not found in context")
	}
	return role, nil
}

type currentUserService struct {
	ctx context.Context
}

func NewCurrentUserService(ctx context.Context) services.CurrentUserService {
	return &currentUserService{
		ctx: ctx,
	}
}

func (s *currentUserService) GetUserID() (int32, error) {
	return getUserID(s.ctx)
}

func (s *currentUserService) GetUserRole() (string, error) {
	return getUserRole(s.ctx)
}

func (s *currentUserService) IsAdmin() bool {
	role, err := getUserRole(s.ctx)
	if err != nil {
		return false
	}
	return role == "admin"
}
