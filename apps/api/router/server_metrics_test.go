package router

import (
	"api/domain/model"
	"api/external/prometheus"
	"api/handler"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type serverMetricsStub struct {
	metrics model.ServerMetrics
	err     error
}

func (stub serverMetricsStub) ReadMetrics(context.Context) (model.ServerMetrics, error) {
	return stub.metrics, stub.err
}

func TestGetMetrics_ReturnsMetrics(t *testing.T) {
	// Arrange
	engine := NewRouter(
		handler.NewStockHandler(stockQuerierStub{}),
		handler.NewDailyPriceChangeHandler(priceChangeCalculatorStub{}),
		handler.NewDailyVolumeRatioHandler(volumeRatioCalculatorStub{}),
		handler.NewWeatherHandler(weatherForecastStub{}),
		handler.NewServerMetricsHandler(serverMetricsStub{metrics: model.ServerMetrics{
			FetchedAt:          time.Date(2026, 11, 7, 0, 0, 0, 0, time.UTC),
			CPUUsagePercent:    12.5,
			MemoryUsagePercent: 75,
		}}),
	)
	recorder := httptest.NewRecorder()

	// Act
	engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/v1/server/metrics", nil))

	// Assert
	assertJSONResponse(t, recorder, http.StatusOK, `{"fetched_at":"2026-11-07T00:00:00Z","cpu_usage_percent":12.5,"memory_usage_percent":75}`)
}

func TestGetMetrics_Errors(t *testing.T) {
	for _, test := range []struct {
		name   string
		err    error
		status int
	}{
		{"UpstreamError", fmt.Errorf("private details: %w", prometheus.ErrQueryFailed), http.StatusBadGateway},
		{"DataUnavailable", fmt.Errorf("private details: %w", prometheus.ErrMetricsUnavailable), http.StatusServiceUnavailable},
		{"UnexpectedError", errors.New("private internal details"), http.StatusInternalServerError},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			engine := NewRouter(
				handler.NewStockHandler(stockQuerierStub{}),
				handler.NewDailyPriceChangeHandler(priceChangeCalculatorStub{}),
				handler.NewDailyVolumeRatioHandler(volumeRatioCalculatorStub{}),
				handler.NewWeatherHandler(weatherForecastStub{}),
				handler.NewServerMetricsHandler(serverMetricsStub{err: test.err}),
			)
			recorder := httptest.NewRecorder()

			// Act
			engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/v1/server/metrics", nil))

			// Assert
			assertJSONResponse(t, recorder, test.status, `{"message":"サーバーのメトリクスを取得できませんでした"}`)
		})
	}
}
