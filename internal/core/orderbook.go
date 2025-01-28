package core

import (
	"sync"
	"time"
)

type Side string

const (
	Buy  Side = "buy"
	Sell Side = "sell"
)

type Order struct {
	ID        string    `json:"id"`
	Symbol    string    `json:"symbol"`
	Side      Side      `json:"side"`
	Quantity  float64   `json:"quantity"`
	Price     float64   `json:"price"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type PriceLevel struct {
	Price    float64  `json:"price"`
	Quantity float64  `json:"quantity"`
	Orders   []*Order `json:"orders"`
}

type Orderbook struct {
	Symbol string
	Bids   []PriceLevel
	Asks   []PriceLevel
	mu     sync.RWMutex
}

func NewOrderbook(symbol string) *Orderbook {
	return &Orderbook{
		Symbol: symbol,
		Bids:   make([]PriceLevel, 0),
		Asks:   make([]PriceLevel, 0),
	}
}

func (ob *Orderbook) AddOrder(order *Order) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	if order.Side == Buy {
		ob.addBid(order)
	} else {
		ob.addAsk(order)
	}
}

func (ob *Orderbook) cancelOrder(orderID string, isBuy bool) bool {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	levels := &ob.Asks
	if isBuy {
		levels = &ob.Bids
	}

	for i := range *levels {
		level := &(*levels)[i]
		for j, order := range level.Orders {
			if order.ID == orderID {
				level.Quantity -= order.Quantity
				level.Orders = append(level.Orders[:j], level.Orders[j+1:]...)

				if level.Quantity == 0 {
					*levels = append((*levels)[:i], (*levels)[i+1:]...)
				}
				return true
			}
		}
	}

	return false
}

func (ob *Orderbook) addBid(order *Order) {
	for i, level := range ob.Bids {
		if level.Price == order.Price {
			ob.Bids[i].Quantity += order.Quantity
			ob.Bids[i].Orders = append(ob.Bids[i].Orders, order)
			return
		}
		if level.Price < order.Price {
			newLevel := PriceLevel{
				Price:    order.Price,
				Quantity: order.Quantity,
				Orders:   []*Order{order},
			}
			ob.Bids = append(ob.Bids, PriceLevel{})
			copy(ob.Bids[i+1:], ob.Bids[i:])
			ob.Bids[i] = newLevel
			return
		}
	}

	ob.Bids = append(ob.Bids, PriceLevel{
		Price:    order.Price,
		Quantity: order.Quantity,
		Orders:   []*Order{order},
	})
}

func (ob *Orderbook) addAsk(order *Order) {
	for i, level := range ob.Asks {
		if level.Price == order.Price {
			ob.Asks[i].Quantity += order.Quantity
			ob.Asks[i].Orders = append(ob.Asks[i].Orders, order)
			return
		}
		if level.Price > order.Price {
			newLevel := PriceLevel{
				Price:    order.Price,
				Quantity: order.Quantity,
				Orders:   []*Order{order},
			}
			ob.Asks = append(ob.Asks, PriceLevel{})
			copy(ob.Asks[i+1:], ob.Asks[i:])
			ob.Asks[i] = newLevel
			return
		}
	}

	ob.Asks = append(ob.Asks, PriceLevel{
		Price:    order.Price,
		Quantity: order.Quantity,
		Orders:   []*Order{order},
	})
}

func (ob *Orderbook) GetBids() []PriceLevel {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	bids := make([]PriceLevel, len(ob.Bids))
	copy(bids, ob.Bids)
	return bids
}

func (ob *Orderbook) GetAsks() []PriceLevel {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	asks := make([]PriceLevel, len(ob.Asks))
	copy(asks, ob.Asks)
	return asks
}
