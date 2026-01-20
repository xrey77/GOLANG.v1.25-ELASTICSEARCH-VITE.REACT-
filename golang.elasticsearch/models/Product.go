package models

import (
	"time"
)

type Product struct {
	ID             string    `json:"id"`
	Category       string    `json:"category"`
	Descriptions   string    `json:"descriptions"`
	Qty            float64   `json:"qty"`
	Unit           string    `json:"unit"`
	Costprice      float64   `json:"costprice"`
	Sellprice      float64   `json:"sellprice"`
	Saleprice      float64   `json:"saleprice"`
	Productpicture *string   `json:"productpicture"`
	Alertstocks    float64   `json:"alertstocks"`
	Criticalstocks float64   `json:"criticalstocks"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
