package model

// Forecast はOpen-Meteoの予報レスポンス。
type Forecast struct {
	Daily DailyForecast `json:"daily"`
}

// DailyForecast は同じ添字が同じ日付を表す日別データ。
type DailyForecast struct {
	Dates []string `json:"time"`
	// 0（快晴）と欠損を区別するためポインターを使う。
	WeatherCodes []*int `json:"weather_code"`
	// 気温・湿度も0と欠損を区別し、欠損はnilのまま保持する。
	TemperatureMax       []*float64 `json:"temperature_2m_max"`
	TemperatureMin       []*float64 `json:"temperature_2m_min"`
	RelativeHumidityMean []*float64 `json:"relative_humidity_2m_mean"`
}
