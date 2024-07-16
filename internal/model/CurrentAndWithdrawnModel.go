package model

import (
	"encoding/json"
	"net/http"

	"github.com/andrey67895/go_diplom_first/internal/helpers"
)

type CurrentAndWithdrawnModel struct {
	Current   *float64 `json:"current"`
	Withdrawn *float64 `json:"withdrawn"`
}

func CreateCurrentAndWithdrawnModelForMarshal(current *float64, withdrawn *float64, w http.ResponseWriter) []byte {
	currentAndWithdrawnModel := CurrentAndWithdrawnModel{Current: current, Withdrawn: withdrawn}
	marshal := currentAndWithdrawnModel.marshal(w)
	return marshal
}

func (c CurrentAndWithdrawnModel) marshal(w http.ResponseWriter) []byte {
	marshal, err := json.Marshal(c)
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return marshal
}
