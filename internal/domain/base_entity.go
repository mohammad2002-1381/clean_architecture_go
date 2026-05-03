package domain

import (
	"time"
)

type IBaseEntity[TID comparable] interface {
	SetID(value TID)
	SetCreatedAt(value time.Time)
	SetUpdatedAt(value time.Time)
	AddNotification(n Notification)
	GetNotifications() []Notification
	ClearNotifications()
}

type BaseEntity[TID comparable] struct {
	ID        TID       `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	notifications []Notification `gorm:"-" json:"-"`
}

func (b *BaseEntity[TID]) SetID(value TID) {
	b.ID = value
}

func (b *BaseEntity[TID]) SetCreatedAt(value time.Time) {
	b.CreatedAt = value
}

func (b *BaseEntity[TID]) SetUpdatedAt(value time.Time) {
	b.UpdatedAt = value
}

func (b *BaseEntity[TID]) AddNotification(n Notification) {
	b.notifications = append(b.notifications, n)
}

func (b *BaseEntity[TID]) GetNotifications() []Notification {
	return b.notifications
}

func (b *BaseEntity[TID]) ClearNotifications() {
	b.notifications = nil
}
