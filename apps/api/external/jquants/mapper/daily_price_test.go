package mapper

import (
	jquantsmodel "api/external/jquants/model"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestToDailyPrices_MapsAdjustedValues(t *testing.T) {
	// Arrange
	adjustedOpen := decimal.RequireFromString("2047.1")
	adjustedHigh := decimal.RequireFromString("2069.2")
	adjustedLow := decimal.RequireFromString("2035.3")
	adjustedClose := decimal.RequireFromString("2045.4")
	adjustedVolume := decimal.RequireFromString("2202500.5")

	input := []jquantsmodel.DailyPrice{
		{
			Date:           "2026-10-15",
			Code:           "86970",
			AdjustedOpen:   &adjustedOpen,
			AdjustedHigh:   &adjustedHigh,
			AdjustedLow:    &adjustedLow,
			AdjustedClose:  &adjustedClose,
			AdjustedVolume: &adjustedVolume,
		},
	}
	wantDate := time.Date(2026, time.October, 15, 0, 0, 0, 0, time.UTC)

	// Act
	dailyPrices, err := ToDailyPrices(input)

	// Assert
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(dailyPrices), 1; got != want {
		t.Fatalf("len(dailyPrices) = %d, want %d", got, want)
	}

	dailyPrice := dailyPrices[0]

	if got, want := dailyPrice.Symbol, "86970"; got != want {
		t.Errorf("Symbol = %q, want %q", got, want)
	}

	if got := dailyPrice.TradeDate; !got.Equal(wantDate) {
		t.Errorf("TradeDate = %v, want %v", got, wantDate)
	}

	assertDecimal(t, "OpenPrice", dailyPrice.OpenPrice, adjustedOpen)
	assertDecimal(t, "HighPrice", dailyPrice.HighPrice, adjustedHigh)
	assertDecimal(t, "LowPrice", dailyPrice.LowPrice, adjustedLow)
	assertDecimal(t, "ClosePrice", dailyPrice.ClosePrice, adjustedClose)
	assertDecimal(t, "Volume", dailyPrice.Volume, adjustedVolume)
}

func TestToDailyPrices_SkipsMissingAdjustedValues(t *testing.T) {
	tests := []struct {
		name   string
		remove func(*jquantsmodel.DailyPrice)
	}{
		{name: "調整済み始値だけが欠損", remove: func(price *jquantsmodel.DailyPrice) { price.AdjustedOpen = nil }},
		{name: "調整済み高値だけが欠損", remove: func(price *jquantsmodel.DailyPrice) { price.AdjustedHigh = nil }},
		{name: "調整済み安値だけが欠損", remove: func(price *jquantsmodel.DailyPrice) { price.AdjustedLow = nil }},
		{name: "調整済み終値だけが欠損", remove: func(price *jquantsmodel.DailyPrice) { price.AdjustedClose = nil }},
		{name: "調整済み出来高だけが欠損", remove: func(price *jquantsmodel.DailyPrice) { price.AdjustedVolume = nil }},
		{name: "調整済み項目がすべて欠損", remove: func(price *jquantsmodel.DailyPrice) {
			*price = jquantsmodel.DailyPrice{Date: price.Date, Code: price.Code}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			value := decimal.RequireFromString("100.5")
			price := jquantsmodel.DailyPrice{
				Date: "2026-10-15", Code: "86970",
				AdjustedOpen: &value, AdjustedHigh: &value, AdjustedLow: &value,
				AdjustedClose: &value, AdjustedVolume: &value,
			}
			test.remove(&price)

			// Act
			got, err := ToDailyPrices([]jquantsmodel.DailyPrice{price})

			// Assert
			if err != nil {
				t.Fatalf("ToDailyPrices() unexpected error: %v", err)
			}
			if len(got) != 0 {
				t.Errorf("ToDailyPrices() = %+v, want no results", got)
			}
		})
	}
}

func TestToDailyPrices_ReturnsOnlyValidPricesFromMixedInput(t *testing.T) {
	// Arrange
	value := decimal.RequireFromString("100.5")
	valid := jquantsmodel.DailyPrice{
		Date: "2026-10-15", Code: "72030",
		AdjustedOpen: &value, AdjustedHigh: &value, AdjustedLow: &value,
		AdjustedClose: &value, AdjustedVolume: &value,
	}
	missingVolume := valid
	missingVolume.Code = "86970"
	missingVolume.AdjustedVolume = nil
	anotherValid := valid
	anotherValid.Code = "67580"
	input := []jquantsmodel.DailyPrice{valid, missingVolume, anotherValid}
	wantSymbols := []string{"72030", "67580"}

	// Act
	got, err := ToDailyPrices(input)

	// Assert
	if err != nil {
		t.Fatalf("ToDailyPrices() unexpected error: %v", err)
	}
	if len(got) != len(wantSymbols) {
		t.Fatalf("len(results) = %d, want %d", len(got), len(wantSymbols))
	}
	for index, price := range got {
		if price.Symbol != wantSymbols[index] {
			t.Errorf("result[%d].Symbol = %s, want %s", index, price.Symbol, wantSymbols[index])
		}
	}
}

func TestToDailyPrices_InvalidDateReturnsError(t *testing.T) {
	// Arrange
	value := decimal.RequireFromString("1.0")

	input := []jquantsmodel.DailyPrice{
		{
			Date:           "invalid-date",
			Code:           "86970",
			AdjustedOpen:   &value,
			AdjustedHigh:   &value,
			AdjustedLow:    &value,
			AdjustedClose:  &value,
			AdjustedVolume: &value,
		},
	}

	// Act
	_, err := ToDailyPrices(input)

	// Assert
	if err == nil {
		t.Fatal("ToDailyPrices() error = nil, want error")
	}
}

func assertDecimal(
	t *testing.T,
	name string,
	actual decimal.Decimal,
	want decimal.Decimal,
) {
	t.Helper()

	if !actual.Equal(want) {
		t.Errorf("%s = %s, want %s", name, actual, want)
	}
}
