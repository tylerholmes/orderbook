package models

import (
	"time"
)

type Order struct {
	ID          string     `db:"id" json:"id"`
	UserID      string     `db:"user_id" json:"user_id"`
	Symbol      string     `db:"symbol" json:"symbol"`
	Side        string     `db:"side" json:"side"`
	OrderType   string     `db:"order_type" json:"order_type"`
	Quantity    float64    `db:"quantity" json:"quantity"`
	FilledQty   float64    `db:"filled_qty" json:"filled_qty"`
	Price       float64    `db:"price" json:"price"`
	Status      string     `db:"status" json:"status"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	CancelledAt *time.Time `db:"cancelled_at" json:"cancelled_at,omitempty"`
	CompletedAt *time.Time `db:"completed_at" json:"completed_at,omitempty"`
}

type Trade struct {
	ID          string    `db:"id" json:"id"`
	BuyOrderID  string    `db:"buy_order_id" json:"buy_order_id"`
	SellOrderID string    `db:"sell_order_id" json:"sell_order_id"`
	Symbol      string    `db:"symbol" json:"symbol"`
	Quantity    float64   `db:"quantity" json:"quantity"`
	Price       float64   `db:"price" json:"price"`
	ExecutedAt  time.Time `db:"executed_at" json:"executed_at"`
}

type User struct {
	ID           string     `db:"id" json:"id"`
	Email        string     `db:"email" json:"email"`
	PasswordHash string     `db:"password_hash" json:"-"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
	LastLoginAt  *time.Time `db:"last_login_at" json:"last_login_at,omitempty"`
}

type Symbol struct {
	ID        string    `db:"id" json:"id"`
	Symbol    string    `db:"symbol" json:"symbol"`
	Name      string    `db:"name" json:"name"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type OrderHistory struct {
	ID        string    `db:"id" json:"id"`
	OrderID   string    `db:"order_id" json:"order_id"`
	Status    string    `db:"status" json:"status"`
	Notes     string    `db:"notes" json:"notes"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
