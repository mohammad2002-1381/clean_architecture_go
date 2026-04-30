package repositories

import (
	"clean_architecture_go/internal/domain"
	"clean_architecture_go/internal/infra"
	"context"
	"fmt"
)

type UserRepository struct {
	infra.BaseRepository[domain.User, int32]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		BaseRepository: *infra.NewBaseRepository[domain.User, int32](),
	}
}

func (u *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	users, err := u.Get(ctx, "email = ?", email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return users[0], nil
}
