package handler

import (
	"api/db"
	"api/openapi"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
			openapi.ErrorResponse{
				Message: "failed to get daily prices",
			})

		return
	}

	responses := make([]openapi.DailyPriceResponse, 0, len(prices))

	for _, price := range prices {
		response := openapi.DailyPriceResponse{
			Symbol:     price.Symbol,
			TradeDate:  price.TradeDate.Format("2006-01-02"),
			OpenPrice:  price.OpenPrice.StringFixed(1),
			HighPrice:  price.HighPrice.StringFixed(1),
			LowPrice:   price.LowPrice.StringFixed(1),
			ClosePrice: price.ClosePrice.StringFixed(1),
			Volume:     price.Volume.StringFixed(1),
		}

		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, responses)
}
