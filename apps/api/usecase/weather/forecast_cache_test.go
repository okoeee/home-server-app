package weather

import (
	"api/domain/model"
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"
)

type forecastReaderFunc func(context.Context) (model.WeatherForecast, error)

func (f forecastReaderFunc) ReadForecast(ctx context.Context) (model.WeatherForecast, error) {
	return f(ctx)
}

func TestGet_CacheLifetime(t *testing.T) {
	for _, test := range []struct {
		name      string
		start     time.Time
		elapsed   time.Duration
		wantReads int
	}{
		{"WithinTTL", time.Date(2026, 11, 1, 9, 0, 0, 0, japanTimezone), 29*time.Minute + 59*time.Second, 1},
		{"ExpiresAt30Minutes", time.Date(2026, 11, 1, 9, 0, 0, 0, japanTimezone), 30 * time.Minute, 2},
		{"ClockMovesBack", time.Date(2026, 11, 1, 9, 0, 0, 0, japanTimezone), -time.Second, 2},
		{"JapanDateChanges", time.Date(2026, 11, 1, 14, 59, 0, 0, time.UTC), time.Minute, 2},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			now := test.start
			reads := 0
			cache := NewForecastCache(forecastReaderFunc(func(context.Context) (model.WeatherForecast, error) {
				reads++
				return model.WeatherForecast{LocationName: "東京都", FetchedAt: now}, nil
			}), func() time.Time { return now })
			first, err := cache.Get(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			now = now.Add(test.elapsed)

			// Act
			got, err := cache.Get(context.Background())

			// Assert
			if err != nil {
				t.Fatal(err)
			}
			if reads != test.wantReads {
				t.Errorf("reads = %d, want %d", reads, test.wantReads)
			}
			wantFetchedAt := first.FetchedAt
			if test.wantReads == 2 {
				wantFetchedAt = now
			}
			if !got.FetchedAt.Equal(wantFetchedAt) {
				t.Errorf("FetchedAt = %v, want %v", got.FetchedAt, wantFetchedAt)
			}
		})
	}
}

func TestGet_ReadFailureDoesNotOverwriteCache(t *testing.T) {
	for _, populated := range []bool{false, true} {
		name := "EmptyCache"
		if populated {
			name = "ExpiredCache"
		}
		t.Run(name, func(t *testing.T) {
			// Arrange
			now := time.Date(2026, 11, 1, 9, 0, 0, 0, japanTimezone)
			readErr := errors.New("upstream failed")
			shouldFail := !populated
			cache := NewForecastCache(forecastReaderFunc(func(context.Context) (model.WeatherForecast, error) {
				if shouldFail {
					return model.WeatherForecast{}, readErr
				}
				return model.WeatherForecast{LocationName: "東京都", FetchedAt: now}, nil
			}), func() time.Time { return now })
			if populated {
				if _, err := cache.Get(context.Background()); err != nil {
					t.Fatal(err)
				}
				now = now.Add(30 * time.Minute)
			}
			previous := cache.forecast
			shouldFail = true

			// Act
			got, err := cache.Get(context.Background())

			// Assert
			if !errors.Is(err, readErr) {
				t.Errorf("error = %v, want %v", err, readErr)
			}
			if !errors.Is(err, ErrForecastUnavailable) {
				t.Errorf("error = %v, want ErrForecastUnavailable", err)
			}
			if !reflect.DeepEqual(got, model.WeatherForecast{}) {
				t.Errorf("returned stale result: %+v", got)
			}
			if cache.forecast != previous {
				t.Error("failed read overwrote cache")
			}
			shouldFail = false
			retried, err := cache.Get(context.Background())
			if err != nil || !retried.FetchedAt.Equal(now) {
				t.Errorf("retry = %+v, %v", retried, err)
			}
		})
	}
}

func TestGet_ConcurrentRequestsShareSuccessfulRead(t *testing.T) {
	// Arrange
	now := time.Date(2026, 11, 1, 9, 0, 0, 0, japanTimezone)
	reads := 0
	cache := NewForecastCache(forecastReaderFunc(func(context.Context) (model.WeatherForecast, error) {
		reads++
		return model.WeatherForecast{LocationName: "東京都", FetchedAt: now}, nil
	}), func() time.Time { return now })
	const requests = 10
	start := make(chan struct{})
	var wait sync.WaitGroup

	// Act
	for range requests {
		wait.Go(func() {
			<-start
			got, err := cache.Get(context.Background())
			if err != nil || got.LocationName != "東京都" {
				t.Errorf("Get() = %+v, %v", got, err)
			}
		})
	}
	close(start)
	wait.Wait()

	// Assert
	if reads != 1 {
		t.Errorf("reads = %d, want 1", reads)
	}
}

func TestGet_CanceledRequestDoesNotRead(t *testing.T) {
	// Arrange
	reads := 0
	cache := NewForecastCache(forecastReaderFunc(func(context.Context) (model.WeatherForecast, error) {
		reads++
		return model.WeatherForecast{}, nil
	}), func() time.Time { return time.Time{} })
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	_, err := cache.Get(ctx)

	// Assert
	if !errors.Is(err, context.Canceled) {
		t.Errorf("error = %v, want context.Canceled", err)
	}
	if reads != 0 {
		t.Errorf("reads = %d, want 0", reads)
	}
}

func TestGet_ReturnedValuesCannotChangeCache(t *testing.T) {
	// Arrange
	now := time.Date(2026, 11, 1, 9, 0, 0, 0, japanTimezone)
	max, min, humidity := 20.5, 10.0, 60.0
	cache := NewForecastCache(forecastReaderFunc(func(context.Context) (model.WeatherForecast, error) {
		return model.WeatherForecast{FetchedAt: now, Forecasts: []model.DailyWeatherForecast{
			{Weather: model.WeatherClear, TemperatureMax: &max, TemperatureMin: &min, RelativeHumidityMean: &humidity},
			{Weather: model.WeatherUnknown},
		}}, nil
	}), func() time.Time { return now })

	// Act
	for range 2 {
		got, err := cache.Get(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		got.Forecasts[0].Weather = model.WeatherRain
		*got.Forecasts[0].TemperatureMax = 99
		*got.Forecasts[0].TemperatureMin = 99
		*got.Forecasts[0].RelativeHumidityMean = 99
	}
	got, err := cache.Get(context.Background())

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	day := got.Forecasts[0]
	if day.Weather != model.WeatherClear || *day.TemperatureMax != 20.5 || *day.TemperatureMin != 10 || *day.RelativeHumidityMean != 60 {
		t.Errorf("cache was modified: %+v", day)
	}
	if got.Forecasts[1].TemperatureMax != nil {
		t.Error("null value was not preserved")
	}
}
