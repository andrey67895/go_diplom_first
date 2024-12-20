package middlewareproject

import (
	"net/http"
	"slices"

	"github.com/andrey67895/go_diplom_first/internal/helpers"
)

var notAuthURLPath = []string{"/api/user/register", "/api/user/login", "/api/ping", "/apidocs", "/openapi.yaml"}

func itCheckAuth(path string) bool {
	return !slices.Contains(notAuthURLPath, path)
}

func AuthValidate(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if itCheckAuth(r.URL.Path) {
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
