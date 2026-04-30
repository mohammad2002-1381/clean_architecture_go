package domain

import (
	"context"
	"clean_architecture_go/internal/domain/extensions"
)

type BaseRepository[TEntity IBaseEntity[TID], TID comparable] interface {
	Add(ctx context.Context, entity *TEntity) (*TEntity, error)
	AddRange(ctx context.Context, entities []*TEntity) error
	Delete(ctx context.Context, entity *TEntity) error
	DeleteRange(ctx context.Context, entities []*TEntity) error
	FindAsync(ctx context.Context, id TID) (*TEntity, error)
	Get(ctx context.Context, query interface{}, args ...interface{}) ([]*TEntity, error)
	Update(ctx context.Context, entity *TEntity) (*TEntity, error)
	UpdateRange(ctx context.Context, entities []*TEntity) error
	UnitOfWork() UnitOfWork
	GetPaged(ctx context.Context, query extensions.PagedQuery, condition interface{}, args ...interface{}) (*extensions.PagedResult[*TEntity], error) 
}
