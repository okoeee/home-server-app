package handler

import (
	"api/db"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type dailyPriceResponse struct {
	Symbol     string `json:"symbol"`
	TradeDate  string `json:"trade_date"`
	OpenPrice  string `json:"open_price"`
	HighPrice  string `json:"high_price"`
	LowPrice   string `json:"low_price"`
	ClosePrice string `json:"close_price"`
	Volume     string `json:"volume"`
}

type StockQuerier interface {
	ListDailyPricesBySymbol(
		context.Context,
		string,
	) ([]db.DailyPrice, error)
}

type StockHandler struct {
	queries StockQuerier
}

func NewStockHandler(queries StockQuerier) *StockHandler {
	return &StockHandler{
		queries: queries,
	}
}

func (handler *StockHandler) GetDailyPrices(c *gin.Context) {
	symbol := c.Param("symbol")

	prices, err := handler.queries.ListDailyPricesBySymbol(
		c.Request.Context(),
		symbol,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "failed to get daily prices",
			})

		return
	}

	responses := make([]dailyPriceResponse, 0, len(prices))

	for _, price := range prices {
		response := dailyPriceResponse{
			Symbol:     price.Symbol,
			TradeDate:  price.TradeDate.Format("2006-01-02"),
			OpenPrice:  price.OpenPrice,
			HighPrice:  price.HighPrice,
			LowPrice:   price.LowPrice,
			ClosePrice: price.ClosePrice,
			Volume:     price.Volume,
		}

		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, responses)
}
