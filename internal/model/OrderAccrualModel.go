package model

type OrdersAccrualModel struct {
	OrderId *string  `json:"order"`
	Accrual *float64 `json:"accrual"`
	Status  *string  `json:"status"`
}
