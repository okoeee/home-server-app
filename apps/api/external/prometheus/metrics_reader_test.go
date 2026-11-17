package prometheus

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestReadMetrics_Results(t *testing.T) {
	at := time.Date(2026, 11, 7, 0, 0, 0, 0, time.UTC)
	for _, test := range []struct {
		name       string
		overrides  map[string][]float64
		wantErr    error
		wantCPU    float64
		wantMemory float64
	}{
		{name: "Normal", wantCPU: 12.5, wantMemory: 75},
		{name: "ZeroUsage", overrides: map[string][]float64{cpuUsageQuery: {0}, memoryAvailableQuery: {800}}},
		{name: "FullUsage", overrides: map[string][]float64{cpuUsageQuery: {100}, memoryAvailableQuery: {0}}, wantCPU: 100, wantMemory: 100},
		{name: "NodeMissing", overrides: map[string][]float64{upQuery: {}}, wantErr: ErrMetricsUnavailable},
		{name: "NodeDown", overrides: map[string][]float64{upQuery: {0}}, wantErr: ErrMetricsUnavailable},
		{name: "CPUInsufficientSamples", overrides: map[string][]float64{cpuUsageQuery: {}}, wantErr: ErrMetricsUnavailable},
		{name: "MemoryMissing", overrides: map[string][]float64{memoryAvailableQuery: {}}, wantErr: ErrMetricsUnavailable},
		{name: "InvalidUp", overrides: map[string][]float64{upQuery: {2}}, wantErr: ErrQueryFailed},
		{name: "CPUOutOfRange", overrides: map[string][]float64{cpuUsageQuery: {101}}, wantErr: ErrQueryFailed},
		{name: "MemoryTotalZero", overrides: map[string][]float64{memoryTotalQuery: {0}}, wantErr: ErrQueryFailed},
		{name: "MemoryAvailableNegative", overrides: map[string][]float64{memoryAvailableQuery: {-1}}, wantErr: ErrQueryFailed},
		{name: "MemoryAvailableExceedsTotal", overrides: map[string][]float64{memoryAvailableQuery: {801}}, wantErr: ErrQueryFailed},
		{name: "MultipleRequiredSeries", overrides: map[string][]float64{memoryTotalQuery: {800, 900}}, wantErr: ErrQueryFailed},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			values := map[string][]float64{
				upQuery: {1}, cpuUsageQuery: {12.5}, memoryTotalQuery: {800}, memoryAvailableQuery: {200},
			}
			for expression, override := range test.overrides {
				values[expression] = override
			}
			client := NewClient("http://prometheus:9090")
			client.httpClient.Transport = roundTripFunc(func(request *http.Request) (*http.Response, error) {
				if request.URL.Query().Get("time") != at.Format(time.RFC3339Nano) {
					t.Errorf("evaluation time = %q, want %s", request.URL.Query().Get("time"), at)
				}
				if _, ok := request.Context().Deadline(); !ok {
					t.Error("request has no deadline")
				}
				expression := request.URL.Query().Get("query")
				queryValues, ok := values[expression]
				if !ok {
					t.Fatalf("unexpected query: %s", expression)
				}
				result := make([]map[string]any, 0, len(queryValues))
				for _, value := range queryValues {
					result = append(result, map[string]any{"metric": map[string]string{}, "value": []any{at.Unix(), strconv.FormatFloat(value, 'f', -1, 64)}})
				}
				body, err := json.Marshal(map[string]any{"status": "success", "data": map[string]any{"resultType": "vector", "result": result}})
				if err != nil {
					t.Fatal(err)
				}
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(string(body)))}, nil
			})
			reader := NewMetricsReader(client, func() time.Time { return at })

			// Act
			metrics, err := reader.ReadMetrics(context.Background())

			// Assert
			if test.wantErr != nil {
				if !errors.Is(err, test.wantErr) {
					t.Errorf("error = %v, want %v", err, test.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !metrics.FetchedAt.Equal(at) || math.Abs(metrics.CPUUsagePercent-test.wantCPU) > 1e-9 || math.Abs(metrics.MemoryUsagePercent-test.wantMemory) > 1e-9 {
				t.Errorf("metrics = %+v, want fetched_at=%s, CPU=%v, memory=%v", metrics, at, test.wantCPU, test.wantMemory)
			}
		})
	}
}

func TestReadMetrics_PropagatesUpstreamFailure(t *testing.T) {
	// Arrange
	client := NewClient("http://prometheus:9090")
	client.httpClient.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	})
	reader := NewMetricsReader(client, time.Now)

	// Act
	_, err := reader.ReadMetrics(context.Background())

	// Assert
	if !errors.Is(err, ErrQueryFailed) {
		t.Errorf("error = %v, want ErrQueryFailed", err)
	}
}
