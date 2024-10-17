package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type Client struct {
	DB     *sql.DB
	Routes *gin.Engine
}

var (
	UNAMEDB string = "postgres"
	PASSDB  string = "password"
	HOSTDB  string = "username"
	DBNAME  string = "orderbook"
)

func (client *Client) CreateConnection() {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", UNAMEDB, PASSDB, HOSTDB, DBNAME)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open a db connections: %w", err)
	}
	client.DB = db
}

func (client *Client) CloseConnection() {
	client.DB.Close()
}

func (client *Client) Run() {
	client.Routes.Run(":8080")
}
