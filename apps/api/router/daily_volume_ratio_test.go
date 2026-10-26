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

type fakeVolumeReader struct {
	volumes      []model.DailyVolume
	history      []model.DailyVolume
	volumesErr   error
	historyErr   error
	checkRequest func(context.Context, time.Time)
}

func (reader fakeVolumeReader) ListVolumesByTradeDate(ctx context.Context, date time.Time) ([]model.DailyVolume, error) {
	if reader.checkRequest != nil {
		reader.checkRequest(ctx, date)
	}
	return reader.volumes, reader.volumesErr
}

func (reader fakeVolumeReader) ListVolumeHistory(ctx context.Context, date time.Time) ([]model.DailyVolume, error) {
	if reader.checkRequest != nil {
		reader.checkRequest(ctx, date)
	}
	return reader.history, reader.historyErr
}

func TestGetDailyVolumeRatios(t *testing.T) {
	tradeDate := time.Date(2026, 10, 19, 0, 0, 0, 0, time.UTC)
	volumes := []model.DailyVolume{
		{Symbol: "12340", TradeDate: tradeDate, Volume: decimal.RequireFromString("150.3")},
	}
	history := make([]model.DailyVolume, 0, 20)
	for index := range 20 {
		history = append(history, model.DailyVolume{
			Symbol:    "12340",
			TradeDate: tradeDate.AddDate(0, 0, -index-1),
			Volume:    decimal.RequireFromString("100.2"),
		})
	}
	tests := []struct {
		name       string
		query      string
		reader     fakeVolumeReader
		wantStatus int
		wantBody   string
		wantCalls  int
	}{
		{
			name:       "出来高と倍率を文字列で返す",
			query:      "?date=2026-10-19",
			reader:     fakeVolumeReader{volumes: volumes, history: history},
			wantStatus: http.StatusOK,
			wantBody:   `[{"symbol":"12340","trade_date":"2026-10-19","volume":"150.3","ratio":"1.5"}]`,
			wantCalls:  2,
		},
		{
			name:       "対象データなしは空配列を返す",
			query:      "?date=2026-10-19",
			wantStatus: http.StatusOK,
			wantBody:   `[]`,
			wantCalls:  2,
		},
		{
			name:       "計算できる銘柄がなければ空配列を返す",
			query:      "?date=2026-10-19",
			reader:     fakeVolumeReader{volumes: volumes},
			wantStatus: http.StatusOK,
			wantBody:   `[]`,
			wantCalls:  2,
		},
		{name: "日付未指定は400", wantStatus: http.StatusBadRequest, wantBody: `{"message":"date must be in YYYY-MM-DD format"}`},
		{name: "不正な日付形式は400", query: "?date=2026/10/19", wantStatus: http.StatusBadRequest, wantBody: `{"message":"date must be in YYYY-MM-DD format"}`},
		{name: "存在しない日付は400", query: "?date=2026-02-30", wantStatus: http.StatusBadRequest, wantBody: `{"message":"date must be in YYYY-MM-DD format"}`},
		{
			name:       "対象日の取得失敗は500",
			query:      "?date=2026-10-19",
			reader:     fakeVolumeReader{volumesErr: errors.New("database error")},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"message":"failed to calculate daily volume ratios"}`,
			wantCalls:  1,
		},
		{
			name:       "履歴の取得失敗は500",
			query:      "?date=2026-10-19",
			reader:     fakeVolumeReader{historyErr: errors.New("database error")},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"message":"failed to calculate daily volume ratios"}`,
			wantCalls:  2,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/daily-volume-ratios"+test.query, nil)
			calls := 0
			reader := test.reader
			reader.checkRequest = func(ctx context.Context, date time.Time) {
				calls++
				if !date.Equal(tradeDate) {
					t.Errorf("date = %v, want %v", date, tradeDate)
				}
				if ctx != request.Context() {
					t.Error("request context was not passed to reader")
				}
			}
			priceCalculator := dailyprice.NewDailyPriceChangeCalculator(fakeComparisonRepository{})
			volumeCalculator := dailyprice.NewDailyVolumeRatioCalculator(reader)
			engine := NewRouter(
				handler.NewStockHandler(fakeQuerier{}),
				handler.NewDailyPriceChangeHandler(priceCalculator),
				handler.NewDailyVolumeRatioHandler(volumeCalculator),
			)
			recorder := httptest.NewRecorder()
			engine.ServeHTTP(recorder, request)
			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d: %s", recorder.Code, test.wantStatus, recorder.Body.String())
			}
			if calls != test.wantCalls {
				t.Errorf("reader calls = %d, want %d", calls, test.wantCalls)
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
