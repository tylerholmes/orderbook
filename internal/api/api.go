package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"orderbook/internal/engine"
	"orderbook/internal/models"
)

func PlaceOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orderbook := engine.NewOrderbook()
	orderbook.AddOrder(&order)
	c.JSON(http.StatusOK, gin.H{"status": "Order placed"})
}

func CancelOrder(c *gin.Context) {
	orderID := c.Param("id")
	// Logic to cancel order
	c.JSON(http.StatusOK, gin.H{"status": "Order canceled", "orderID": orderID})
}

func GetOrderBook(c *gin.Context) {
	stockID := c.Param("stock_id")
	// Logic to get order book
	c.JSON(http.StatusOK, gin.H{"stock_id": stockID, "buy_orders": "", "sell_orders": ""})
}

func ExecuteTrade(c *gin.Context) {
	// Logic to execute trade
	c.JSON(http.StatusOK, gin.H{"status": "Trade executed"})
}
