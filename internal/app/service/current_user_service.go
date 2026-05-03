package service

import (
	"context"
	"errors"
	"go-ca/internal/domain"
)

type currentUserService struct{}

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

func NewCurrentUserService() *currentUserService {
	return &currentUserService{}
}

func (s *currentUserService) GetUserID(ctx context.Context) (uint, error) {
	val := ctx.Value(UserIDKey)
	id, ok := val.(uint)
	if !ok {
		return 0, errors.New("unauthorized")
	}
	return id, nil
}

func (s *currentUserService) GetUserRole(ctx context.Context) (domain.UserRoleType, error) {
	val := ctx.Value(RoleKey)
	role, ok := val.(domain.UserRoleType)
	if !ok {
		return "", errors.New("role not found or invalid type in context")
	}
	return role, nil
}
