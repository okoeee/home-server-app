package router

import (
	"api/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/ping", handler.Ping)

	return router
}
