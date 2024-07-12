package model

type CurrentAndWithdrawnModel struct {
	Current   *float64 `json:"current"`
	Withdrawn *string  `json:"withdrawn"`
}
