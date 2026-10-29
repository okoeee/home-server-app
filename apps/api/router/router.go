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

	v1 := router.Group("/v1")
	v1.GET("/daily-price-changes/latest", dailyPriceChangeHandler.GetLatestDailyPriceChanges)
	v1.GET("/daily-volume-ratios", dailyVolumeRatioHandler.GetDailyVolumeRatios)

	v1.GET(
		"/stocks/:symbol/prices",
		stockHandler.GetDailyPrices,
	)

	return router
}
