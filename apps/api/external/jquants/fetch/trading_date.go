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

func (client *JQuantsClient) FetchTradingDays(
	ctx context.Context,
	from time.Time,
	to time.Time,
) (model.ListedTradingCalendar, error) {

	endpoint, err := url.Parse(client.baseURL + marketsCalendar)
	if err != nil {
		return model.ListedTradingCalendar{}, fmt.Errorf("error parsing url: %w", err)
	}

	query := endpoint.Query()

	// 営業日のみ
	// [https://jpx-jquants.com/ja/spec/mkt-cal/holiday-division]
	query.Set("hol_div", "1")
	query.Set("from", from.Format(time.DateOnly))
	query.Set("to", to.Format(time.DateOnly))

	endpoint.RawQuery = query.Encode()

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint.String(),
		nil,
	)
	if err != nil {
		return model.ListedTradingCalendar{}, fmt.Errorf("error creating request: %w", err)
	}

	request.Header.Set("x-api-key", client.apiKey)

	response, err := client.sendRequestWithInterval(ctx, request)
	if err != nil {
		return model.ListedTradingCalendar{}, fmt.Errorf("error doing request: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return model.ListedTradingCalendar{}, fmt.Errorf("error fetching trading date: %d %s", response.StatusCode, response.Status)
	}

	var result model.ListedTradingCalendar

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return model.ListedTradingCalendar{}, fmt.Errorf("error decoding response: %w", err)
	}

	return result, nil

}
