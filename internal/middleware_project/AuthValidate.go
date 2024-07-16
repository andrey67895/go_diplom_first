package middleware_project

import (
	"net/http"
	"slices"

	"github.com/andrey67895/go_diplom_first/internal/helpers"
)

var not_auth_url_path = []string{"/api/user/register", "/api/user/login"}

func itNotAuth(path string) bool {
	return !slices.Contains(not_auth_url_path, path)
}

func AuthValidate(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if itNotAuth(r.URL.Path) {
			cookie, err := r.Cookie("Token")
			if err != nil {
				helpers.TLog.Error(err.Error() + " : пользователь не аутентифицирован!")
				http.Error(w, "Пользователь не аутентифицирован!", http.StatusUnauthorized)
				return
			}
			_, err = helpers.DecodeJWT(cookie.Value)
			if err != nil {
				helpers.TLog.Error(err.Error() + " : пользователь не аутентифицирован!")
				http.Error(w, "Пользователь не аутентифицирован!", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
