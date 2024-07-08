package model

import "time"

type OrdersModel struct {
	OrdersID   *int64     `json:"number"`
	Login      *string    `json:"-"`
	Accrual    *int64     `json:"accrual"`
	Status     *string    `json:"status"`
	UploadedAT *time.Time `json:"uploaded_at"`
}
