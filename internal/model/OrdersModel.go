package model

import (
	"fmt"
	"time"
)

type OrdersModel struct {
	OrdersID   *string    `json:"number"`
	Login      *string    `json:"-"`
	Accrual    *float64   `json:"accrual"`
	Status     *string    `json:"status"`
	UploadedAT *time.Time `json:"uploaded_at"`
}

func (o OrdersModel) IsConflictByLogin(login string) error {
	if *o.Login != login {
		return fmt.Errorf("номер заказа уже был загружен другим пользователем")
	}
	return nil
}
