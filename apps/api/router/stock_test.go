package router

import (
	"api/db"
	"api/handler"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestGetDailyPrices_Scenarios(t *testing.T) {
	prices := []db.DailyPrice{{
		Symbol: "13010", TradeDate: time.Date(2026, 10, 19, 0, 0, 0, 0, time.UTC),
		OpenPrice: decimal.RequireFromString("100.1"), HighPrice: decimal.RequireFromString("110.2"),
		LowPrice: decimal.RequireFromString("90.3"), ClosePrice: decimal.RequireFromString("105.4"),
		Volume: decimal.RequireFromString("1200.5"),
	}}
	tests := []struct {
		name       string
		prices     []db.DailyPrice
		err        error
		wantStatus int
		wantBody   string
	}{
		{
			name:   "指定銘柄の日足を生成モデルからJSONへ変換する",
			prices: prices, wantStatus: http.StatusOK,
			wantBody: `[{"symbol":"13010","trade_date":"2026-10-19","open_price":"100.1","high_price":"110.2","low_price":"90.3","close_price":"105.4","volume":"1200.5"}]`,
		},
		{name: "対象なしは空配列", wantStatus: http.StatusOK, wantBody: `[]`},
		{
			name: "取得失敗は500", err: errors.New("read failed"),
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"message":"failed to get daily prices"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			queries := stockQuerierStub{pricesBySymbol: map[string][]db.DailyPrice{"13010": test.prices}, err: test.err}
			engine := NewRouter(
				handler.NewStockHandler(queries),
				handler.NewDailyPriceChangeHandler(priceChangeCalculatorStub{}),
				handler.NewDailyVolumeRatioHandler(volumeRatioCalculatorStub{}),
			)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/v1/stocks/13010/prices", nil)

			// Act
			engine.ServeHTTP(recorder, request)

			// Assert
			assertJSONResponse(t, recorder, test.wantStatus, test.wantBody)
		})
	}
}
