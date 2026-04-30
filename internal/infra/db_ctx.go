// infra/unit_of_work.go
package infra

import (
	"clean_architecture_go/internal/domain"
	"context"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

type UnitOfWork struct {
	db       *gorm.DB
	commands []domain.DbCommand
	tx       *gorm.DB
}

func NewUnitOfWork() domain.UnitOfWork {
	return &UnitOfWork{
		db:       dbInstance,
		commands: make([]domain.DbCommand, 0),
	}
}

func (u *UnitOfWork) Create(ctx context.Context, entity interface{}) {
	u.commands = append(u.commands, domain.DbCommand{
		Operation: "create",
		Entity:    entity,
	})
}

func (u *UnitOfWork) Save(ctx context.Context, entity interface{}) {
	u.commands = append(u.commands, domain.DbCommand{
		Operation: "save",
		Entity:    entity,
	})
}

func (u *UnitOfWork) Delete(ctx context.Context, entity interface{}) {
	u.commands = append(u.commands, domain.DbCommand{
		Operation: "delete",
		Entity:    entity,
	})
}

func (u *UnitOfWork) SaveEntitiesAsync(ctx context.Context) (int, error) {
	if len(u.commands) == 0 {
		return 0, nil
	}

	now := time.Now().UTC()
	u.setTimestamps(now)

	tx := u.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	count := 0
	for _, cmd := range u.commands {
		var err error

		switch cmd.Operation {
		case "create":
			err = tx.Create(cmd.Entity).Error
		case "save":
			err = tx.Save(cmd.Entity).Error
		case "delete":
			err = tx.Delete(cmd.Entity).Error
		}

		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("command %s failed: %w", cmd.Operation, err)
		}
		count++
	}

	if err := tx.Commit().Error; err != nil {
		return 0, fmt.Errorf("failed to commit: %w", err)
	}

	if err := u.dispatchDomainEvents(ctx); err != nil {
		return count, fmt.Errorf("failed to dispatch domain events: %w", err)
	}

	u.commands = make([]domain.DbCommand, 0)

	return count, nil
}

func (u *UnitOfWork) setTimestamps(now time.Time) {
	for _, cmd := range u.commands {
		if cmd.Operation == "create" || cmd.Operation == "save" {
			if baseEntity, ok := cmd.Entity.(interface{ SetTimestamps(time.Time) }); ok {
				baseEntity.SetTimestamps(now)
			}
		}
	}
}

func (u *UnitOfWork) dispatchDomainEvents(ctx context.Context) error {
	for _, cmd := range u.commands {
		if eventSource, ok := cmd.Entity.(*domain.BaseEntity[any]); ok {
			if err := eventSource.DispatchEvents(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *UnitOfWork) First(ctx context.Context, dest interface{}, query interface{}, args ...interface{}) error {
	return u.db.WithContext(ctx).Where(query, args...).First(dest).Error
}

func (u *UnitOfWork) Find(ctx context.Context, dest interface{}, query interface{}, args ...interface{}) error {
	db := u.db.WithContext(ctx)
	if query != nil {
		db = db.Where(query, args...)
	}
	return db.Find(dest).Error
}

func (u *UnitOfWork) Count(ctx context.Context, model interface{}, count *int64) error {
	return u.db.WithContext(ctx).Model(model).Count(count).Error
}

func (u *UnitOfWork) Offset(offset int) domain.UnitOfWork {
	u.db = u.db.Offset(offset)
	return u
}

func (u *UnitOfWork) Limit(limit int) domain.UnitOfWork {
 (u *UnitOfWork) GetDB() *gorm.DB {
	return u.db
}

func (u *UnitOfWork) Where(query interface{}, args ...interface{}) domain.UnitOfWork {
	u.db = u.db.Where(query, args...)
	return u
}
	u.db = u.db.Limit(limit)
	return u
}

func (u *UnitOfWork) Model(model interface{}) domain.UnitOfWork {
	u.db = u.db.Model(model)
	return u
}

func