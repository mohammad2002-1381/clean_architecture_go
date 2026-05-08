package domain

import (
	"context"
	"reflect"
)

type Notification interface {
	IsNotification()
}

type NotificationHandler[T Notification] interface {
	Handle(ctx context.Context, event T) error
}

type EventHandler func(ctx context.Context, event Notification) error

var eventRegistry = make(map[reflect.Type][]EventHandler)

func RegisterHandler[T Notification](handler NotificationHandler[T]) {
	eventType := reflect.TypeOf((*T)(nil)).Elem()

	wrapper := func(ctx context.Context, ev Notification) error {
		return handler.Handle(ctx, ev.(T))
	}

	eventRegistry[eventType] = append(eventRegistry[eventType], wrapper)
}

func Dispatch(ctx context.Context, event Notification) error {
	eventType := reflect.TypeOf(event)

	handlers, exists := eventRegistry[eventType]
	if !exists {
		return nil
	}

	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			return err
		}
	}
	return nil
}
