package handler

import (
	"api/domain/model"
	"api/external/prometheus"
	"api/openapi"
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServerMetricsReader interface {
	ReadMetrics(context.Context) (model.ServerMetrics, error)
}

type ServerMetricsHandler struct {
	reader ServerMetricsReader
}

func NewServerMetricsHandler(reader ServerMetricsReader) *ServerMetricsHandler {
	return &ServerMetricsHandler{reader: reader}
}

func (handler *ServerMetricsHandler) GetMetrics(c *gin.Context) {
	metrics, err := handler.reader.ReadMetrics(c.Request.Context())
	if err != nil {
		log.Printf("get server metrics: %v", err)
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, prometheus.ErrQueryFailed):
			status = http.StatusBadGateway
		case errors.Is(err, prometheus.ErrMetricsUnavailable):
			status = http.StatusServiceUnavailable
		}
		c.JSON(status, openapi.ErrorResponse{Message: "サーバーのメトリクスを取得できませんでした"})
		return
	}
	c.JSON(http.StatusOK, openapi.ServerMetricsResponse{
		FetchedAt:          metrics.FetchedAt,
		CpuUsagePercent:    metrics.CPUUsagePercent,
		MemoryUsagePercent: metrics.MemoryUsagePercent,
	})
}
