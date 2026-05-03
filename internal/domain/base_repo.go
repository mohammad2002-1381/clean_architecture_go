package domain

import (
	"context"
)

type BaseRepository[T any, TID comparable] interface {
	Add(context.Context, *T) error
	Find(context.Context, TID) (*T, error)
	First(context.Context) (*T, error)
	Get(context.Context) ([]*T, error)
	Update(context.Context, *T) error
	Delete(context.Context, TID) error
	Where(context.Context, *T) BaseRepository[T, TID]
	Or(context.Context, *T) BaseRepository[T, TID]
	GetPaged(ctx context.Context, params PaginationParams) (PaginatedResult[T], error)
}
