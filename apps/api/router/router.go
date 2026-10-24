package router

import (
	"api/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	stockHandler *handler.StockHandler,
	dailyPriceChangeHandler *handler.DailyPriceChangeHandler,
) *gin.Engine {
	router := gin.Default()

	router.GET("/ping", handler.Ping)
	router.GET("/daily-price-changes", dailyPriceChangeHandler.GetDailyPriceChanges)

	router.GET(
		"/stocks/:symbol/prices",
		stockHandler.GetDailyPrices,
	)

	return router
}
