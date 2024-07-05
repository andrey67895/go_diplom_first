package transport

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/andrey67895/go_diplom_first/internal/database"
	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/model"
)

func SaveOrders(w http.ResponseWriter, req *http.Request) {
	b, err := io.ReadAll(req.Body)
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, "Неверный формат номера заказа!", http.StatusUnprocessableEntity)
		return
	}
	orderId, err := strconv.Atoi(string(b))
	if !helpers.LuhnValid(orderId) || err != nil {
		http.Error(w, "Неверный формат номера заказа!", http.StatusUnprocessableEntity)
		return
	}

}

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

	auth, err := database.DBStorage.GetAuth(*tModel.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := database.DBStorage.CreateAuth(tModel)
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

func AuthUser(w http.ResponseWriter, req *http.Request) {
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

	auth, err := database.DBStorage.GetAuth(*tModel.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.TLog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	if auth != nil {
		if *auth.Password == helpers.EncodeHash(*tModel.Password) {
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
			http.Error(w, "неверная пара логин/пароль", http.StatusUnauthorized)
			return
		}
	}

}
