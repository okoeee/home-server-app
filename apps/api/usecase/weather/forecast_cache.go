package weather

import (
	"api/domain/model"
	"context"
	"fmt"
	"sync"
	"time"
)

const forecastCacheTTL = 30 * time.Minute

var japanTimezone = time.FixedZone("Asia/Tokyo", 9*60*60)

// ForecastReader は取得・変換・検証済みの予報を返す。
type ForecastReader interface {
	ReadForecast(context.Context) (model.WeatherForecast, error)
}

type ForecastCache struct {
	reader   ForecastReader
	now      func() time.Time
	mu       sync.Mutex
	forecast *model.WeatherForecast
}

func NewForecastCache(reader ForecastReader, now func() time.Time) *ForecastCache {
	return &ForecastCache{reader: reader, now: now}
}

func (cache *ForecastCache) Get(ctx context.Context) (model.WeatherForecast, error) {
	// 取得中もロックし、後続リクエストは取得が終わったキャッシュを再利用する。
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return model.WeatherForecast{}, err
	}
	now := cache.now()
	if cache.forecast != nil {
		age := now.Sub(cache.forecast.FetchedAt)
		sameDay := now.In(japanTimezone).Format(time.DateOnly) ==
			cache.forecast.FetchedAt.In(japanTimezone).Format(time.DateOnly)
		if age >= 0 && age < forecastCacheTTL && sameDay {
			return cloneForecast(*cache.forecast), nil
		}
	}

	forecast, err := cache.reader.ReadForecast(ctx)
	if err != nil {
		return model.WeatherForecast{}, fmt.Errorf("read weather forecast: %w", err)
	}
	// 成功した結果だけを保存する。返却値からキャッシュを書き換えられないようコピーする。
	cached := cloneForecast(forecast)
	cache.forecast = &cached
	return forecast, nil
}

func cloneForecast(forecast model.WeatherForecast) model.WeatherForecast {
	days := make([]model.DailyWeatherForecast, len(forecast.Forecasts))
	for i, day := range forecast.Forecasts {
		days[i] = day
		days[i].TemperatureMax = cloneFloat(day.TemperatureMax)
		days[i].TemperatureMin = cloneFloat(day.TemperatureMin)
		days[i].RelativeHumidityMean = cloneFloat(day.RelativeHumidityMean)
	}
	forecast.Forecasts = days
	return forecast
}

func cloneFloat(value *float64) *float64 {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
