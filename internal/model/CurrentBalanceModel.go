package model

import (
	"fmt"
)

type CurrentBalanceModel struct {
	Login   *string
	Balance *float64
}

func (c CurrentBalanceModel) IsValidByWithdrawn(withdrawn float64) error {
	if *c.Balance < withdrawn {
		return fmt.Errorf("на счету недостаточно средств")
	}
	return nil
}
