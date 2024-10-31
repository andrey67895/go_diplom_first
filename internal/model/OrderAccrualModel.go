package model

type OrdersAccrualModel struct {
	OrderID *string  `json:"order"`
	Accrual *float64 `json:"accrual"`
	Status  *string  `json:"status"`
}
