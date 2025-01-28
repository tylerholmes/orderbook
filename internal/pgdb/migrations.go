package database

const schema = `
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP WITH TIME ZONE
);

-- Symbols table
CREATE TABLE IF NOT EXISTS symbols (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    symbol VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Orders table
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id), -- Made nullable for testing
    symbol VARCHAR(20) NOT NULL REFERENCES symbols(symbol),
    side VARCHAR(4) NOT NULL CHECK (side IN ('buy', 'sell')),
    order_type VARCHAR(10) NOT NULL CHECK (order_type IN ('market', 'limit')),
    quantity DECIMAL(18,8) NOT NULL CHECK (quantity > 0),
    filled_qty DECIMAL(18,8) NOT NULL DEFAULT 0 CHECK (filled_qty >= 0),
    price DECIMAL(18,8) CHECK (price > 0),
    status VARCHAR(20) NOT NULL CHECK (
        status IN ('pending', 'partial', 'filled', 'cancelled', 'rejected')
    ),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    cancelled_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT quantity_check CHECK (filled_qty <= quantity)
);

-- Trades table
CREATE TABLE IF NOT EXISTS trades (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    buy_order_id UUID NOT NULL REFERENCES orders(id),
    sell_order_id UUID NOT NULL REFERENCES orders(id),
    symbol VARCHAR(20) NOT NULL REFERENCES symbols(symbol),
    quantity DECIMAL(18,8) NOT NULL CHECK (quantity > 0),
    price DECIMAL(18,8) NOT NULL CHECK (price > 0),
    executed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Order History table
CREATE TABLE IF NOT EXISTS order_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id),
    status VARCHAR(20) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_symbol ON orders(symbol);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_trades_buy_order_id ON trades(buy_order_id);
CREATE INDEX IF NOT EXISTS idx_trades_sell_order_id ON trades(sell_order_id);
CREATE INDEX IF NOT EXISTS idx_trades_symbol ON trades(symbol);
CREATE INDEX IF NOT EXISTS idx_order_history_order_id ON order_history(order_id);

-- Triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_symbols_updated_at
    BEFORE UPDATE ON symbols
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert initial test data
DO $$ 
BEGIN
    -- Insert test symbols if they don't exist
    INSERT INTO symbols (symbol, name) 
    VALUES 
        ('AAPL', 'Apple Inc.'),
        ('GOOGL', 'Alphabet Inc.'),
        ('MSFT', 'Microsoft Corporation'),
        ('AMZN', 'Amazon.com, Inc.')
    ON CONFLICT (symbol) DO NOTHING;

    -- Insert test user if it doesn't exist
    INSERT INTO users (id, email, password_hash)
    VALUES (
        '00000000-0000-0000-0000-000000000000',
        'test@example.com',
        '$2a$10$test_hash_for_testing_purposes'
    )
    ON CONFLICT (email) DO NOTHING;
END $$;
`

// Additional migrations functions if needed
func (d *Database) Migrate() error {
	_, err := d.db.Exec(schema)
	if err != nil {
		return err
	}
	return nil
}

// Function to reset the database (useful for testing)
func (d *Database) ResetDatabase() error {
	resetSQL := `
    DROP TABLE IF EXISTS order_history CASCADE;
    DROP TABLE IF EXISTS trades CASCADE;
    DROP TABLE IF EXISTS orders CASCADE;
    DROP TABLE IF EXISTS symbols CASCADE;
    DROP TABLE IF EXISTS users CASCADE;
    `

	_, err := d.db.Exec(resetSQL)
	if err != nil {
		return err
	}

	// Re-run migrations
	return d.Migrate()
}
