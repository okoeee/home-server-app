package router

import (
	"api/domain/model"
	"api/external/openmeteo"
	openmeteomodel "api/external/openmeteo/model"
	"api/handler"
	"api/usecase/weather"
	"bytes"
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type weatherForecastStub struct {
	err error
}

func (stub weatherForecastStub) Get(context.Context) (model.WeatherForecast, error) {
	return model.WeatherForecast{}, stub.err
}

type forecastFetcherStub struct {
	forecast openmeteomodel.Forecast
	err      error
	reads    int
}

func (stub *forecastFetcherStub) FetchForecast(context.Context) (openmeteomodel.Forecast, error) {
	stub.reads++
	return stub.forecast, stub.err
}

func newWeatherTestRouter(getter handler.WeatherForecastGetter) *gin.Engine {
	return NewRouter(
		handler.NewStockHandler(stockQuerierStub{}),
		handler.NewDailyPriceChangeHandler(priceChangeCalculatorStub{}),
		handler.NewDailyVolumeRatioHandler(volumeRatioCalculatorStub{}),
		handler.NewWeatherHandler(getter),
	)
}

func TestGetForecast_ReturnsCachedForecast(t *testing.T) {
	// Arrange
	fetchedAt := time.Date(2026, 11, 1, 0, 30, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))
	now := fetchedAt
	clock := func() time.Time { return now }
	clear, rain, maxTemp, minTemp, humidity := 0, 61, 21.5, 0.0, 62.0
	fetcher := &forecastFetcherStub{forecast: openmeteomodel.Forecast{Daily: openmeteomodel.DailyForecast{
		Dates:                []string{"2026-11-01", "2026-11-02", "2026-11-03", "2026-11-04", "2026-11-05"},
		WeatherCodes:         []*int{&clear, &rain, nil, nil, nil},
		TemperatureMax:       []*float64{&maxTemp, nil, nil, nil, nil},
		TemperatureMin:       []*float64{&minTemp, nil, nil, nil, nil},
		RelativeHumidityMean: []*float64{&humidity, nil, nil, nil, nil},
	}}}
	reader := openmeteo.NewForecastReader(fetcher, clock)
	cache := weather.NewForecastCache(reader, clock)
	engine := newWeatherTestRouter(cache)
	want := `{
		"location_name":"東京都", "fetched_at":"2026-11-01T00:30:00+09:00",
		"forecasts":[
			{"date":"2026-11-01","weather":"clear","temperature_max":21.5,"temperature_min":0,"relative_humidity_mean":62},
			{"date":"2026-11-02","weather":"rain","temperature_max":null,"temperature_min":null,"relative_humidity_mean":null},
			{"date":"2026-11-03","weather":"unknown","temperature_max":null,"temperature_min":null,"relative_humidity_mean":null},
			{"date":"2026-11-04","weather":"unknown","temperature_max":null,"temperature_min":null,"relative_humidity_mean":null},
			{"date":"2026-11-05","weather":"unknown","temperature_max":null,"temperature_min":null,"relative_humidity_mean":null}
		]
	}`

	// Act
	first, second := httptest.NewRecorder(), httptest.NewRecorder()
	engine.ServeHTTP(first, httptest.NewRequest(http.MethodGet, "/v1/weather/forecast", nil))
	now = now.Add(time.Minute)
	engine.ServeHTTP(second, httptest.NewRequest(http.MethodGet, "/v1/weather/forecast", nil))

	// Assert
	for _, response := range []*httptest.ResponseRecorder{first, second} {
		assertJSONResponse(t, response, http.StatusOK, want)
		if !strings.HasPrefix(response.Header().Get("Content-Type"), "application/json") {
			t.Errorf("Content-Type = %q", response.Header().Get("Content-Type"))
		}
	}
	if fetcher.reads != 1 {
		t.Errorf("fetches = %d, want 1", fetcher.reads)
	}
}

func TestGetForecast_UpstreamErrorsReturnBadGateway(t *testing.T) {
	for _, test := range []struct {
		name    string
		err     error
		wantLog string
	}{
		{"FetchFailure", errors.New("private upstream details"), "private upstream details"},
		{"InvalidForecast", nil, "daily forecast arrays must contain 5 entries each"},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			var logs bytes.Buffer
			originalWriter := log.Writer()
			log.SetOutput(&logs)
			t.Cleanup(func() { log.SetOutput(originalWriter) })
			clock := func() time.Time { return time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC) }
			fetcher := &forecastFetcherStub{err: test.err}
			reader := openmeteo.NewForecastReader(fetcher, clock)
			engine := newWeatherTestRouter(weather.NewForecastCache(reader, clock))
			recorder := httptest.NewRecorder()

			// Act
			engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/v1/weather/forecast", nil))

			// Assert
			assertJSONResponse(t, recorder, http.StatusBadGateway, `{"message":"天気予報を取得できませんでした"}`)
			if !strings.Contains(logs.String(), "get weather forecast:") || !strings.Contains(logs.String(), test.wantLog) {
				t.Errorf("logs = %q, want forecast error including %q", logs.String(), test.wantLog)
			}
		})
	}
}

func TestGetForecast_InternalErrorReturnsInternalServerError(t *testing.T) {
	// Arrange
	engine := newWeatherTestRouter(weatherForecastStub{err: errors.New("private internal details")})
	recorder := httptest.NewRecorder()

	// Act
	engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/v1/weather/forecast", nil))

	// Assert
	assertJSONResponse(t, recorder, http.StatusInternalServerError, `{"message":"天気予報を取得できませんでした"}`)
}
