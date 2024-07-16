package model

import (
	"encoding/json"
	"io"
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

func (tModel WithdrawnBalanceModel) Marshal() []byte {
	marshal, _ := json.Marshal(tModel)
	return marshal
}

func WithdrawnBalanceModelDecode(w http.ResponseWriter, body io.ReadCloser) (*WithdrawnBalanceModel, error) {
	var tModel WithdrawnBalanceModel
	err := json.NewDecoder(body).Decode(&tModel)
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	if tModel.Order == nil {
		http.Error(w, "Неверный формат запроса!", http.StatusUnprocessableEntity)
		return nil, err
	}
	orderID, err := strconv.Atoi(*tModel.Order)
	if !helpers.LuhnValid(orderID) || err != nil {
		http.Error(w, "Неверный формат номера заказа!", http.StatusUnprocessableEntity)
		return nil, err
	}
	return &tModel, nil
}
