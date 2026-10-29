package handler

import (
	"api/domain/model"
	"api/openapi"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DailyVolumeRatioCalculator interface {
	Calculate(context.Context, time.Time) ([]model.DailyVolumeRatio, error)
}

type DailyVolumeRatioHandler struct {
	calculator DailyVolumeRatioCalculator
}

func NewDailyVolumeRatioHandler(calculator DailyVolumeRatioCalculator) *DailyVolumeRatioHandler {
	return &DailyVolumeRatioHandler{calculator: calculator}
}

func (handler *DailyVolumeRatioHandler) GetDailyVolumeRatios(c *gin.Context) {
	tradeDate, err := time.Parse(time.DateOnly, c.Query("date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.ErrorResponse{Message: "date must be in YYYY-MM-DD format"})
		return
	}

	ratios, err := handler.calculator.Calculate(c.Request.Context(), tradeDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{Message: "failed to calculate daily volume ratios"})
		return
	}

	responses := make([]openapi.DailyVolumeRatioResponse, 0, len(ratios))
	for _, ratio := range ratios {
		responses = append(responses, openapi.DailyVolumeRatioResponse{
			Symbol:    ratio.Symbol,
			TradeDate: ratio.TradeDate.Format(time.DateOnly),
			Volume:    ratio.Volume.StringFixed(1),
			Ratio:     ratio.Ratio.String(),
		})
	}
	c.JSON(http.StatusOK, responses)
}
