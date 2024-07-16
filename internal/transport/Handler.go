package transport

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/andrey67895/go_diplom_first/internal/database"
	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/model"
)

func GetWithdrawalsHistory(w http.ResponseWriter, req *http.Request) {
	cookie, _ := req.Cookie("Token")
	login, _ := helpers.DecodeJWT(cookie.Value)

	withdrawnHistory, err := database.DBStorage.GetWithdrawnBalanceByLogin(login)
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, "Ошибка сервера!", http.StatusInternalServerError)
		return

	}
	w.Header().Set("Content-Type", "application/json")
	if len(*withdrawnHistory) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	tOrders := *withdrawnHistory

	sort.Slice(tOrders, func(i, j int) bool {
		return tOrders[i].ProcessedAT.After(*tOrders[j].ProcessedAT)
	})

	marshal, err := json.Marshal(tOrders)
	if err != nil {
		http.Error(w, "Ошибка записи ответа", http.StatusNotFound)
		return
	}
	_, errWrite := w.Write(marshal)
	if errWrite != nil {
		return
	}
	w.WriteHeader(http.StatusOK)

}

func WithdrawBalance(w http.ResponseWriter, req *http.Request) {
	cookie, _ := req.Cookie("Token")
	login, _ := helpers.DecodeJWT(cookie.Value)
	var tModel model.WithdrawnBalanceModel
	err := json.NewDecoder(req.Body).Decode(&tModel)
	if err != nil {
		helpers.TLog.Error(err.Error())
		return
	}
	orderID, err := strconv.Atoi(*tModel.Order)
	if !helpers.LuhnValid(orderID) || err != nil {
		http.Error(w, "Неверный формат номера заказа!", http.StatusUnprocessableEntity)
		return
	}

	currentBalanceModel, err := database.DBStorage.GetCurrentBalanceByLogin(login)
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if *currentBalanceModel.Balance < *tModel.Withdrawn {
		http.Error(w, "На счету недостаточно средств", http.StatusPaymentRequired)
		return
	}
	err = database.DBStorage.WithdrawnBalanceSumByLogin(model.WithdrawnBalanceModel{Login: &login, Order: tModel.Order, Withdrawn: tModel.Withdrawn})
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetBalance(w http.ResponseWriter, req *http.Request) {
	cookie, _ := req.Cookie("Token")
	login, _ := helpers.DecodeJWT(cookie.Value)
	currentBalanceModel, err := database.DBStorage.GetCurrentBalanceByLogin(login)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tFloat := 0.0
			currentBalanceModel = &model.CurrentBalanceModel{Balance: &tFloat, Login: &login}
		} else {
			helpers.TLog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	withdrawnBalanceSum, err := database.DBStorage.GetWithdrawnBalanceSumByLogin(login)
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	currentAndWithdrawnModel := model.CurrentAndWithdrawnModel{Current: currentBalanceModel.Balance, Withdrawn: withdrawnBalanceSum}
	marshal, err := json.Marshal(currentAndWithdrawnModel)
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshal)
}

func GetOrders(w http.ResponseWriter, req *http.Request) {
	cookie, _ := req.Cookie("Token")
	login, _ := helpers.DecodeJWT(cookie.Value)

	orders, err := database.DBStorage.GetOrdersByLogin(login)
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, "Ошибка сервера!", http.StatusInternalServerError)
		return

	}
	w.Header().Set("Content-Type", "application/json")
	if len(*orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	tOrders := *orders

	sort.Slice(tOrders, func(i, j int) bool {
		return tOrders[i].UploadedAT.After(*tOrders[j].UploadedAT)
	})

	marshal, err := json.Marshal(tOrders)
	if err != nil {
		http.Error(w, "Ошибка записи ответа", http.StatusNotFound)
		return
	}
	_, errWrite := w.Write(marshal)
	if errWrite != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
}

func SaveOrders(w http.ResponseWriter, req *http.Request) {
	cookie, _ := req.Cookie("Token")
	login, _ := helpers.DecodeJWT(cookie.Value)

	b, err := io.ReadAll(req.Body)
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, "Неверный формат номера заказа!", http.StatusUnprocessableEntity)
		return
	}
	orderID, err := strconv.Atoi(string(b))
	if !helpers.LuhnValid(orderID) || err != nil {
		http.Error(w, "Неверный формат номера заказа!", http.StatusUnprocessableEntity)
		return
	}
	orders, err := database.DBStorage.GetOrdersByOrderID(string(b))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tNumber := string(b)
			err := database.DBStorage.CreateOrders(model.OrdersModel{
				OrdersID: &tNumber,
				Login:    &login})
			if err != nil {
				helpers.TLog.Error(err.Error())
				http.Error(w, "Ошибка сервера!", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusAccepted)
			return
		} else {
			helpers.TLog.Error(err.Error())
			http.Error(w, "Ошибка сервера!", http.StatusInternalServerError)
			return
		}
	}
	if *orders.Login == login {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		http.Error(w, "Номер заказа уже был загружен другим пользователем!", http.StatusConflict)
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
