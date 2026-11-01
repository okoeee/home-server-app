package mapper

import (
	domainmodel "api/domain/model"
	openmeteomodel "api/external/openmeteo/model"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func pointer[T any](value T) *T { return &value }

func validForecast() openmeteomodel.Forecast {
	return openmeteomodel.Forecast{Daily: openmeteomodel.DailyForecast{
		Dates:                []string{"2026-11-01", "2026-11-02", "2026-11-03", "2026-11-04", "2026-11-05"},
		WeatherCodes:         []*int{pointer(0), pointer(61), nil, pointer(3), pointer(999)},
		TemperatureMax:       []*float64{pointer(21.5), nil, nil, pointer(18.0), pointer(17.0)},
		TemperatureMin:       []*float64{pointer(0.0), pointer(-1.5), nil, pointer(9.0), pointer(8.0)},
		RelativeHumidityMean: []*float64{pointer(62.0), nil, nil, pointer(0.0), pointer(100.0)},
	}}
}

func TestToWeatherForecast_ConvertsFiveDaysInJapanTime(t *testing.T) {
	// Arrange
	forecast := validForecast()
	// UTCでは10月31日だが、日本時間では11月1日。
	fetchedAt := time.Date(2026, 10, 31, 15, 30, 0, 0, time.UTC)
	want := domainmodel.WeatherForecast{
		LocationName: "東京都",
		FetchedAt:    time.Date(2026, 11, 1, 0, 30, 0, 0, japanTimezone),
		Forecasts: []domainmodel.DailyWeatherForecast{
			{Date: time.Date(2026, 11, 1, 0, 0, 0, 0, japanTimezone), Weather: domainmodel.WeatherClear,
				TemperatureMax: pointer(21.5), TemperatureMin: pointer(0.0), RelativeHumidityMean: pointer(62.0)},
			{Date: time.Date(2026, 11, 2, 0, 0, 0, 0, japanTimezone), Weather: domainmodel.WeatherRain,
				TemperatureMin: pointer(-1.5)},
			{Date: time.Date(2026, 11, 3, 0, 0, 0, 0, japanTimezone), Weather: domainmodel.WeatherUnknown},
			{Date: time.Date(2026, 11, 4, 0, 0, 0, 0, japanTimezone), Weather: domainmodel.WeatherOvercast,
				TemperatureMax: pointer(18.0), TemperatureMin: pointer(9.0), RelativeHumidityMean: pointer(0.0)},
			{Date: time.Date(2026, 11, 5, 0, 0, 0, 0, japanTimezone), Weather: domainmodel.WeatherUnknown,
				TemperatureMax: pointer(17.0), TemperatureMin: pointer(8.0), RelativeHumidityMean: pointer(100.0)},
		},
	}

	// Act
	got, err := ToWeatherForecast(forecast, fetchedAt)

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToWeatherForecast() = %#v, want %#v", got, want)
	}
}

func TestToWeatherForecast_InvalidDailyData(t *testing.T) {
	for _, test := range []struct {
		name   string
		change func(*openmeteomodel.DailyForecast)
	}{
		{"MissingDaily", func(d *openmeteomodel.DailyForecast) { *d = openmeteomodel.DailyForecast{} }},
		{"TooFewDates", func(d *openmeteomodel.DailyForecast) { d.Dates = d.Dates[:4] }},
		{"TooManyDates", func(d *openmeteomodel.DailyForecast) { d.Dates = append(d.Dates, "2026-11-06") }},
		{"WeatherCodesLength", func(d *openmeteomodel.DailyForecast) { d.WeatherCodes = d.WeatherCodes[:4] }},
		{"MaxTemperaturesLength", func(d *openmeteomodel.DailyForecast) { d.TemperatureMax = d.TemperatureMax[:4] }},
		{"MinTemperaturesLength", func(d *openmeteomodel.DailyForecast) { d.TemperatureMin = d.TemperatureMin[:4] }},
		{"HumidityLength", func(d *openmeteomodel.DailyForecast) { d.RelativeHumidityMean = d.RelativeHumidityMean[:4] }},
		{"InvalidDate", func(d *openmeteomodel.DailyForecast) { d.Dates[1] = "2026-11-31" }},
		{"DoesNotStartToday", func(d *openmeteomodel.DailyForecast) { d.Dates[0] = "2026-10-31" }},
		{"DuplicateDate", func(d *openmeteomodel.DailyForecast) { d.Dates[1] = d.Dates[0] }},
		{"DatesOutOfOrder", func(d *openmeteomodel.DailyForecast) { d.Dates[1], d.Dates[2] = d.Dates[2], d.Dates[1] }},
		{"HumidityBelowZero", func(d *openmeteomodel.DailyForecast) { d.RelativeHumidityMean[0] = pointer(-0.1) }},
		{"HumidityAbove100", func(d *openmeteomodel.DailyForecast) { d.RelativeHumidityMean[0] = pointer(100.1) }},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			forecast := validForecast()
			test.change(&forecast.Daily)
			fetchedAt := time.Date(2026, 11, 1, 9, 0, 0, 0, japanTimezone)

			// Act
			got, err := ToWeatherForecast(forecast, fetchedAt)

			// Assert
			if err == nil {
				t.Fatal("expected invalid daily data to return an error")
			}
			if !reflect.DeepEqual(got, domainmodel.WeatherForecast{}) {
				t.Errorf("got partial forecast on error: %#v", got)
			}
		})
	}
}

func TestToWeatherCondition_MapsCodes(t *testing.T) {
	for _, test := range []struct {
		codes []int
		want  domainmodel.WeatherCondition
	}{
		{[]int{0}, domainmodel.WeatherClear},
		{[]int{1}, domainmodel.WeatherMainlyClear},
		{[]int{2}, domainmodel.WeatherPartlyCloudy},
		{[]int{3}, domainmodel.WeatherOvercast},
		{[]int{45, 48}, domainmodel.WeatherFog},
		{[]int{51, 53, 55}, domainmodel.WeatherDrizzle},
		{[]int{56, 57}, domainmodel.WeatherFreezingDrizzle},
		{[]int{61, 63, 65}, domainmodel.WeatherRain},
		{[]int{66, 67}, domainmodel.WeatherFreezingRain},
		{[]int{71, 73, 75}, domainmodel.WeatherSnow},
		{[]int{77}, domainmodel.WeatherSnowGrains},
		{[]int{80, 81, 82}, domainmodel.WeatherRainShowers},
		{[]int{85, 86}, domainmodel.WeatherSnowShowers},
		{[]int{95, 97}, domainmodel.WeatherThunderstorm},
		{[]int{96, 99}, domainmodel.WeatherThunderstormWithHail},
		{[]int{-1, 999}, domainmodel.WeatherUnknown},
	} {
		for _, code := range test.codes {
			t.Run(strconv.Itoa(code), func(t *testing.T) {
				// Arrange
				input := pointer(code)

				// Act
				got := toWeatherCondition(input)

				// Assert
				if got != test.want {
					t.Errorf("toWeatherCondition(%d) = %s, want %s", code, got, test.want)
				}
			})
		}
	}
}
