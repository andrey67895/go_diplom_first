package model

import "net/http"

type CurrentBalanceModel struct {
	Login   *string
	Balance *float64
}

func (c CurrentBalanceModel) IsValidByWithdrawn(withdrawn float64, w http.ResponseWriter) {
	if *c.Balance < withdrawn {
		http.Error(w, "На счету недостаточно средств", http.StatusPaymentRequired)
	}
}
