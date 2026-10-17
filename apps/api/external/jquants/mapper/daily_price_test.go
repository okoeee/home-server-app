package mapper

import (
	jquantsmodel "api/external/jquants/model"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestToDailyPrices(t *testing.T) {
	adjustedOpen := decimal.RequireFromString("2047.1")
	adjustedHigh := decimal.RequireFromString("2069.2")
	adjustedLow := decimal.RequireFromString("2035.3")
	adjustedClose := decimal.RequireFromString("2045.4")
	adjustedVolume := decimal.RequireFromString("2202500.5")

	dailyPrices, err := ToDailyPrices([]jquantsmodel.DailyPrice{
		{
			Date:           "2026-10-15",
			Code:           "86970",
			AdjustedOpen:   &adjustedOpen,
			AdjustedHigh:   &adjustedHigh,
			AdjustedLow:    &adjustedLow,
			AdjustedClose:  &adjustedClose,
			AdjustedVolume: &adjustedVolume,
		},
	})
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

	wantDate := time.Date(2026, time.October, 15, 0, 0, 0, 0, time.UTC)
	if got := dailyPrice.TradeDate; !got.Equal(wantDate) {
		t.Errorf("TradeDate = %v, want %v", got, wantDate)
	}

	assertDecimal(t, "OpenPrice", dailyPrice.OpenPrice, adjustedOpen)
	assertDecimal(t, "HighPrice", dailyPrice.HighPrice, adjustedHigh)
	assertDecimal(t, "LowPrice", dailyPrice.LowPrice, adjustedLow)
	assertDecimal(t, "ClosePrice", dailyPrice.ClosePrice, adjustedClose)
	assertDecimal(t, "Volume", dailyPrice.Volume, adjustedVolume)
}

func TestToDailyPricesSkipsMissingAdjustedValues(t *testing.T) {
	dailyPrices, err := ToDailyPrices([]jquantsmodel.DailyPrice{
		{
			Date: "2026-10-15",
			Code: "86970",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if got := len(dailyPrices); got != 0 {
		t.Fatalf("len(dailyPrices) = %d, want 0", got)
	}
}

func TestToDailyPricesReturnsErrorForInvalidDate(t *testing.T) {
	value := decimal.RequireFromString("1.0")

	_, err := ToDailyPrices([]jquantsmodel.DailyPrice{
		{
			Date:           "invalid-date",
			Code:           "86970",
			AdjustedOpen:   &value,
			AdjustedHigh:   &value,
			AdjustedLow:    &value,
			AdjustedClose:  &value,
			AdjustedVolume: &value,
		},
	})
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
