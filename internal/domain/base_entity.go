// domain/base_entity.go
package domain

import (
	"clean_architecture_go/internal/pkg/mediatr"
	"context"
	"sync"
	"time"
)

type IBaseEntity[TID comparable] interface {
	verifyEntity()
}

type DomainEventEntry struct {
	Event   mediatr.Notification
	Handler mediatr.NotificationHandler[mediatr.Notification]
}

type BaseEntity[TID comparable] struct {
	ID        TID       `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	domainEvents []DomainEventEntry `gorm:"-" json:"-"`
	mu           sync.Mutex         `gorm:"-" json:"-"`
}

func (b BaseEntity[TID]) verifyEntity() {}

func (b *BaseEntity[TID]) AddDomainEvent(event mediatr.Notification, handler mediatr.NotificationHandler[mediatr.Notification]) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.domainEvents = append(b.domainEvents, DomainEventEntry{
		Event:   event,
		Handler: handler,
	})
}

func (b *BaseEntity[TID]) RemoveDomainEvent(event mediatr.Notification) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i, e := range b.domainEvents {
		if e.Event == event {
			b.domainEvents = append(b.domainEvents[:i], b.domainEvents[i+1:]...)
			return
		}
	}
}

func (b *BaseEntity[TID]) ClearDomainEvents() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.domainEvents = nil
}

func (b *BaseEntity[TID]) DispatchEvents(ctx context.Context) error {
	b.mu.Lock()
	events := b.domainEvents
	b.domainEvents = nil
	b.mu.Unlock()

	for _, entry := range events {
		if err := entry.Handler.Handle(ctx, entry.Event); err != nil {
			return err
		}
	}
	return nil
}

func (b *BaseEntity[TID]) SetTimestamps(now time.Time) {
	if b.CreatedAt.IsZero() {
		b.CreatedAt = now
	}
	b.UpdatedAt = now
}
