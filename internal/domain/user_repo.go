package domain

import "context"

type UserRepository interface {
	BaseRepository[User, int32]
	GetByEmail(ctx context.Context, email string) (*User, error)
}
