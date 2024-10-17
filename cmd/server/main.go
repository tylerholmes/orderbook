package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"orderbook/internal/routes"
)

func main() {
	r := gin.Default()

	routes.Run(r)

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
