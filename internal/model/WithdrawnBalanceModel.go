package model

import "time"

type WithdrawnBalanceModel struct {
	Login       *string   `json:"-"`
	Order       *string   `json:"order"`
	ProcessedAT time.Time `json:"processed_at"`
	Withdrawn   *float64  `json:"sum"`
}
