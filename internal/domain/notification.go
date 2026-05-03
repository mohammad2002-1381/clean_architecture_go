package domain

import "context"

type Notification interface {
	EventHandler(context.Context)
}