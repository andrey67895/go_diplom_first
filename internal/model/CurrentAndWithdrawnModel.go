package model

import (
	"encoding/json"

	"github.com/andrey67895/go_diplom_first/internal/helpers"
)

type CurrentAndWithdrawnModel struct {
	Current   *float64 `json:"current"`
	Withdrawn *float64 `json:"withdrawn"`
}

func CreateCurrentAndWithdrawnModelForMarshal(current *float64, withdrawn *float64) ([]byte, error) {
	currentAndWithdrawnModel := CurrentAndWithdrawnModel{Current: current, Withdrawn: withdrawn}
	marshal, err := currentAndWithdrawnModel.marshal()
	return marshal, err
}

func (c CurrentAndWithdrawnModel) marshal() ([]byte, error) {
	marshal, err := json.Marshal(c)
	if err != nil {
		helpers.TLog.Error(err.Error())
		return nil, err
	}
	return marshal, nil
}
