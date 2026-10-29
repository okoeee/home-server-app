package router

import (
	"api/domain/model"
	"api/handler"
	"api/usecase/dailyprice"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestGetLatestDailyPriceChanges(t *testing.T) {
	date := time.Date(2026, 10, 19, 0, 0, 0, 0, time.UTC)
	comparison := model.DailyPriceComparison{
		Symbol: "12340", TradeDate: date,
		ClosePrice:         decimal.RequireFromString("110.5"),
		PreviousTradeDate:  date.AddDate(0, 0, -3),
		PreviousClosePrice: decimal.NewFromInt(100),
	}
	invalid := comparison
	invalid.PreviousClosePrice = decimal.Zero
	tests := []struct {
		name           string
		repository     fakeComparisonRepository
		wantStatus     int
		wantBody       string
		wantListCalled bool
	}{
		{
			name:           "最新の保存日を使って日付と倍率を返す",
			repository:     fakeComparisonRepository{latestDate: date, latestFound: true, comparisons: []model.DailyPriceComparison{comparison}},
			wantStatus:     http.StatusOK,
			wantBody:       `{"trade_date":"2026-10-19","changes":[{"symbol":"12340","trade_date":"2026-10-19","close_price":"110.5","previous_trade_date":"2026-10-16","previous_close_price":"100.0","change_rate":"10.5"}]}`,
			wantListCalled: true,
		},
		{
			name:       "DBが空なら日付なしと空配列を返す",
			wantStatus: http.StatusOK,
			wantBody:   `{"trade_date":null,"changes":[]}`,
		},
		{
			name:           "比較元がない場合も最新日を維持する",
			repository:     fakeComparisonRepository{latestDate: date, latestFound: true},
			wantStatus:     http.StatusOK,
			wantBody:       `{"trade_date":"2026-10-19","changes":[]}`,
			wantListCalled: true,
		},
		{
			name:       "最新日の取得失敗は500",
			repository: fakeComparisonRepository{latestErr: errors.New("database error")},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"message":"failed to calculate latest daily price changes"}`,
		},
		{
			name:           "日足の取得失敗は500",
			repository:     fakeComparisonRepository{latestDate: date, latestFound: true, err: errors.New("database error")},
			wantStatus:     http.StatusInternalServerError,
			wantBody:       `{"message":"failed to calculate latest daily price changes"}`,
			wantListCalled: true,
		},
		{
			name:           "計算失敗は500",
			repository:     fakeComparisonRepository{latestDate: date, latestFound: true, comparisons: []model.DailyPriceComparison{invalid}},
			wantStatus:     http.StatusInternalServerError,
			wantBody:       `{"message":"failed to calculate latest daily price changes"}`,
			wantListCalled: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/v1/daily-price-changes/latest", nil)
			listCalled := false
			repository := test.repository
			repository.listByTradeDate = func(ctx context.Context, gotDate time.Time) {
				listCalled = true
				if !gotDate.Equal(date) {
					t.Errorf("date = %v, want latest stored date %v", gotDate, date)
				}
				if ctx != request.Context() {
					t.Error("request context was not passed to repository")
				}
			}
			calculator := dailyprice.NewDailyPriceChangeCalculator(repository)
			volumeCalculator := dailyprice.NewDailyVolumeRatioCalculator(fakeVolumeReader{})
			engine := NewRouter(handler.NewStockHandler(fakeQuerier{}), handler.NewDailyPriceChangeHandler(calculator), handler.NewDailyVolumeRatioHandler(volumeCalculator))
			recorder := httptest.NewRecorder()
			engine.ServeHTTP(recorder, request)
			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d: %s", recorder.Code, test.wantStatus, recorder.Body.String())
			}
			if listCalled != test.wantListCalled {
				t.Errorf("list called = %v, want %v", listCalled, test.wantListCalled)
			}
			var got, want any
			if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal([]byte(test.wantBody), &want); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("body = %s, want %s", recorder.Body.String(), test.wantBody)
			}
		})
	}
}
