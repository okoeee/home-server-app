package handler

import (
	"api/usecase/dailyprice"
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

type DailyPriceChangeHandler struct {
	calculator *dailyprice.DailyPriceChangeCalculator
}

func NewDailyPriceChangeHandler(calculator *dailyprice.DailyPriceChangeCalculator) *DailyPriceChangeHandler {
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

	c.JSON(http.StatusOK, responses)
}
