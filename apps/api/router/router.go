package router

import (
	"api/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	stockHandler *handler.StockHandler,
	dailyPriceChangeHandler *handler.DailyPriceChangeHandler,
	dailyVolumeRatioHandler *handler.DailyVolumeRatioHandler,
) *gin.Engine {
	router := gin.Default()

	router.GET("/ping", handler.Ping)
	router.GET("/daily-price-changes", dailyPriceChangeHandler.GetDailyPriceChanges)
	router.GET("/daily-price-changes/latest", dailyPriceChangeHandler.GetLatestDailyPriceChanges)
	router.GET("/daily-volume-ratios", dailyVolumeRatioHandler.GetDailyVolumeRatios)

	router.GET(
		"/stocks/:symbol/prices",
		stockHandler.GetDailyPrices,
	)

	return router
}
