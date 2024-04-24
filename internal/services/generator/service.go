package generator

import (
	"fmt"

	"example.com/m/internal/user"
)

type Tokens struct {
	Access   string
	Refrersh string
}

func NewTokens(access string, refresh string) *Tokens {
	return &Tokens{Access: access, Refrersh: refresh}
}

func GetTokens(user user.User, createTime int64) (t Tokens, u user.CreateUserDTO, err error) {
	jwt, time, err := NewJWT(user, createTime)
	if err != nil {
		return t, u, fmt.Errorf("error: %v", err)
	}

	refresh, err := NewRefreshToken()
	if err != nil {
		return t, u, fmt.Errorf("error: %v", err)
	}

	u.GUID = user.GUID
	u.AccessIssuedAt = time
	u.RefreshToken = refresh

	return *NewTokens(jwt, refresh), u, nil
}
