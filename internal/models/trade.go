package models

import "time"

type Trade struct {
	BuyOrder  *Order
	SellOrder *Order
	Quantity  int
	Price     float64
	Timestamp time.Time
}
