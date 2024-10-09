package main

import (
	"net/http"
	"orderbook/internal/api"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	api.RegisterRoutes(r)

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
