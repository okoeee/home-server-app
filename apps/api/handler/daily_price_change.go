package handler

import (
	"api/domain/model"
	"api/openapi"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DailyPriceChangeCalculator interface {
	CalculateLatest(context.Context) (model.LatestDailyPriceChanges, error)
}

type DailyPriceChangeHandler struct {
	calculator DailyPriceChangeCalculator
}

func NewDailyPriceChangeHandler(calculator DailyPriceChangeCalculator) *DailyPriceChangeHandler {
	return &DailyPriceChangeHandler{calculator: calculator}
}

func (handler *DailyPriceChangeHandler) GetLatestDailyPriceChanges(c *gin.Context) {
	result, err := handler.calculator.CalculateLatest(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{Message: "failed to calculate latest daily price changes"})
		return
	}
	var tradeDate *string
	if result.TradeDate != nil {
		formatted := result.TradeDate.Format(time.DateOnly)
		tradeDate = &formatted
	}
	c.JSON(http.StatusOK, openapi.LatestDailyPriceChangesResponse{
		TradeDate: tradeDate,
		Changes:   toDailyPriceChangeResponses(result.Changes),
	})
}

func toDailyPriceChangeResponses(changes []model.DailyPriceChange) []openapi.DailyPriceChangeResponse {
	responses := make([]openapi.DailyPriceChangeResponse, 0, len(changes))
	for _, change := range changes {
		responses = append(responses, openapi.DailyPriceChangeResponse{
			Symbol:             change.Symbol,
			CompanyName:        change.CompanyName,
			TradeDate:          change.TradeDate.Format(time.DateOnly),
			ClosePrice:         change.ClosePrice.StringFixed(1),
			PreviousTradeDate:  change.PreviousTradeDate.Format(time.DateOnly),
			PreviousClosePrice: change.PreviousClosePrice.StringFixed(1),
			ChangeRate:         change.ChangeRate.String(),
		})
	}

	return responses
}
