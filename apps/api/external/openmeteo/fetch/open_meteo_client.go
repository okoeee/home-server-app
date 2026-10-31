package fetch

import (
	"api/external/openmeteo/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	forecastURL    = "https://api.open-meteo.com/v1/forecast"
	requestTimeout = 10 * time.Second

	tokyoCityHallLatitude  = "35.689444"
	tokyoCityHallLongitude = "139.691667"
	forecastTimezone       = "Asia/Tokyo"
	forecastDays           = "5"
)

type OpenMeteoClient struct {
	httpClient *http.Client
}

func NewClient() *OpenMeteoClient {
	return &OpenMeteoClient{
		httpClient: &http.Client{Timeout: requestTimeout},
	}
}

// FetchForecast は東京都庁を代表地点として、今日を含む5日間の予報を取得する。
// 日付・配列の整合性確認とアプリ内モデルへの変換は呼び出し側で行う。
func (client *OpenMeteoClient) FetchForecast(ctx context.Context) (model.Forecast, error) {
	query := url.Values{}
	query.Set("latitude", tokyoCityHallLatitude)
	query.Set("longitude", tokyoCityHallLongitude)
	query.Set("timezone", forecastTimezone)
	query.Set("forecast_days", forecastDays)
	query.Set("temperature_unit", "celsius")
	query.Set("daily", "weather_code,temperature_2m_max,temperature_2m_min,relative_humidity_2m_mean")

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, forecastURL+"?"+query.Encode(), nil)
	if err != nil {
		return model.Forecast{}, fmt.Errorf("create forecast request: %w", err)
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return model.Forecast{}, fmt.Errorf("fetch forecast: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return model.Forecast{}, fmt.Errorf("fetch forecast: unexpected HTTP status %d", response.StatusCode)
	}

	var forecast model.Forecast
	if err := json.NewDecoder(response.Body).Decode(&forecast); err != nil {
		return model.Forecast{}, fmt.Errorf("decode forecast response: %w", err)
	}

	return forecast, nil
}
