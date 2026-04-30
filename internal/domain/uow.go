// domain/unit_of_work.go
package domain

import (
	"context"

	"gorm.io/gorm"
)

type DbCommand struct {
	Operation string      // "create", "update", "delete"
	Entity    interface{} // The entity to operate on
}

type UnitOfWork interface {
	// Queue operations (don't execute yet)
	Create(ctx context.Context, entity interface{})
	Save(ctx context.Context, entity interface{})
	Delete(ctx context.Context, entity interface{})

	// Execute all queued operations
	SaveEntitiesAsync(ctx context.Context) (int, error)

	// Query operations (execute immediately)
	First(ctx context.Context, dest interface{}, query interface{}, args ...interface{}) error
	Find(ctx context.Context, dest interface{}, query interface{}, args ...interface{}) error

	Count(ctx context.Context, model interface{}, count *int64) error
	Offset(offset int) UnitOfWork
	Limit(limit int) UnitOfWork
	Model(model interface{}) UnitOfWork
	Where(query interface{}, args ...interface{}) UnitOfWork
	GetDB() *gorm.DB
}
