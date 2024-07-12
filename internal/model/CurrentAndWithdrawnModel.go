package model

type CurrentAndWithdrawnModel struct {
	Current   *float64 `json:"current"`
	Withdrawn *float64 `json:"withdrawn"`
}
