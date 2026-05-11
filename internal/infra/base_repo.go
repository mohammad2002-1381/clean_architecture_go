package infra

import (
	"context"
	"errors"
	"reflect"
	"time"

	"go-ca/internal/domain"

	"gorm.io/gorm"
)

type BaseRepository[T domain.IBaseEntity[TID], TID comparable] struct {
	db *gorm.DB
}

// NewBaseRepository creates and returns a new *BaseRepository instance.
// It returns a pointer to ensure consistent behavior with method receivers
// and to allow chainable methods to return new instances that implement the interface.
func NewBaseRepository[T domain.IBaseEntity[TID], TID comparable](db *gorm.DB) BaseRepository[T, TID] {
	return BaseRepository[T, TID]{ // Return a pointer to the struct
		db: db,
	}
}

func (r *BaseRepository[T, TID]) setTimestamps(entity T, isInsert bool) {
	now := time.Now()
	val := reflect.ValueOf(entity).Elem()

	if isInsert {
		createdAtField := val.FieldByName("CreatedAt")
		if createdAtField.IsValid() && createdAtField.CanSet() {
			createdAtField.Set(reflect.ValueOf(now))
		}
	}

	updatedAtField := val.FieldByName("UpdatedAt")
	if updatedAtField.IsValid() && updatedAtField.CanSet() {
		updatedAtField.Set(reflect.ValueOf(now))
	}
}

func (r *BaseRepository[T, TID]) Add(ctx context.Context, entity T) error {
	// 1: Insertion => Set both CreatedAt and UpdatedAt
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return err
	}
	r.saveEntityAsync(ctx, entity, true)
	return nil
}

func (r *BaseRepository[T, TID]) Update(ctx context.Context, entity T) error {
	// 2: Update => Set only UpdatedAt
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		return err
	}
	r.saveEntityAsync(ctx, entity, false)
	return nil
}

func (r *BaseRepository[T, TID]) Find(ctx context.Context, id TID) (T, error) {
	var entity T
	// r.db here will be the current session, which might have filters applied if chained.
	err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, nil
		}
		return entity, err
	}
	return entity, nil
}

func (r *BaseRepository[T, TID]) First(ctx context.Context) (T, error) {
	var entity T
	// r.db here will be the current session, which might have filters applied if chained.
	err := r.db.WithContext(ctx).First(&entity).Error
	return entity, err
}

func (r *BaseRepository[T, TID]) Get(ctx context.Context) ([]T, error) {
	var entities []T
	// r.db here will be the current session, which might have filters applied if chained.
	err := r.db.WithContext(ctx).Find(&entities).Error
	return entities, err
}

func (r *BaseRepository[T, TID]) Delete(ctx context.Context, id TID) error {
	var entity T
	// r.db here will be the current session, which might have filters applied if chained.
	return r.db.WithContext(ctx).Delete(&entity, "id = ?", id).Error
}

// Where returns a NEW BaseRepository instance with the WHERE condition applied.
// The original repository remains unchanged (immutable).
func (r *BaseRepository[T, TID]) Where(ctx context.Context, filter T) domain.BaseRepository[T, TID] {
	// Create a new GORM session with the context and the filter.
	newGormDB := r.db.WithContext(ctx).Where(filter)
	// Return a *NEW* BaseRepository instance that encapsulates this new GORM session.
	// This ensures immutability for the original 'r' and allows chaining.
	return &BaseRepository[T, TID]{db: newGormDB}
}

func (r *BaseRepository[T, TID]) saveEntityAsync(ctx context.Context, entity T, inserted bool) {
	r.setTimestamps(entity, inserted)

	events := entity.GetNotifications()
	for _, event := range events {
		domain.Dispatch(ctx, event)
	}
	entity.ClearNotifications()
}

// Or returns a NEW BaseRepository instance with the OR condition applied.
// The original repository remains unchanged (immutable).
func (r *BaseRepository[T, TID]) Or(ctx context.Context, filter T) domain.BaseRepository[T, TID] {
	// Create a new GORM session with the context and the OR filter.
	newGormDB := r.db.WithContext(ctx).Or(filter)
	// Return a *NEW* BaseRepository instance that encapsulates this new GORM session.
	return &BaseRepository[T, TID]{db: newGormDB}
}

func (r *BaseRepository[T, TID]) GetPaged(ctx context.Context, params domain.PaginationParams) (domain.PaginatedResult[T], error) {
	var items []T
	var total int64

	// r.db here will be the current session, which might have filters applied if chained.
	db := r.db.WithContext(ctx).Model(new(T))

	if err := db.Count(&total).Error; err != nil {
		return domain.PaginatedResult[T]{}, err
	}

	if params.Sort != "" {
		order := params.Sort
		if params.Desc {
			order += " desc"
		}
		db = db.Order(order)
	}

	if !params.DisablePaging {
		if params.PageNumber < 1 {
			params.PageNumber = 1
		}
		if params.PageSize < 1 {
			params.PageSize = 10
		}

		offset := (params.PageNumber - 1) * params.PageSize
		db = db.Offset(offset).Limit(params.PageSize)
	}

	if err := db.Find(&items).Error; err != nil {
		return domain.PaginatedResult[T]{}, err
	}

	return domain.PaginatedResult[T]{
		Items:      items,
		TotalCount: total,
		PageNumber: params.PageNumber,
		PageSize:   params.PageSize,
	}, nil
}
