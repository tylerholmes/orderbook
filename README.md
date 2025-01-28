# Orderbook

A simple orderbook system implemented in Go. Features limit and market orders, price matching engine, and Postgres database for persistence.

## Overview

This project implements a basic trading platform with the following features:
- Limit and Market orders
- Real-time price matching engine
- Postgres database for persistence
- REST API for order placement and retrieval
- Real market data using Alpha Vantage
- Automatic trade execution

## Tech stack
- Go with Echo framework
- Postgres
- Docker

### Setup and running

1. Clone the repo

2. Create a `.env` file with the following content:
```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=orderbook
DB_PASSWORD=orderbook
DB_NAME=orderbook
ALPHA_VANTAGE_API_KEY=your_api_key_here
```
3. Run `docker-compose up` to start the app and the database

### API Endpoints
The API is at `http://localhost:8080`

#### Place an Order
```bash
curl -X POST http://localhost:8080/orders \
-H "Content-Type: application/json" \
-d '{
    "symbol": "NVDA",
    "side": "buy",
    "quantity": 10,
    "price": 118.58,
    "order_type": "limit"
}'
```

#### Get All Orders
```bash
curl http://localhost:8080/orders
```

#### Get Orderbook for Symbol
```bash
curl http://localhost:8080/orderbook/AAPL
```

#### Get Market Price
```bash
curl http://localhost:8080/market-price/AAPL
```

#### Cancel Order
```bash
curl -X DELETE http://localhost:8080/orders/{orderID}
```