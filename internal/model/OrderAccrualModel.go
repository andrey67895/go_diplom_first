package model

type OrdersAccrualModel struct {
	OrderId *string `json:"order"`
	Accrual *int64  `json:"accrual"`
	Status  *string `json:"status"`
}
