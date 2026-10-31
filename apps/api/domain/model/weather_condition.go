package model

type WeatherCondition string

const (
	WeatherClear                WeatherCondition = "clear"
	WeatherMainlyClear          WeatherCondition = "mainly_clear"
	WeatherPartlyCloudy         WeatherCondition = "partly_cloudy"
	WeatherOvercast             WeatherCondition = "overcast"
	WeatherFog                  WeatherCondition = "fog"
	WeatherDrizzle              WeatherCondition = "drizzle"
	WeatherFreezingDrizzle      WeatherCondition = "freezing_drizzle"
	WeatherRain                 WeatherCondition = "rain"
	WeatherFreezingRain         WeatherCondition = "freezing_rain"
	WeatherSnow                 WeatherCondition = "snow"
	WeatherSnowGrains           WeatherCondition = "snow_grains"
	WeatherRainShowers          WeatherCondition = "rain_showers"
	WeatherSnowShowers          WeatherCondition = "snow_showers"
	WeatherThunderstorm         WeatherCondition = "thunderstorm"
	WeatherThunderstormWithHail WeatherCondition = "thunderstorm_with_hail"
	WeatherUnknown              WeatherCondition = "unknown"
)
