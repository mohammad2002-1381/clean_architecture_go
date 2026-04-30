package jwtservice

import (
	"clean_architecture_go/internal/app/services"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenService struct {
	secretKey       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewJWTTokenService(secret string) services.TokenService {
	println("secret: " + secret)
	if secret == "" {
		secret = "my-hardcoded-secret-key-123456" // Hardcoded for testing
	}
	println("secret: " + secret)
	return &JWTTokenService{
		secretKey:       []byte(secret),
		accessTokenTTL:  15 * time.Minute,
		refreshTokenTTL: 7 * 24 * time.Hour,
	}
}

func (s *JWTTokenService) GenerateToken(userID int32, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(s.accessTokenTTL).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *JWTTokenService) GenerateRefreshToken(userID int32) (string, error) {
	bytes := make([]byte, 64)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *JWTTokenService) ValidateToken(tokenString string) (*services.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	userID, _ := claims["user_id"].(float64)
	role, _ := claims["role"].(string)
	exp, _ := claims["exp"].(float64)

	return &services.TokenClaims{
		UserID:    int32(userID),
		Role:      role,
		ExpiresAt: time.Unix(int64(exp), 0),
	}, nil
}
