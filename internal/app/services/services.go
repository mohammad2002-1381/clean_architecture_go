package services

import (
	"context"
	"time"
)

type TokenService interface {
	GenerateToken(userID int32, role string) (string, error)
	GenerateRefreshToken(userID int32) (string, error)
	ValidateToken(token string) (*TokenClaims, error)
}

type TokenClaims struct {
	UserID    int32
	Role      string
	ExpiresAt time.Time
}

type PasswordService interface {
	Hash(password string) (string, error)
	Verify(password, hash string) bool
}

type CurrentUser struct {
	UserID   int32
	Email    string
	Role     string
	IsActive bool
}

type CurrentUserService interface {
	GetUserID() (int32, error)
	GetUserRole() (string, error)
	IsAdmin() bool
}

type Job func(ctx context.Context) error

type BackgroundJobService interface {
	Enqueue(job Job)
	EnqueueWithTimeout(job Job) bool
	PendingJobs() int
	Shutdown()
}
