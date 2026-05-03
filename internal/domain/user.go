package domain

type UserRoleType string

const (
	RoleAdmin UserRoleType = "admin"
	RoleUser  UserRoleType = "user"
)

type User struct {
	BaseEntity[uint]
	FirstName    string       `gorm:"column:first_name;not null"`
	LastName     string       `gorm:"column:last_name;not null"`
	Email        string       `gorm:"column:email;unique;not null"`
	PasswordHash string       `gorm:"column:password_hash;not null"`
	Role         UserRoleType `gorm:"column:role;default:user"`
	IsActive     bool         `gorm:"column:is_active;default:true"`
}

func NewUser(firstName, lastName, email, passwordHash string, role UserRoleType, isActive bool) User {
	return User{
		FirstName: firstName,
		LastName: lastName,
		Email: email,
		PasswordHash: passwordHash,
		Role: role,
		IsActive: isActive,
	}
}

func (u *User) ActiveUser() {
	u.IsActive = true
}

func (u *User) DisactiveUser() {
	u.IsActive = false
}

func (u *User) SetFirstName(value string) {
	u.FirstName = value
}

func (u *User) SetLastName(value string) {
	u.LastName = value
}

func (u *User) SetEmail(value string) {
	u.Email = value
}

func (u *User) SetRole(value UserRoleType) {
	u.Role = value
}