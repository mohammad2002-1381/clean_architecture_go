package dto

type AuthDto struct {
	User         UserDTO	`json:"user"`
	Token        string		`json:"token"`
	RefreshToken string		`json:"refresh_token"`
}

func NewAuthDto(user UserDTO, token string, refreshToken string) AuthDto {
	return AuthDto{
		User:         user,
		Token:        token,
		RefreshToken: refreshToken,
	}
}
