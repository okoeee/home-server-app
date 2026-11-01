package mapper

import (
	domainmodel "api/domain/model"
	openmeteomodel "api/external/openmeteo/model"
	"fmt"
	"time"
)

// 日本時間はUTC+9固定。タイムゾーンDBを持たないコンテナでも利用できる。
var japanTimezone = time.FixedZone("Asia/Tokyo", 9*60*60)

// ToWeatherForecast は取得日時の日本時間での暦日を基準に、今日から5日分を検証・変換する。
// 日別データの仕様: https://open-meteo.com/en/docs#daily-parameter-definition
func ToWeatherForecast(forecast openmeteomodel.Forecast, fetchedAt time.Time) (domainmodel.WeatherForecast, error) {
	daily := forecast.Daily
	if len(daily.Dates) != domainmodel.WeatherForecastDays ||
		len(daily.WeatherCodes) != domainmodel.WeatherForecastDays ||
		len(daily.TemperatureMax) != domainmodel.WeatherForecastDays ||
		len(daily.TemperatureMin) != domainmodel.WeatherForecastDays ||
		len(daily.RelativeHumidityMean) != domainmodel.WeatherForecastDays {
		return domainmodel.WeatherForecast{}, fmt.Errorf("daily forecast arrays must contain %d entries each", domainmodel.WeatherForecastDays)
	}

	localTime := fetchedAt.In(japanTimezone)
	today := time.Date(localTime.Year(), localTime.Month(), localTime.Day(), 0, 0, 0, 0, japanTimezone)
	result := domainmodel.WeatherForecast{
		LocationName: "東京都",
		FetchedAt:    localTime,
		Forecasts:    make([]domainmodel.DailyWeatherForecast, 0, domainmodel.WeatherForecastDays),
	}
	for i, dateString := range daily.Dates {
		date, err := time.ParseInLocation(time.DateOnly, dateString, japanTimezone)
		if err != nil {
			return domainmodel.WeatherForecast{}, fmt.Errorf("parse forecast date %q: %w", dateString, err)
		}
		if !date.Equal(today.AddDate(0, 0, i)) {
			return domainmodel.WeatherForecast{}, fmt.Errorf("forecast date %q must be day %d from today", dateString, i)
		}

		humidity := daily.RelativeHumidityMean[i]
		if humidity != nil && (*humidity < 0 || *humidity > 100) {
			return domainmodel.WeatherForecast{}, fmt.Errorf("relative humidity for %s must be between 0 and 100", dateString)
		}

		result.Forecasts = append(result.Forecasts, domainmodel.DailyWeatherForecast{
			Date:                 date,
			Weather:              toWeatherCondition(daily.WeatherCodes[i]),
			TemperatureMax:       daily.TemperatureMax[i],
			TemperatureMin:       daily.TemperatureMin[i],
			RelativeHumidityMean: humidity,
		})
	}
	return result, nil
}

// 天気コードの仕様: https://open-meteo.com/en/docs#weathervariables
func toWeatherCondition(code *int) domainmodel.WeatherCondition {
	if code == nil {
		return domainmodel.WeatherUnknown
	}
	switch *code {
	case 0:
		return domainmodel.WeatherClear
	case 1:
		return domainmodel.WeatherMainlyClear
	case 2:
		return domainmodel.WeatherPartlyCloudy
	case 3:
		return domainmodel.WeatherOvercast
	case 45, 48:
		return domainmodel.WeatherFog
	case 51, 53, 55:
		return domainmodel.WeatherDrizzle
	case 56, 57:
		return domainmodel.WeatherFreezingDrizzle
	case 61, 63, 65:
		return domainmodel.WeatherRain
	case 66, 67:
		return domainmodel.WeatherFreezingRain
	case 71, 73, 75:
		return domainmodel.WeatherSnow
	case 77:
		return domainmodel.WeatherSnowGrains
	case 80, 81, 82:
		return domainmodel.WeatherRainShowers
	case 85, 86:
		return domainmodel.WeatherSnowShowers
	case 95, 97:
		return domainmodel.WeatherThunderstorm
	case 96, 99:
		return domainmodel.WeatherThunderstormWithHail
	default:
		return domainmodel.WeatherUnknown
	}
}
