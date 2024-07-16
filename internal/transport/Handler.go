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

	withdrawnHistory := services.GetWithdrawnBalanceAndSortByLogin(login, w)
	if len(*withdrawnHistory) == 0 {
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
	tModel := model.WithdrawnBalanceModelDecode(w, req)
	currentBalanceModel := services.GetCurrentBalanceByLogin(login, w)
	currentBalanceModel.IsValidByWithdrawn(*tModel.Withdrawn, w)
	tWithdrawnBalanceModel := model.WithdrawnBalanceModel{Login: &login, Order: tModel.Order, Withdrawn: tModel.Withdrawn}
	services.WithdrawnBalanceByLogin(tWithdrawnBalanceModel, w)
	w.WriteHeader(http.StatusOK)
}

func GetBalance(w http.ResponseWriter, req *http.Request) {
	cookie, _ := req.Cookie("Token")
	login, _ := helpers.DecodeJWT(cookie.Value)
	currentBalanceModel := services.GetCurrentBalanceByLogin(login, w)
	withdrawnBalanceSum := services.GetWithdrawnBalanceSum(login, w)
	marshal := model.CreateCurrentAndWithdrawnModelForMarshal(currentBalanceModel.Balance, withdrawnBalanceSum, w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshal)
}

func GetOrders(w http.ResponseWriter, req *http.Request) {
	cookie, _ := req.Cookie("Token")
	login, _ := helpers.DecodeJWT(cookie.Value)
	orders := services.GetOrdersAndSortByLogin(login, w)
	if len(*orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	marshal, err := json.Marshal(orders)
	helpers.TLog.Info("NENENE ::: ", marshal)
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
	orderID := services.GetOrderIDAndValid(w, req)
	tModel := model.OrdersModel{OrdersID: &orderID, Login: &login}
	orders := services.GetOrderByOrderIDOrCreate(tModel, w)
	if orders != nil {
		orders.IsConflictByLogin(login, w)
	}
}

func UserRegister(w http.ResponseWriter, req *http.Request) {
	tModel := model.UserModelDecode(w, req)
	tModel.IsValid(w)
	services.GetAuth(tModel, w, true)
}

func AuthUser(w http.ResponseWriter, req *http.Request) {
	tModel := model.UserModelDecode(w, req)
	tModel.IsValid(w)
	auth := services.GetAuth(tModel, w, false)
	if *auth.Password == helpers.EncodeHash(*tModel.Password) {
		helpers.CreateAndSetJWTCookieInHTTP(*tModel.Login, w)
	}
}
