package router

import (
	"api/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	stockHandler *handler.StockHandler,
) *gin.Engine {
	router := gin.Default()

	router.GET("/ping", handler.Ping)

	router.GET(
		"/stocks/:symbol/prices",
		stockHandler.GetDailyPrices,
	)

	return router
}
