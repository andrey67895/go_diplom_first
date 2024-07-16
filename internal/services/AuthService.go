package services

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/andrey67895/go_diplom_first/internal/database"
	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/model"
)

func GetAuth(login string, w http.ResponseWriter) *model.UserModel {
	auth, err := database.DBStorage.GetAuth(login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fail := "неверная пара логин/пароль"
			helpers.TLog.Error(fail)
			http.Error(w, fail, http.StatusUnauthorized)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	return auth
}
