package fetch

import (
	"api/external/jquants/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// 上場銘柄を1ページ取得
func (c *JQuantsClient) fetchEquityMaster(
	ctx context.Context,
	paginationKey string,
) (model.ListedEquityMaster, error) {

	endpoint, err := url.Parse(c.baseURL + equitiesMasterPath)
	if err != nil {
		return model.ListedEquityMaster{}, fmt.Errorf("error parsing endpoint: %w", err)
	}

	query := endpoint.Query()

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
		return model.ListedEquityMaster{}, fmt.Errorf("error creating request: %w", err)
	}

	request.Header.Set("x-api-key", c.apiKey)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return model.ListedEquityMaster{}, fmt.Errorf("error doing request: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return model.ListedEquityMaster{}, fmt.Errorf("unexpected response status: %s", response.Status)
	}

	var result model.ListedEquityMaster

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return model.ListedEquityMaster{}, fmt.Errorf("error decoding response: %w", err)
	}

	return result, nil
}

func (c *JQuantsClient) FetchEquityMasterList(
	ctx context.Context,
) ([]model.EquityMaster, error) {

	var allEquityMasters []model.EquityMaster

	paginationKey := ""

	for {
		response, err := c.fetchEquityMaster(ctx, paginationKey)
		if err != nil {
			return nil, err
		}

		allEquityMasters = append(allEquityMasters, response.Data...)

		if response.PaginationKey == "" {
			break
		}

		paginationKey = response.PaginationKey
	}

	return allEquityMasters, nil
}
