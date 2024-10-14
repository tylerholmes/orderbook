package engine

import (
	"container/heap"
	"sync"

	"orderbook/internal/models"
	"orderbook/internal/priorityqueue"
)

type Orderbook struct {
	BuyOrders  *priorityqueue.PriorityQueue
	SellOrders *priorityqueue.PriorityQueue
	mutex      sync.Mutex
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
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	if order.Type == "buy" {
		heap.Push(ob.BuyOrders, order)
	} else if order.Type == "sell" {
		heap.Push(ob.SellOrders, order)
	}

	ob.MatchOrder()
}

func (ob *Orderbook) CancelOrder(order *models.Order) {
	if order != nil {
		order.Status = "CANCELLED"
		order.Quantity = 0
	}
}

func (ob *Orderbook) MatchOrder() {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	for ob.BuyOrders.Len() > 0 && ob.SellOrders.Len() > 0 {
		topBuyOrder := ob.BuyOrders.Peek().(*models.Order)
		topSellOrder := ob.SellOrders.Peek().(*models.Order)

		if topBuyOrder.Price >= topSellOrder.Price {
			trade := ExecuteTrade(topBuyOrder, topSellOrder)
			TradeLogger(trade)

			if topBuyOrder.Quantity == 0 {
				ob.BuyOrders.Pop()
			}
			if topSellOrder.Quantity == 0 {
				ob.SellOrders.Pop()
			}
		}

	}

}
