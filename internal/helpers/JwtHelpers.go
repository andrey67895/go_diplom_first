package helpers

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var sampleSecretKey = []byte("GoDiplomKey")

func GenerateJWT(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(sampleSecretKey)

	if err != nil {
		TLog.Error("Ошибка генерации токена: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}
