package services

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/andrey67895/go_diplom_first/internal/database"
	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/model"
)

func GetAuth(login string, w http.ResponseWriter, create bool) *model.UserModel {
	auth, err := database.DBStorage.GetAuth(login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if create {
				helpers.CreateAndSetJWTCookieInHTTP(login, w)
			} else {
				fail := "неверная пара логин/пароль"
				helpers.TLog.Error(fail)
				http.Error(w, fail, http.StatusUnauthorized)
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	if auth != nil && create {
		http.Error(w, "Пользователь уже существует", http.StatusConflict)
	} else if auth == nil && !create {
		http.Error(w, "неверная пара логин/пароль", http.StatusUnauthorized)
	}
	return auth
}
