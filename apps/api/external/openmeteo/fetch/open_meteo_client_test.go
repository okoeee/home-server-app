package fetch

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

type trackingBody struct {
	io.Reader
	closed bool
}

func (body *trackingBody) Close() error {
	body.closed = true
	return nil
}

func TestFetchForecast_ReturnsDailyDataAndPreservesNulls(t *testing.T) {
	// Arrange
	body := &trackingBody{Reader: strings.NewReader(`{
		"daily": {
			"time": ["2026-10-31", "2026-11-01", "2026-11-02", "2026-11-03", "2026-11-04"],
			"weather_code": [0, null, 61, 3, 1],
			"temperature_2m_max": [21.5, null, 18.3, 19.4, 20.6],
			"temperature_2m_min": [0, null, 13.5, 12.8, 12.1],
			"relative_humidity_2m_mean": [62, null, 81, 70, 59]
		}
	}`)}
	client := NewClient()
	client.httpClient.Transport = roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Method != http.MethodGet || request.URL.Scheme != "https" ||
			request.URL.Host != "api.open-meteo.com" || request.URL.Path != "/v1/forecast" {
			t.Errorf("unexpected request: %s %s", request.Method, request.URL)
		}
		for key, want := range map[string]string{
			"latitude": "35.689444", "longitude": "139.691667",
			"timezone": "Asia/Tokyo", "forecast_days": "5", "temperature_unit": "celsius",
			"daily": "weather_code,temperature_2m_max,temperature_2m_min,relative_humidity_2m_mean",
		} {
			if got := request.URL.Query().Get(key); got != want {
				t.Errorf("query[%s] = %q, want %q", key, got, want)
			}
		}
		return &http.Response{StatusCode: http.StatusOK, Body: body}, nil
	})

	// Act
	forecast, err := client.FetchForecast(context.Background())

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	daily := forecast.Daily
	if len(daily.Dates) != 5 || daily.Dates[0] != "2026-10-31" {
		t.Fatalf("Dates = %v, want five dates starting 2026-10-31", daily.Dates)
	}
	if len(daily.WeatherCodes) != 5 || daily.WeatherCodes[0] == nil || *daily.WeatherCodes[0] != 0 || daily.WeatherCodes[1] != nil {
		t.Errorf("WeatherCodes = %v, want clear sky and null preserved", daily.WeatherCodes)
	}
	for name, check := range map[string]struct {
		values []*float64
		want   float64
	}{
		"TemperatureMax":       {daily.TemperatureMax, 21.5},
		"TemperatureMin":       {daily.TemperatureMin, 0},
		"RelativeHumidityMean": {daily.RelativeHumidityMean, 62},
	} {
		if len(check.values) != 5 || check.values[0] == nil || *check.values[0] != check.want || check.values[1] != nil {
			t.Errorf("%s = %v, want first value %v and second null", name, check.values, check.want)
		}
	}
	if !body.closed {
		t.Error("response body was not closed")
	}
}

func TestFetchForecast_ResponseErrors(t *testing.T) {
	for _, test := range []struct {
		name   string
		status int
		body   string
		want   string
	}{
		{"ServerError", http.StatusInternalServerError, `internal details`, "unexpected HTTP status 500"},
		{"InvalidJSON", http.StatusOK, `{`, "decode forecast response"},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			body := &trackingBody{Reader: strings.NewReader(test.body)}
			client := NewClient()
			client.httpClient.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: test.status, Body: body}, nil
			})

			// Act
			_, err := client.FetchForecast(context.Background())

			// Assert
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Errorf("error = %v, want %q", err, test.want)
			}
			if !body.closed {
				t.Error("response body was not closed")
			}
		})
	}
}
