package fetch

import (
	"api/external/jquants/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func (client *JQuantsClient) fetchDailyPrice(
	ctx context.Context,
	date time.Time,
	paginationKey string,
) (model.ListedDailyPrices, error) {

	endpoint, err := url.Parse(client.baseURL + equitiesBarsDaily)
	if err != nil {
		return model.ListedDailyPrices{}, fmt.Errorf("failed to parse endpoint: %w", err)
	}

	query := endpoint.Query()

	query.Set("date", date.Format(time.DateOnly))
	if paginationKey != "" {
		query.Set("pagination_key", paginationKey)
	}

	endpoint.RawQuery = query.Encode()

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint.String(),
		nil,
	)
	if err != nil {
		return model.ListedDailyPrices{}, fmt.Errorf("failed to create request: %w", err)
	}

	request.Header.Set("x-api-key", client.apiKey)

	response, err := client.httpClient.Do(request)
	if err != nil {
		return model.ListedDailyPrices{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return model.ListedDailyPrices{}, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var listedDailyPrices model.ListedDailyPrices

	err = json.NewDecoder(response.Body).Decode(&listedDailyPrices)
	if err != nil {
		return model.ListedDailyPrices{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return listedDailyPrices, nil

}

func (client *JQuantsClient) FetchDailyPrices(
	ctx context.Context,
	date time.Time,
) ([]model.DailyPrice, error) {

	var allListedDailyPrices []model.DailyPrice

	paginationKey := ""

	for {

		listedDailyPrice, err := client.fetchDailyPrice(ctx, date, paginationKey)
		if err != nil {
			return nil, err
		}

		allListedDailyPrices = append(allListedDailyPrices, listedDailyPrice.Data...)

		paginationKey = listedDailyPrice.PaginationKey

		if paginationKey == "" {
			break
		}
	}

	return allListedDailyPrices, nil

}
