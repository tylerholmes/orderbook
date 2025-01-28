package handlers

import (
	"net/http"
	"orderbook/internal/core"
	"orderbook/internal/marketdata"
	"orderbook/internal/models"
	database "orderbook/internal/pgdb"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ServerInterface interface {
	PostOrders(ctx echo.Context) error
	GetOrders(ctx echo.Context) error
	GetOrderbookSymbol(ctx echo.Context, symbol string) error
	CancelOrder(ctx echo.Context, orderID string) error
	GetMarketPrice(ctx echo.Context, symbol string) error
}

type OrderRequest struct {
	Symbol    string  `json:"symbol"`
	Side      string  `json:"side"`
	Quantity  float64 `json:"quantity"`
	Price     float64 `json:"price"`
	OrderType string  `json:"order_type"`
	UserID    string  `json:"user_id,omitempty"`
}

type OrderbookAPI struct {
	db             *database.Database
	matchingEngine *core.MatchingEngine
	marketData     *marketdata.AlphaVantage
}

const (
	TestUserID = "00000000-0000-0000-0000-000000000000"
)

const (
	OrderTypeLimit  = "limit"
	OrderTypeMarket = "market"
)

func RegisterHandlers(e *echo.Echo, api *OrderbookAPI) {
	e.POST("/orders", api.PostOrders)
	e.GET("/orders", api.GetOrders)
	e.GET("/orderbook/:symbol", func(c echo.Context) error {
		return api.GetOrderbookSymbol(c, c.Param("symbol"))
	})
	e.DELETE("/orders/:orderID", func(c echo.Context) error {
		return api.CancelOrder(c, c.Param("orderID"))
	})
	e.GET("/market-price/:symbol", func(c echo.Context) error {
		return api.GetMarketPrice(c, c.Param("symbol"))
	})
}

func NewOrderbookAPI(db *database.Database, alphaVantageKey string) *OrderbookAPI {
	return &OrderbookAPI{
		db:             db,
		matchingEngine: core.NewMatchingEngine(),
		marketData:     marketdata.NewAlphaVantage(alphaVantageKey),
	}
}

func (api *OrderbookAPI) PostOrders(ctx echo.Context) error {
	var orderRequest OrderRequest
	if err := ctx.Bind(&orderRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if orderRequest.OrderType != OrderTypeLimit && orderRequest.OrderType != OrderTypeMarket {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid order_type: must be 'limit' or 'market'")
	}

	if orderRequest.Side != "buy" && orderRequest.Side != "sell" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid side: must be 'buy' or 'sell'")
	}

	if orderRequest.Quantity <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "quantity must be greater than 0")
	}

	if orderRequest.OrderType == OrderTypeLimit && orderRequest.Price <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "price must be greater than 0 for limit orders")
	}

	userID := orderRequest.UserID
	if userID == "" {
		userID = TestUserID
	}

	if _, err := uuid.Parse(userID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user_id format")
	}

	order := &core.Order{
		ID:        uuid.New().String(),
		Symbol:    orderRequest.Symbol,
		Side:      core.Side(orderRequest.Side),
		Quantity:  orderRequest.Quantity,
		Price:     orderRequest.Price,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	api.matchingEngine.ProcessOrder(order)

	dbOrder := &models.Order{
		ID:        order.ID,
		UserID:    userID,
		Symbol:    orderRequest.Symbol,
		Side:      orderRequest.Side,
		OrderType: orderRequest.OrderType,
		Quantity:  orderRequest.Quantity,
		FilledQty: 0,
		Price:     orderRequest.Price,
		Status:    order.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := api.db.SaveOrder(dbOrder); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, order)
}

func (api *OrderbookAPI) GetOrders(ctx echo.Context) error {
	orders, err := api.db.GetOrders()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, orders)
}

func (api *OrderbookAPI) GetOrderbookSymbol(ctx echo.Context, symbol string) error {
	ob := api.matchingEngine.GetOrderbook(symbol)
	if ob == nil {
		return echo.NewHTTPError(http.StatusNotFound, "orderbook not found")
	}

	response := struct {
		Symbol string            `json:"symbol"`
		Bids   []core.PriceLevel `json:"bids"`
		Asks   []core.PriceLevel `json:"asks"`
	}{
		Symbol: ob.Symbol,
		Bids:   ob.Bids,
		Asks:   ob.Asks,
	}

	return ctx.JSON(http.StatusOK, response)
}

func (api *OrderbookAPI) CancelOrder(ctx echo.Context, orderID string) error {
	if err := api.matchingEngine.CancelOrder(orderID); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	if err := api.db.UpdateOrderStatus(orderID, "cancelled"); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}

func (api *OrderbookAPI) GetMarketPrice(ctx echo.Context, symbol string) error {
	quote, err := api.marketData.GetQuote(symbol)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, quote)
}
