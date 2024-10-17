package routes

import "github.com/gin-gonic/gin"

func Run(r *gin.Engine) {
	r.POST("/orders", api.PlaceOrder)
	r.DELETE("/orders/:id"), api.CancelOrder)
	r.GET("/orderbook/:stock_id", api.GetOrderBook)
	r.POST("/trade", api.ExecuteTrade)
}