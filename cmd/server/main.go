package main

import (
	"log"
	"orderbook/internal/handlers"
	database "orderbook/internal/pgdb"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.NewDatabase(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Initialize(); err != nil {
		log.Fatal(err)
	}

	api := handlers.NewOrderbookAPI(db, os.Getenv("ALPHA_VANTAGE_API_KEY"))

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	handlers.RegisterHandlers(e, api)

	e.Logger.Fatal(e.Start(":8080"))
}
