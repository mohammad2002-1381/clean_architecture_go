package domain

import (
	"clean_architecture_go/internal/domain/events"
	"time"
)

type UserRoleType string

const (
	RoleAdmin UserRoleType = "admin"
	RoleUser  UserRoleType = "user"
)

type User struct {
	BaseEntity[int32]
	FirstName    string       `gorm:"column:first_name;not null"`
	LastName     string       `gorm:"column:last_name;not null"`
	Email        string       `gorm:"column:email;unique;not null"`
	PasswordHash string       `gorm:"column:password_hash;not null"`
	Role         UserRoleType `gorm:"column:role;default:user"`
	IsActive     bool         `gorm:"column:is_active;default:true"`
}

func NewUser(firstName, lastName, email string, role UserRoleType, passwordHash string, isActive bool) *User {
	now := time.Now().UTC()
	user := &User{
		BaseEntity: BaseEntity[int32]{
			CreatedAt: now,
			UpdatedAt: now,
		},
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		IsActive:     isActive,
	}

	// Add domain event
	user.AddDomainEvent(events.UserRegisteredEvent{
		UserID:    0, // Will be set after DB insert
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      string(role),
	}, events.UserRegisteredEventHandler {}) // Handler will be set by UnitOfWork

	return user
}

func (u *User) UpdateBeforeSave() {
	u.UpdatedAt = time.Now().UTC()
}
