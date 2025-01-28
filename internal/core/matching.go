package core

import (
	"fmt"
	"sync"
)

type MatchingEngine struct {
	orderbooks map[string]*Orderbook
	mu         sync.RWMutex
}

func NewMatchingEngine() *MatchingEngine {
	return &MatchingEngine{
		orderbooks: make(map[string]*Orderbook),
	}
}

func (m *MatchingEngine) ProcessOrder(order *Order) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ob, exists := m.orderbooks[order.Symbol]
	if !exists {
		ob = NewOrderbook(order.Symbol)
		m.orderbooks[order.Symbol] = ob
	}

	if order.Side == Buy {
		m.matchBuyOrder(ob, order)
	} else {
		m.matchSellOrder(ob, order)
	}

	if order.Quantity > 0 {
		order.Status = "pending"
		ob.AddOrder(order)
	}
}

func (m *MatchingEngine) CancelOrder(orderID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, ob := range m.orderbooks {
		if ob.cancelOrder(orderID, true) || ob.cancelOrder(orderID, false) {
			return nil
		}
	}

	return fmt.Errorf("order not found")
}

func (m *MatchingEngine) GetOrderbook(symbol string) *Orderbook {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.orderbooks[symbol]
}

func (m *MatchingEngine) matchBuyOrder(ob *Orderbook, order *Order) {
	for i := 0; i < len(ob.Asks); i++ {
		if ob.Asks[i].Price > order.Price {
			break
		}

		if order.Quantity == 0 {
			break
		}

		level := &ob.Asks[i]
		matchQuantity := min(order.Quantity, level.Quantity)

		order.Quantity -= matchQuantity
		level.Quantity -= matchQuantity

		if level.Quantity == 0 {
			ob.Asks = append(ob.Asks[:i], ob.Asks[i+1:]...)
			i--
		}
	}
}

func (m *MatchingEngine) matchSellOrder(ob *Orderbook, order *Order) {
	for i := 0; i < len(ob.Bids); i++ {
		if ob.Bids[i].Price < order.Price {
			break
		}

		if order.Quantity == 0 {
			break
		}

		level := &ob.Bids[i]
		matchQuantity := min(order.Quantity, level.Quantity)

		order.Quantity -= matchQuantity
		level.Quantity -= matchQuantity

		if level.Quantity == 0 {
			ob.Bids = append(ob.Bids[:i], ob.Bids[i+1:]...)
			i--
		}
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
