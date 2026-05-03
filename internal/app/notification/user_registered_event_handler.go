package notification

import "context"

type UserRegisteredEvent struct {
	UserID   uint
	OldEmail string
	NewEmail string
}

func (u *UserRegisteredEvent) EventHandler(ctx context.Context) {
	println("user registered successfully")
}
