package model

import (
	"encoding/json"
	"fmt"
	"io"
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

func WithdrawnBalanceModelDecode(body io.ReadCloser) (*WithdrawnBalanceModel, error) {
	var tModel WithdrawnBalanceModel
	err := json.NewDecoder(body).Decode(&tModel)
	if err != nil {
		helpers.TLog.Error(err.Error())
		return nil, err
	}
	if tModel.Order == nil {
		err := fmt.Errorf("неверный формат запроса")
		return &tModel, err
	}
	orderID, err := strconv.Atoi(*tModel.Order)
	if !helpers.LuhnValid(orderID) || err != nil {
		err := fmt.Errorf("неверный формат запроса")
		return &tModel, err
	}
	return &tModel, nil
}
