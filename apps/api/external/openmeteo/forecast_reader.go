package openmeteo

import (
	domainmodel "api/domain/model"
	"api/external/openmeteo/mapper"
	openmeteomodel "api/external/openmeteo/model"
	"context"
	"time"
)

type ForecastFetcher interface {
	FetchForecast(context.Context) (openmeteomodel.Forecast, error)
}

type ForecastReader struct {
	fetcher ForecastFetcher
	now     func() time.Time
}

func NewForecastReader(fetcher ForecastFetcher, now func() time.Time) *ForecastReader {
	return &ForecastReader{fetcher: fetcher, now: now}
}

func (reader *ForecastReader) ReadForecast(ctx context.Context) (domainmodel.WeatherForecast, error) {
	forecast, err := reader.fetcher.FetchForecast(ctx)
	if err != nil {
		return domainmodel.WeatherForecast{}, err
	}
	return mapper.ToWeatherForecast(forecast, reader.now())
}
