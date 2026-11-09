package prometheus

import (
	"api/domain/model"
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	nodeLabels           = `job="node",instance="node_exporter:9100"`
	upQuery              = `up{` + nodeLabels + `}`
	cpuUsageQuery        = `100 * (1 - avg(rate(node_cpu_seconds_total{` + nodeLabels + `,mode="idle"}[1m])))`
	memoryTotalQuery     = `node_memory_MemTotal_bytes{` + nodeLabels + `}`
	memoryAvailableQuery = `node_memory_MemAvailable_bytes{` + nodeLabels + `}`
)

var ErrMetricsUnavailable = errors.New("server metrics unavailable")

type MetricsReader struct {
	client *Client
	now    func() time.Time
}

func NewMetricsReader(client *Client, now func() time.Time) *MetricsReader {
	return &MetricsReader{client: client, now: now}
}

func (reader *MetricsReader) ReadMetrics(ctx context.Context) (model.ServerMetrics, error) {
	// 複数クエリ全体を10秒以内に制限する。各リクエストで期限を延長しない。
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()
	at := reader.now()

	values := make(map[string]float64, 4)
	for _, expression := range []string{upQuery, cpuUsageQuery, memoryTotalQuery, memoryAvailableQuery} {
		value, found, err := reader.readOne(ctx, expression, at)
		if err != nil {
			return model.ServerMetrics{}, err
		}
		if !found {
			return model.ServerMetrics{}, fmt.Errorf("%w: missing result for %s", ErrMetricsUnavailable, expression)
		}
		values[expression] = value
	}
	if values[upQuery] == 0 {
		return model.ServerMetrics{}, fmt.Errorf("%w: node exporter scrape failed", ErrMetricsUnavailable)
	}
	if values[upQuery] != 1 || values[cpuUsageQuery] < 0 || values[cpuUsageQuery] > 100 ||
		values[memoryTotalQuery] <= 0 || values[memoryAvailableQuery] < 0 || values[memoryAvailableQuery] > values[memoryTotalQuery] {
		return model.ServerMetrics{}, fmt.Errorf("%w: metrics outside valid ranges", ErrQueryFailed)
	}
	return model.ServerMetrics{
		FetchedAt:          at,
		CPUUsagePercent:    values[cpuUsageQuery],
		MemoryUsagePercent: 100 * (1 - values[memoryAvailableQuery]/values[memoryTotalQuery]),
	}, nil
}

func (reader *MetricsReader) readOne(ctx context.Context, expression string, at time.Time) (float64, bool, error) {
	result, err := reader.client.query(ctx, expression, at)
	if err != nil {
		return 0, false, err
	}
	if len(result) == 0 {
		return 0, false, nil
	}
	if len(result) != 1 {
		return 0, false, fmt.Errorf("%w: expected one series for %s", ErrQueryFailed, expression)
	}
	return result[0], true, nil
}
