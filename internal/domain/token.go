package domain

import "time"

type Token struct {
	BaseEntity[uint]
	Token        string    `gorm:"column:token;not null"`
	RefreshToken string    `gorm:"column:refresh_token;not null"`
	Expires      time.Time `gorm:"column:expires;not null"`
	UserID       uint     `gorm:"column:user_id;not null;index"`
	User         *User     `gorm:"foreignKey:UserID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
}

func NewToken(token, refreshToken string, userID uint) *Token {
	now := time.Now().UTC()
	return &Token{
		BaseEntity: BaseEntity[uint]{
			CreatedAt: now,
			UpdatedAt: now,
		},
		Token:        token,
		RefreshToken: refreshToken,
		UserID:       userID,
		Expires:      now.AddDate(0, 0, 1),
	}
}

func (t *Token) IsActive() bool {
	return t.Expires.After(time.Now().UTC())
}

func (t *Token) SetToken(value string) {
	t.Token = value
}

func (t *Token) SetRefreshToken(value string) {
	t.RefreshToken = value
}
