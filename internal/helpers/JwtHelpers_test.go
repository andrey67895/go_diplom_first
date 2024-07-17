package helpers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateJWTAndCheckError(t *testing.T) {
	login := "TEST_LOGIN"
	type args struct {
		login string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive test #1",
			args: args{
				login: login,
			},
			want: login,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			token := GenerateJWTAndCheckError(test.args.login, w)
			login, err := DecodeJWT(token)
			assert.Equal(t, nil, err)
			assert.Equal(t, test.want, login, "GenerateJWTAndCheckError is not assert")
		})
	}
}

func TestCreateAndSetJWTCookieInHTTP(t *testing.T) {
	login := "TEST_LOGIN"
	type args struct {
		login string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test #1",
			args: args{
				login: login,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			CreateAndSetJWTCookieInHTTP(test.args.login, w)
			result := w.Result()
			for _, cookie := range result.Cookies() {
				if cookie.Name == "Token" {
					token := cookie.Value
					assert.NotNil(t, token)
					login, err := DecodeJWT(token)
					assert.Equal(t, nil, err)
					assert.Equal(t, test.args.login, login, "GenerateJWTAndCheckError is not assert")
					break
				}
			}
			result.Body.Close()
		})
	}
}

func TestDecodeJWT(t *testing.T) {
	type args struct {
		tokenString string
	}
	type want struct {
		login string
		err   string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "negative test #1",
			args: args{
				tokenString: "ERROR_TOKEN",
			},

			want: want{
				login: "",
				err:   "ошибка разбора: token contains an invalid number of segments",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tLogin, err := DecodeJWT(test.args.tokenString)

			assert.NotNil(t, err)
			assert.Truef(t, strings.Contains(err.Error(), test.want.err), "Ошибка DecodeJWT(%v)", test.args.tokenString)
			assert.Equal(t, test.want.login, tLogin, "Ошибка DecodeJWT(%v)", test.args.tokenString)
		})
	}
}
