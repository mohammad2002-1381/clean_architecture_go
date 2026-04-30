// infra/base_repository.go
package infra

import (
	"clean_architecture_go/internal/domain"
	"clean_architecture_go/internal/domain/extensions"
	"context"
	"errors"

	"gorm.io/gorm"
)

type BaseRepository[TEntity domain.IBaseEntity[TID], TID comparable] struct {
	uow domain.UnitOfWork
}

func NewBaseRepository[TEntity domain.IBaseEntity[TID], TID comparable]() *BaseRepository[TEntity, TID] {
	return &BaseRepository[TEntity, TID]{
		uow: NewUnitOfWork(),
	}
}

func (r *BaseRepository[TEntity, TID]) UnitOfWork() domain.UnitOfWork {
	return r.uow
}

func (r *BaseRepository[TEntity, TID]) Add(ctx context.Context, entity *TEntity) (*TEntity, error) {
	r.uow.Create(ctx, entity)
	return entity, nil
}

func (r *BaseRepository[TEntity, TID]) AddRange(ctx context.Context, entities []*TEntity) error {
	for _, entity := range entities {
		r.uow.Create(ctx, entity)
	}
	return nil
}

func (r *BaseRepository[TEntity, TID]) Update(ctx context.Context, entity *TEntity) (*TEntity, error) {
	r.uow.Save(ctx, entity)
	return entity, nil
}

func (r *BaseRepository[TEntity, TID]) UpdateRange(ctx context.Context, entities []*TEntity) error {
	for _, entity := range entities {
		r.uow.Save(ctx, entity)
	}
	return nil
}

// Delete queues a single entity for deletion
func (r *BaseRepository[TEntity, TID]) Delete(ctx context.Context, entity *TEntity) error {
	r.uow.Delete(ctx, entity)
	return nil
}

func (r *BaseRepository[TEntity, TID]) DeleteRange(ctx context.Context, entities []*TEntity) error {
	for _, entity := range entities {
		r.uow.Delete(ctx, entity)
	}
	return nil
}

func (r *BaseRepository[TEntity, TID]) FindAsync(ctx context.Context, id TID) (*TEntity, error) {
	var entity TEntity
	err := r.uow.First(ctx, &entity, "id = ?", id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[TEntity, TID]) Get(ctx context.Context, query interface{}, args ...interface{}) ([]*TEntity, error) {
	var entities []*TEntity
	err := r.uow.Find(ctx, &entities, query, args...)
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *BaseRepository[TEntity, TID]) GetPaged(ctx context.Context, query extensions.PagedQuery, condition interface{}, args ...interface{}) (*extensions.PagedResult[*TEntity], error) {
	query.Normalize()

	db := r.uow.GetDB().WithContext(ctx).Model(new(TEntity))

	if condition != nil {
		db = db.Where(condition, args...)
	}

	var totalCount int64
	if err := db.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if !query.DisablePaging {
		db = db.Offset(query.Offset()).Limit(query.Limit())
	}

	var items []*TEntity
	if err := db.Find(&items).Error; err != nil {
		return nil, err
	}

	result := extensions.NewPagedResult(items, int(totalCount), query)
	return &result, nil
}