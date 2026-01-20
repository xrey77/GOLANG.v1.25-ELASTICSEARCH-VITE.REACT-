package models

import "time"

type Sale struct {
	ID        string    `json:"id"`
	Amount    float64   `json:"amount"`
	Salesdate time.Time `json:"salesdate"`
}
