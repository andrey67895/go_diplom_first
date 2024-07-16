package model

import (
	"net/http"
	"time"
)

type OrdersModel struct {
	OrdersID   *string    `json:"number"`
	Login      *string    `json:"-"`
	Accrual    *float64   `json:"accrual"`
	Status     *string    `json:"status"`
	UploadedAT *time.Time `json:"uploaded_at"`
}

func (o OrdersModel) IsConflictByLogin(login string, w http.ResponseWriter) {
	if *o.Login == login {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Номер заказа уже был загружен другим пользователем!", http.StatusConflict)
	}
}
