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

func TestGetDailyVolumeRatios_Scenarios(t *testing.T) {
	date := time.Date(2026, 10, 19, 0, 0, 0, 0, time.UTC)
	ratios := []model.DailyVolumeRatio{
		{Symbol: "12340", TradeDate: date, Volume: decimal.RequireFromString("150.3"), Ratio: decimal.RequireFromString("1.5")},
	}
	tests := []struct {
		name       string
		query      string
		ratios     []model.DailyVolumeRatio
		err        error
		wantStatus int
		wantBody   string
	}{
		{
			name:  "計算済みの出来高と倍率を文字列で返す",
			query: "?date=2026-10-19", ratios: ratios,
			wantStatus: http.StatusOK,
			wantBody:   `[{"symbol":"12340","trade_date":"2026-10-19","volume":"150.3","ratio":"1.5"}]`,
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
			wantBody:   `{"message":"failed to calculate daily volume ratios"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			calculator := volumeRatioCalculatorStub{
				ratiosByDate: map[time.Time][]model.DailyVolumeRatio{date: test.ratios},
				err:          test.err,
			}
			engine := newTestRouter(priceChangeCalculatorStub{}, calculator)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/v1/daily-volume-ratios"+test.query, nil)

			// Act
			engine.ServeHTTP(recorder, request)

			// Assert
			assertJSONResponse(t, recorder, test.wantStatus, test.wantBody)
		})
	}
}
