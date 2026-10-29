package router

import (
	"api/domain/model"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestGetDailyPriceChanges_Scenarios(t *testing.T) {
	date := time.Date(2026, 10, 19, 0, 0, 0, 0, time.UTC)
	changes := []model.DailyPriceChange{{
		Symbol: "12340", TradeDate: date, ClosePrice: decimal.RequireFromString("110.5"),
		PreviousTradeDate: date.AddDate(0, 0, -3), PreviousClosePrice: decimal.NewFromInt(100),
		ChangeRate: decimal.RequireFromString("10.5"),
	}}
	tests := []struct {
		name       string
		query      string
		changes    []model.DailyPriceChange
		err        error
		wantStatus int
		wantBody   string
	}{
		{
			name:  "指定日の計算済み結果を文字列で返す",
			query: "?date=2026-10-19", changes: changes,
			wantStatus: http.StatusOK,
			wantBody:   `[{"symbol":"12340","trade_date":"2026-10-19","close_price":"110.5","previous_trade_date":"2026-10-16","previous_close_price":"100.0","change_rate":"10.5"}]`,
		},
		{
			name:       "結果なしはnullではなく空配列で返す",
			query:      "?date=2026-10-19",
			wantStatus: http.StatusOK, wantBody: `[]`,
		},
		{name: "日付未指定は400", wantStatus: http.StatusBadRequest, wantBody: `{"message":"date must be in YYYY-MM-DD format"}`},
		{name: "不正な日付形式は400", query: "?date=2026/10/19", wantStatus: http.StatusBadRequest, wantBody: `{"message":"date must be in YYYY-MM-DD format"}`},
		{name: "存在しない日付は400", query: "?date=2026-02-30", wantStatus: http.StatusBadRequest, wantBody: `{"message":"date must be in YYYY-MM-DD format"}`},
		{
			name:  "Usecase失敗は500",
			query: "?date=2026-10-19", err: errors.New("calculation failed"),
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"message":"failed to calculate daily price changes"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			calculator := priceChangeCalculatorStub{
				changesByDate: map[time.Time][]model.DailyPriceChange{date: test.changes},
				err:           test.err,
			}
			engine := newTestRouter(calculator, volumeRatioCalculatorStub{})
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/v1/daily-price-changes"+test.query, nil)

			// Act
			engine.ServeHTTP(recorder, request)

			// Assert
			assertJSONResponse(t, recorder, test.wantStatus, test.wantBody)
		})
	}
}
