package passwordservice

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordService struct {
	cost int
}

func NewBcryptPasswordService() *BcryptPasswordService {
	return &BcryptPasswordService{
		cost: bcrypt.DefaultCost, // 10
	}
}

func NewBcryptPasswordServiceWithCost(cost int) *BcryptPasswordService {
	return &BcryptPasswordService{
		cost: cost,
	}
}

func (s *BcryptPasswordService) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s *BcryptPasswordService) Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}