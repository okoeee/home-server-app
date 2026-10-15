package mapper

import (
	jquantsmodel "api/external/jquants/model"
	"testing"
	"time"
)

func TestToStocks(t *testing.T) {
	equityMasters := []jquantsmodel.EquityMaster{
		{
			Date:                "2026-10-15",
			Code:                "72030",
			CompanyName:         "トヨタ自動車",
			Sector33Code:        "3700",
			Sector33Name:        "輸送用機器",
			MarketCode:          "0111",
			MarketName:          "プライム",
			ProductCategoryCode: "011",
		},
	}

	stocks, err := ToStocks(equityMasters)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(stocks), 1; got != want {
		t.Fatalf("len(stocks) = %d, want %d", got, want)
	}

	stock := stocks[0]

	if got, want := stock.Symbol, "72030"; got != want {
		t.Errorf("Symbol = %q, want %q", got, want)
	}

	if got, want := stock.CompanyName, "トヨタ自動車"; got != want {
		t.Errorf("CompanyName = %q, want %q", got, want)
	}

	wantDate := time.Date(2026, time.October, 15, 0, 0, 0, 0, time.UTC)
	if got := stock.EffectiveDate; !got.Equal(wantDate) {
		t.Errorf("EffectiveDate = %v, want %v", got, wantDate)
	}

	if got, want := stock.MarketName, "プライム"; got != want {
		t.Errorf("MarketName = %q, want %q", got, want)
	}
}

func TestToStocksReturnsErrorForInvalidDate(t *testing.T) {
	equityMasters := []jquantsmodel.EquityMaster{
		{
			Date: "invalid-date",
			Code: "72030",
		},
	}

	_, err := ToStocks(equityMasters)
	if err == nil {
		t.Fatal("ToStocks() error = nil, want error")
	}
}
