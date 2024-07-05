package transport

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/andrey67895/go_diplom_first/internal/database"
	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/model"
)

func UserRegister(w http.ResponseWriter, req *http.Request) {
	var tModel model.UserModel
	err := json.NewDecoder(req.Body).Decode(&tModel)
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, "Ошибка десериализации!", http.StatusBadRequest)
		return
	}
	err = tModel.IsValid()
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	auth, err := database.DbStorage.GetAuth(*tModel.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := database.DbStorage.CreateAuth(tModel)
			if err != nil {
				helpers.TLog.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			token, err := helpers.GenerateJWT(*tModel.Login)
			if err != nil {
				helpers.TLog.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			cookie := &http.Cookie{
				Name:     "Token",
				Value:    token,
				Secure:   false,
				HttpOnly: true,
				MaxAge:   300,
			}
			http.SetCookie(w, cookie)
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	if auth != nil {
		http.Error(w, "Пользователь уже существует", http.StatusConflict)
		return
	}

}
