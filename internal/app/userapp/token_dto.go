package userapp

type TokenDTO struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func NewTokenDTO(token, refreshToken string) TokenDTO {
	return TokenDTO{
		Token:        token,
		RefreshToken: refreshToken,
	}
}
