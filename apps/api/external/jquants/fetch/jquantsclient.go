package fetch

import (
	"errors"
	"net/http"
	"time"
)

const (
	defaultBaseURL = "https://api.jquants.com/v2"

	equitiesMasterPath = "/equities/master"
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
