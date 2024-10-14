package engine

import (
	"orderbook/internal/models"
	"time"
)

func NewTrade(buyOrder, sellOrder *models.Order, price float64, quantity int) *models.Trade {
	trade := &models.Trade{
		BuyOrder:  buyOrder,
		SellOrder: sellOrder,
		Quantity:  quantity,
		Price:     price,
		Timestamp: time.Now(),
	}

	buyOrder.Quantity -= quantity
	sellOrder.Quantity -= quantity

	if buyOrder.Quantity == 0 {
		buyOrder.Status = "FILLED"
	}
	if sellOrder.Quantity == 0 {
		sellOrder.Status = "FILLED"
	}

	return trade
}

func ExecuteTrade(buyOrder, sellOrder *models.Order) *models.Trade {
	tradePrice := sellOrder.Price
	tradeQuantity := sellOrder.Quantity
	if sellOrder.Quantity > buyOrder.Quantity {
		tradeQuantity = buyOrder.Quantity
	}

	trade := NewTrade(buyOrder, sellOrder, tradePrice, tradeQuantity)

	return trade
}

func TradeLogger(trade *models.Trade) {
	println("Executed Trade:")
	println("Buy Order ID: ", trade.BuyOrder.ID)
	println("Sell Order ID: ", trade.SellOrder.ID)
	println("Price: ", trade.Price)
	println("Quantity: ", trade.Quantity)
	println("Timestamp: ", trade.Timestamp)
}
