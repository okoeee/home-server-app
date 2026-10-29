package handler

import (
	"api/openapi"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {

	c.JSON(http.StatusOK, openapi.PingResponse{
		Message: "pong",
	})

}
