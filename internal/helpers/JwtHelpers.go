package helpers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var sampleSecretKey = []byte("GoDiplomKey")

type JwtToken struct {
	Login *string
}

func DecodeJWT(tokenString string) (string, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", t.Header["alg"])
		}
		return sampleSecretKey, nil
	}
	var claims = jwt.RegisteredClaims{}

	parsedToken, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc)
	if err != nil {
		TLog.Error("Ошибка разбора: ", err)
		return "", err
	}

	if !parsedToken.Valid {
		err := "недействительный токен"
		TLog.Error(err)
		return "", fmt.Errorf(err)
	}
	return claims.Subject, nil

}

func generateJWTAndCheckError(login string, w http.ResponseWriter) string {
	token, err := generateJWT(login)
	if err != nil {
		TLog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return token
}

func CreateAndSetJWTCookieInHttp(login string, w http.ResponseWriter) {
	token := generateJWTAndCheckError(login, w)
	cookie := &http.Cookie{
		Name:     "Token",
		Value:    token,
		Secure:   false,
		HttpOnly: true,
		MaxAge:   300,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func generateJWT(username string) (string, error) {
	var claims jwt.RegisteredClaims
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 30).UTC())
	claims.Subject = username
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(sampleSecretKey)

	if err != nil {
		TLog.Error("Ошибка генерации токена: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}
