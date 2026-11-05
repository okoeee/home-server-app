package handler

import (
	"api/domain/model"
	"api/openapi"
	"api/usecase/weather"
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type WeatherForecastGetter interface {
	Get(context.Context) (model.WeatherForecast, error)
}

type WeatherHandler struct {
	getter WeatherForecastGetter
}

func NewWeatherHandler(getter WeatherForecastGetter) *WeatherHandler {
	return &WeatherHandler{getter: getter}
}

func (handler *WeatherHandler) GetForecast(c *gin.Context) {
	forecast, err := handler.getter.Get(c.Request.Context())
	if err != nil {
		log.Printf("get weather forecast: %v", err)
		status := http.StatusInternalServerError
		if errors.Is(err, weather.ErrForecastUnavailable) {
			status = http.StatusBadGateway
		}
		c.JSON(status, openapi.ErrorResponse{Message: "天気予報を取得できませんでした"})
		return
	}

	days := make([]openapi.DailyWeatherForecastResponse, 0, len(forecast.Forecasts))
	for _, day := range forecast.Forecasts {
		days = append(days, openapi.DailyWeatherForecastResponse{
			Date:                 day.Date.Format(time.DateOnly),
			Weather:              openapi.WeatherCondition(day.Weather),
			TemperatureMax:       day.TemperatureMax,
			TemperatureMin:       day.TemperatureMin,
			RelativeHumidityMean: day.RelativeHumidityMean,
		})
	}
	c.JSON(http.StatusOK, openapi.WeatherForecastResponse{
		LocationName: forecast.LocationName,
		FetchedAt:    forecast.FetchedAt,
		Forecasts:    days,
	})
}
