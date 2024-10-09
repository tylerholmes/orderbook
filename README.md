# Order Book

## Overview
The **Order Book Project** is a stock trading application designed to simulate a real-time order book system, allowing users to place buy and sell orders for stocks. Built with Go and the Gin framework for the backend, React for the frontend, PostgreSQL for the database,

## Project Structure
```
├── cmd
│   └── server
│       └── main.go  # Entrypoint
├── internal
│   └── api            # api route handlers 
│   └── db             # database initialization and functions
│   └── engine         # core logic for the order book
│   └── models         # database models
│   └── priorityqueue  # priority queue implementation
├── pkg
│   └── middleware     # auth 
│   └── utils          # logging and error utils
```

## Project Roadmap
- [ ] core logic for orders
- [ ] create all api endpoints and route handling
- [ ] setup pgdb
- [ ] create db tables
- [ ] concurrency handling 
- [ ] setup user authentication
- [ ] create react frontend
- [ ] add docker to project
- [ ] deploy to gcp/github pages
- [ ] fetch real-time market data
- [ ] test with simulated market activity
- [ ] implement more advanced order matching algorithms (pro-rata, weighted, etc.)
- [ ] add support for different order types (stop loss, fill or kill, etc.)
- [ ] implement redis to optimize latency