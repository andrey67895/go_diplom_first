package middlewareproject

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

var notAuthPath = []string{"/api/user/register", "/api/user/login"}
var authPath = []string{"/api/user/orders", "/api/user/balance", "/api/user/balance/withdraw", "/api/user/withdrawals"}

func Test_itAuth(t *testing.T) {
	type args struct {
		paths []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "positive test #1",
			args: args{
				paths: notAuthPath,
			},
			want: false,
		},
		{
			name: "positive test #2",
			args: args{
				paths: authPath,
			},
			want: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, path := range test.args.paths {
				assert.Equalf(t, test.want, itCheckAuth(path), "Ошибка в методе: %s", path)
			}
		})
	}
}

func TestAuthValidate(t *testing.T) {
	type args struct {
		paths []string
		auth  bool
		valid bool
	}
	type want struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive test #1",
			args: args{
				paths: notAuthPath,
				auth:  false,
				valid: false,
			},
			want: want{
				code: 200,
			},
		},
		{
			name: "positive test #2",
			args: args{
				paths: notAuthPath,
				auth:  true,
				valid: false,
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "positive test #3",
			args: args{
				paths: notAuthPath,
				auth:  false,
				valid: true,
			},
			want: want{
				code: 200,
			},
		},
		{
			name: "positive test #4",
			args: args{
				paths: notAuthPath,
				auth:  true,
				valid: true,
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "positive test #5",
			args: args{
				paths: authPath,
				auth:  true,
				valid: true,
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "negative test #1",
			args: args{
				paths: authPath,
				auth:  false,
				valid: false,
			},
			want: want{
				code: http.StatusUnauthorized,
			},
		},
		{
			name: "negative test #2",
			args: args{
				paths: authPath,
				auth:  false,
				valid: true,
			},
			want: want{
				code: http.StatusUnauthorized,
			},
		},
		{
			name: "negative test #3",
			args: args{
				paths: authPath,
				auth:  true,
				valid: false,
			},
			want: want{
				code: http.StatusUnauthorized,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, path := range test.args.paths {
				req := httptest.NewRequest("GET", path, nil)
				w := httptest.NewRecorder()
				r := chi.NewRouter()
				r.Use(AuthValidate)
				if test.args.auth {
					var token string
					if test.args.valid {
						token = helpers.GenerateJWTAndCheckError("TEST", w)
					} else {
						token = "QWERTY"
					}
					req.AddCookie(&http.Cookie{Name: "Token", Value: token})
				}
				fn := func(w http.ResponseWriter, req *http.Request) {}
				tHandler := AuthValidate(http.HandlerFunc(fn))

				tHandler.ServeHTTP(w, req)

				result := w.Result()
				assert.Equal(t, test.want.code, result.StatusCode)
			}
		})
	}
}
