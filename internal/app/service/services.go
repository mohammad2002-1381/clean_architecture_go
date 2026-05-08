package service

import (
	"context"
	"go-ca/internal/domain"
	"time"
)

type TokenClaims struct {
	UserID    uint
	Role      string
	ExpiresAt time.Time
}

type JWTService interface {
	GenerateToken(uint, string) (string, error)
	GenerateRefreshToken(uint) (string, error)
	ValidateToken(string) (*TokenClaims, error)
}

type PasswordService interface {
	Hash(string) (string, error)
	Verify(string, string) bool
}

type CurrentUserService interface {
	GetUserID(context.Context) (uint, error)
	GetUserRole(context.Context) (domain.UserRoleType, error)
	GetToken(context.Context) (string, error)
}

func NewJwtService(secret string) JWTService {
	return newJWTTokenService(secret)
}

func NewPasswordService() PasswordService {
	return newbcryptPasswordService()
}
