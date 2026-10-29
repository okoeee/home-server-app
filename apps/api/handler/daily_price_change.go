package handler

import (
	"api/domain/model"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type dailyPriceChangeResponse struct {
	Symbol             string `json:"symbol"`
	TradeDate          string `json:"trade_date"`
	ClosePrice         string `json:"close_price"`
	PreviousTradeDate  string `json:"previous_trade_date"`
	PreviousClosePrice string `json:"previous_close_price"`
	ChangeRate         string `json:"change_rate"`
}

type latestDailyPriceChangesResponse struct {
	TradeDate *string                    `json:"trade_date"`
	Changes   []dailyPriceChangeResponse `json:"changes"`
}

type DailyPriceChangeCalculator interface {
	Calculate(context.Context, time.Time) ([]model.DailyPriceChange, error)
	CalculateLatest(context.Context) (model.LatestDailyPriceChanges, error)
}

type DailyPriceChangeHandler struct {
	calculator DailyPriceChangeCalculator
}

func NewDailyPriceChangeHandler(calculator DailyPriceChangeCalculator) *DailyPriceChangeHandler {
	return &DailyPriceChangeHandler{calculator: calculator}
}

func (handler *DailyPriceChangeHandler) GetDailyPriceChanges(c *gin.Context) {
	tradeDate, err := time.Parse(time.DateOnly, c.Query("date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "date must be in YYYY-MM-DD format"})
		return
	}

	changes, err := handler.calculator.Calculate(c.Request.Context(), tradeDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to calculate daily price changes"})
		return
	}

	c.JSON(http.StatusOK, toDailyPriceChangeResponses(changes))
}

func (handler *DailyPriceChangeHandler) GetLatestDailyPriceChanges(c *gin.Context) {
	result, err := handler.calculator.CalculateLatest(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to calculate latest daily price changes"})
		return
	}
	var tradeDate *string
	if result.TradeDate != nil {
		formatted := result.TradeDate.Format(time.DateOnly)
		tradeDate = &formatted
	}
	c.JSON(http.StatusOK, latestDailyPriceChangesResponse{
		TradeDate: tradeDate,
		Changes:   toDailyPriceChangeResponses(result.Changes),
	})
}

func toDailyPriceChangeResponses(changes []model.DailyPriceChange) []dailyPriceChangeResponse {
	responses := make([]dailyPriceChangeResponse, 0, len(changes))
	for _, change := range changes {
		responses = append(responses, dailyPriceChangeResponse{
			Symbol:             change.Symbol,
			TradeDate:          change.TradeDate.Format(time.DateOnly),
			ClosePrice:         change.ClosePrice.StringFixed(1),
			PreviousTradeDate:  change.PreviousTradeDate.Format(time.DateOnly),
			PreviousClosePrice: change.PreviousClosePrice.StringFixed(1),
			ChangeRate:         change.ChangeRate.String(),
		})
	}

	return responses
}
