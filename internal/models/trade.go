package models

type Trade struct {
	BuyOrder  *Order
	SellOrder *Order
	Quantity  int
	Price     float64
}
