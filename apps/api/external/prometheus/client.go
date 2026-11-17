package prometheus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const requestTimeout = 10 * time.Second

var ErrQueryFailed = errors.New("prometheus query failed")

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Timeout: requestTimeout},
	}
}

// query の応答仕様: https://prometheus.io/docs/prometheus/latest/querying/api/#instant-queries
func (client *Client) query(
	ctx context.Context,
	expression string,
	at time.Time,
) ([]float64, error) {
	parameters := url.Values{}
	parameters.Set("query", expression)
	parameters.Set("time", at.Format(time.RFC3339Nano))
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, client.baseURL+"/api/v1/query?"+parameters.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: create request: %w", ErrQueryFailed, err)
	}
	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%w: send request: %w", ErrQueryFailed, err)
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: HTTP status %d", ErrQueryFailed, response.StatusCode)
	}

	var payload struct {
		Status string `json:"status"`
		Error  string `json:"error"`
		Data   struct {
			ResultType string `json:"resultType"`
			Result     []struct {
				Value []json.RawMessage `json:"value"`
			} `json:"result"`
		} `json:"data"`
	}
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&payload); err != nil {
		return nil, fmt.Errorf("%w: decode response: %w", ErrQueryFailed, err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, fmt.Errorf("%w: unexpected trailing response data", ErrQueryFailed)
	}
	if payload.Status != "success" || payload.Data.ResultType != "vector" || payload.Data.Result == nil {
		return nil, fmt.Errorf("%w: unexpected response status or result: %s", ErrQueryFailed, payload.Error)
	}

	result := make([]float64, 0, len(payload.Data.Result))
	for _, item := range payload.Data.Result {
		if len(item.Value) != 2 {
			return nil, fmt.Errorf("%w: sample must contain timestamp and value", ErrQueryFailed)
		}
		var timestamp *float64
		var valueString *string
		if err := json.Unmarshal(item.Value[0], &timestamp); err != nil || timestamp == nil {
			return nil, fmt.Errorf("%w: invalid sample timestamp", ErrQueryFailed)
		}
		if err := json.Unmarshal(item.Value[1], &valueString); err != nil || valueString == nil {
			return nil, fmt.Errorf("%w: invalid sample value", ErrQueryFailed)
		}
		value, err := strconv.ParseFloat(*valueString, 64)
		if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
			return nil, fmt.Errorf("%w: invalid numeric value %q", ErrQueryFailed, *valueString)
		}
		result = append(result, value)
	}
	return result, nil
}
