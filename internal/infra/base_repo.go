package infra

import (
	"context"
	"errors"
	"go-ca/internal/domain"

	"gorm.io/gorm"
)

type BaseRepository[T any, TID comparable] struct {
	db *gorm.DB
}

func NewBaseRepository[T any, TID comparable](db *gorm.DB) BaseRepository[T, TID] {
	return BaseRepository[T, TID]{
		db: db,
	}
}

func (r *BaseRepository[T, TID]) dispatchEvents(ctx context.Context, entity *T) {
	// Cast the entity to IBaseEntity to access notifications
	// if baseEntity, ok := any(entity).(domain.IBaseEntity[TID]); ok {
	// 	payloads := baseEntity.GetNotifications()

	// 	if len(payloads) > 0 {
	// 		allSuccess := true

	// 		for _, payload := range payloads {
	// 			// Execute the handler attached to the notification
	// 			if err := payload.Handler(ctx, payload.Notification); err != nil {
	// 				// In a production app, you might want to log this error
	// 				allSuccess = false
	// 			}
	// 		}

	// 		// Clear notifications if all handlers succeeded
	// 		// (matches your previous dispatcher logic)
	// 		if allSuccess {
	// 			baseEntity.ClearNotifications()
	// 		}
	// 	}
	// }
}

func (r *BaseRepository[T, TID]) Add(ctx context.Context, entity *T) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return err
	}
	r.dispatchEvents(ctx, entity)
	return nil
}

// func (r *BaseRepository[T, TID]) saveEntityAsync(ctx context.Context, entity *T) error {
// 	entity
// }

func (r *BaseRepository[T, TID]) Update(ctx context.Context, entity *T) error {
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		return err
	}
	r.dispatchEvents(ctx, entity)
	return nil
}

func (r *BaseRepository[T, TID]) Find(ctx context.Context, id TID) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T, TID]) First(ctx context.Context) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity).Error
	return &entity, err
}

func (r *BaseRepository[T, TID]) Get(ctx context.Context) ([]*T, error) {
	var entities []*T
	err := r.db.WithContext(ctx).Find(&entities).Error
	return entities, err
}

func (r *BaseRepository[T, TID]) Delete(ctx context.Context, id TID) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, "id = ?", id).Error
}

func (r *BaseRepository[T, TID]) Where(ctx context.Context, filter *T) domain.BaseRepository[T, TID] {
	r.db = r.db.WithContext(ctx).Where(filter)
	return r
}

func (r *BaseRepository[T, TID]) Or(ctx context.Context, filter *T) domain.BaseRepository[T, TID] {
	r.db = r.db.WithContext(ctx).Or(filter)
	return r
}

func (r *BaseRepository[T, TID]) GetPaged(ctx context.Context, params domain.PaginationParams) (domain.PaginatedResult[T], error) {
	var items []T
	var total int64

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
