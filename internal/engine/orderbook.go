package engine

import (
	"container/heap"
	"orderbook/internal/models"
	"orderbook/internal/priorityqueue"
)

type Orderbook struct {
	BuyOrders  *priorityqueue.PriorityQueue
	SellOrders *priorityqueue.PriorityQueue
}

func NewOrderbook() *Orderbook {
	buyPQ := priorityqueue.InitPriorityQueue()
	sellPQ := priorityqueue.InitPriorityQueue()

	return &Orderbook{
		BuyOrders:  buyPQ,
		SellOrders: sellPQ,
	}
}

func (ob *Orderbook) AddOrder(order *models.Order) {
	if order.Type == "buy" {
		heap.Push(ob.BuyOrders, order)
	} else if order.Type == "sell" {
		heap.Push(ob.SellOrders, order)
	}
}

//TODO add the logic for removing/updating orders in order book
