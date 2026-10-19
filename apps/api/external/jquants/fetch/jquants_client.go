package fetch

import (
	"context"
	"errors"
	"net/http"
	"time"
)

const (
	defaultBaseURL = "https://api.jquants.com/v2"

	equitiesMasterPath = "/equities/master"

	equitiesBarsDaily = "/equities/bars/daily"

	marketsCalendar = "/markets/calendar"
)

type JQuantsClient struct {
	baseURL string

	apiKey string

	httpClient *http.Client
}

func NewClient(apiKey string) (*JQuantsClient, error) {
	if apiKey == "" {
		return nil, errors.New("API key is empty")
	}

	return &JQuantsClient{
		baseURL: defaultBaseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil

}

// Lightプランの上限は60リクエスト/分のため、余裕を持って1.1秒待機する
const requestInterval = 1100 * time.Millisecond

func (client *JQuantsClient) sendRequestWithInterval(
	ctx context.Context,
	request *http.Request,
) (*http.Response, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(requestInterval):
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil

}
