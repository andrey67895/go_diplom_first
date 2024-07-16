package model

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/andrey67895/go_diplom_first/internal/helpers"
)

type WithdrawnBalanceModel struct {
	Login       *string    `json:"-"`
	Order       *string    `json:"order"`
	ProcessedAT *time.Time `json:"processed_at"`
	Withdrawn   *float64   `json:"sum"`
}

func WithdrawnBalanceModelDecode(w http.ResponseWriter, r *http.Request) WithdrawnBalanceModel {
	var tModel WithdrawnBalanceModel
	err := json.NewDecoder(r.Body).Decode(&tModel)
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	orderID, err := strconv.Atoi(*tModel.Order)
	if !helpers.LuhnValid(orderID) || err != nil {
		http.Error(w, "Неверный формат номера заказа!", http.StatusUnprocessableEntity)
	}
	return tModel
}
