package transport

import (
	"encoding/json"
	"net/http"

	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/model"
	"github.com/andrey67895/go_diplom_first/internal/services"
)

func GetWithdrawalsHistory(w http.ResponseWriter, req *http.Request) {
	cookie, _ := req.Cookie("Token")
	login, _ := helpers.DecodeJWT(cookie.Value)

	withdrawnHistory, err := services.GetWithdrawnBalanceAndSortByLogin(login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(withdrawnHistory) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	marshal, err := json.Marshal(withdrawnHistory)
	if err != nil {
		http.Error(w, "Ошибка записи ответа", http.StatusNotFound)
		return
	}
	w.Write(marshal)
	w.WriteHeader(http.StatusOK)
}

func WithdrawBalance(w http.ResponseWriter, req *http.Request) {
	cookie, _ := req.Cookie("Token")
	login, _ := helpers.DecodeJWT(cookie.Value)
	tModel, err := model.WithdrawnBalanceModelDecode(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	currentBalanceModel, err := services.GetCurrentBalanceByLogin(login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := currentBalanceModel.IsValidByWithdrawn(*tModel.Withdrawn); err != nil {
		http.Error(w, err.Error(), http.StatusPaymentRequired)
		return
	}
	tWithdrawnBalanceModel := model.WithdrawnBalanceModel{Login: &login, Order: tModel.Order, Withdrawn: tModel.Withdrawn}
	if err := services.WithdrawnBalanceByLogin(tWithdrawnBalanceModel); err != nil {
		http.Error(w, err.Error(), http.StatusPaymentRequired)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetBalance(w http.ResponseWriter, req *http.Request) {
	cookie, _ := req.Cookie("Token")
	login, _ := helpers.DecodeJWT(cookie.Value)
	currentBalanceModel, err := services.GetCurrentBalanceByLogin(login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	withdrawnBalanceSum, err := services.GetWithdrawnBalanceSum(login)

	marshal, err := model.CreateCurrentAndWithdrawnModelForMarshal(currentBalanceModel.Balance, withdrawnBalanceSum)
	if err != nil {
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
	orders, err := services.GetOrdersAndSortByLogin(login)
	if err != nil {
		http.Error(w, "Ошибка сервера!", http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	marshal, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, "Ошибка записи ответа", http.StatusInternalServerError)
		return
	}
	w.Write(marshal)
	w.WriteHeader(http.StatusOK)
}

func SaveOrders(w http.ResponseWriter, req *http.Request) {
	cookie, _ := req.Cookie("Token")
	login, _ := helpers.DecodeJWT(cookie.Value)
	orderID, err := services.GetOrderIDAndValid(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if orderID != nil {
		tModel := model.OrdersModel{OrdersID: orderID, Login: &login}
		orders, err := services.GetOrderByOrderIDOrCreate(tModel)
		if err == nil {
			http.Error(w, "Ошибка сервера!", http.StatusInternalServerError)
			return
		} else {
			if orders == nil {
				w.WriteHeader(http.StatusAccepted)
				return
			} else {
				if err := orders.IsConflictByLogin(login); err != nil {
					http.Error(w, err.Error(), http.StatusConflict)
					return
				} else {
					w.WriteHeader(http.StatusOK)
				}
			}

		}
	}
}

func UserRegister(w http.ResponseWriter, req *http.Request) {
	tModel, err := model.UserModelDecode(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := tModel.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, apierr := services.GetAuth(*tModel, true)
	if apierr != nil {
		http.Error(w, apierr.Error.Error(), apierr.Status)
		return
	} else {
		token, err := helpers.CreateTokenInHTTP(*tModel.Login)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		helpers.SetCookie(token, w)
		w.WriteHeader(http.StatusOK)
	}

}

func AuthUser(w http.ResponseWriter, req *http.Request) {
	tModel, err := model.UserModelDecode(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := tModel.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	auth, apierr := services.GetAuth(*tModel, false)

	if apierr != nil {
		http.Error(w, apierr.Error.Error(), apierr.Status)
		return
	} else {
		if *auth.Password == helpers.EncodeHash(*tModel.Password) {
			token, err := helpers.CreateTokenInHTTP(*tModel.Login)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			helpers.SetCookie(token, w)
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
			return
		}

	}

}
