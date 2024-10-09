package models

import "time"

type Order struct {
	ID        int
	Symbol    string
	Type      string
	Price     float64
	Quantity  int
	Timestamp time.Time
}
