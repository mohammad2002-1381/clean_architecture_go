package domain

import (
	"time"
)

type Token struct {
	BaseEntity[int32]
	Token        string    `gorm:"column:token;not null"`
	RefreshToken string    `gorm:"column:refresh_token;not null"`
	Expires      time.Time `gorm:"column:expires;not null"`
	UserID       int32     `gorm:"column:user_id;not null;index"`
	User         *User     `gorm:"foreignKey:UserID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
}

func NewToken(token, refreshToken string, userID int32) *Token {
	now := time.Now().UTC()
	return &Token{
		BaseEntity: BaseEntity[int32]{
			CreatedAt: now,
			UpdatedAt: now,
		},
		Token:        token,
		RefreshToken: refreshToken,
		UserID:       userID,
		Expires:      now.AddDate(0, 1, 0),
	}
}

func (t *Token) IsActive() bool {
	return t.Expires.After(time.Now().UTC())
}