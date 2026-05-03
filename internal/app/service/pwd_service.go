package service

import (
	"golang.org/x/crypto/bcrypt"
)

type bcryptPasswordService struct {
	cost int
}

func newbcryptPasswordService() *bcryptPasswordService {
	return &bcryptPasswordService{
		cost: bcrypt.DefaultCost, // 10
	}
}

func NewbcryptPasswordServiceWithCost(cost int) *bcryptPasswordService {
	return &bcryptPasswordService{
		cost: cost,
	}
}

func (s *bcryptPasswordService) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s *bcryptPasswordService) Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}