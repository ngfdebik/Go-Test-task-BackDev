package generator

import (
	"os"
	"strconv"

	"example.com/m/internal/user"
	jwt "github.com/dgrijalva/jwt-go"
)

func NewJWT(user user.User, timer int64) (t string, time string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:  user.GUID,
		IssuedAt: timer,
	})
	time = strconv.Itoa(int(timer))
	t, err = token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "JWT token signing", "", err
	}

	return t, time, err
}
