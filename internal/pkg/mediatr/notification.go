package mediatr

import (
	"clean_architecture_go/internal/pkg/error"
	"context"
)

type Notification interface{}

type NotificationHandler[TNotification Notification] interface {
	Handle(ctx context.Context, request TNotification) *error.RequestError
}
