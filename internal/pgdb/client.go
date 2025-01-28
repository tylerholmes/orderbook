package database

import (
	"database/sql"
	"errors"
	"fmt"
	"orderbook/internal/core"
	"orderbook/internal/models"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(host, port, user, password, dbname string) (*Database, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Initialize() error {
	_, err := d.db.Exec(schema)
	return err
}

// in internal/pgdb/client.go
func (d *Database) SaveOrder(order *models.Order) error {
	query := `
        INSERT INTO orders (
            id, user_id, symbol, side, order_type, 
            quantity, price, status, created_at, updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $9)
    `

	_, err := d.db.Exec(
		query,
		order.ID,
		order.UserID,
		order.Symbol,
		order.Side,
		order.OrderType,
		order.Quantity,
		order.Price,
		order.Status,
		time.Now(),
	)

	return err
}

func (d *Database) GetOrders() ([]core.Order, error) {
	query := `SELECT id, symbol, side, quantity, price, status, created_at FROM orders`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []core.Order
	for rows.Next() {
		var order core.Order
		err := rows.Scan(
			&order.ID,
			&order.Symbol,
			&order.Side,
			&order.Quantity,
			&order.Price,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (d *Database) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $4)
		RETURNING id
	`
	return d.db.QueryRow(
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		time.Now(),
	).Scan(&user.ID)
}

func (d *Database) CreateOrder(order *models.Order) error {
	query := `
		INSERT INTO orders (
			id, user_id, symbol, side, order_type, quantity,
			price, status, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $9)
		RETURNING id
	`
	return d.db.QueryRow(
		query,
		order.ID,
		order.UserID,
		order.Symbol,
		order.Side,
		order.OrderType,
		order.Quantity,
		order.Price,
		order.Status,
		time.Now(),
	).Scan(&order.ID)
}

func (d *Database) UpdateOrderStatus(orderID string, status string) error {
	query := `
		UPDATE orders
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	result, err := d.db.Exec(query, status, orderID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("order not found")
	}

	return nil
}

func (d *Database) CreateTrade(trade *models.Trade) error {
	query := `
		INSERT INTO trades (
			id, buy_order_id, sell_order_id, symbol,
			quantity, price, executed_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	return d.db.QueryRow(
		query,
		trade.ID,
		trade.BuyOrderID,
		trade.SellOrderID,
		trade.Symbol,
		trade.Quantity,
		trade.Price,
		time.Now(),
	).Scan(&trade.ID)
}

func (d *Database) CreateSymbol(symbol *models.Symbol) error {
	query := `
		INSERT INTO symbols (id, symbol, name, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING id
	`
	return d.db.QueryRow(
		query,
		symbol.ID,
		symbol.Symbol,
		symbol.Name,
		symbol.IsActive,
		time.Now(),
	).Scan(&symbol.ID)
}

func (d *Database) AddOrderHistory(history *models.OrderHistory) error {
	query := `
		INSERT INTO order_history (id, order_id, status, notes, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return d.db.QueryRow(
		query,
		history.ID,
		history.OrderID,
		history.Status,
		history.Notes,
		time.Now(),
	).Scan(&history.ID)
}

func (d *Database) GetOrdersByUserID(userID string) ([]models.Order, error) {
	query := `
		SELECT id, user_id, symbol, side, order_type, quantity,
			   filled_qty, price, status, created_at, updated_at,
			   cancelled_at, completed_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := d.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Symbol,
			&order.Side,
			&order.OrderType,
			&order.Quantity,
			&order.FilledQty,
			&order.Price,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.CancelledAt,
			&order.CompletedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
