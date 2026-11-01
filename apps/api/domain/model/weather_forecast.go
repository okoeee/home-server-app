package model

import "time"

const WeatherForecastDays = 5

type WeatherForecast struct {
	LocationName string
	FetchedAt    time.Time
	Forecasts    []DailyWeatherForecast
}

type DailyWeatherForecast struct {
	Date                 time.Time
	Weather              WeatherCondition
	TemperatureMax       *float64
	TemperatureMin       *float64
	RelativeHumidityMean *float64
}
