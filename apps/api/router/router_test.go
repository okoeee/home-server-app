package router

import (
	"api/db"
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

type fakeQuerier struct{}

type fakeComparisonRepository struct {
	comparisons     []model.DailyPriceComparison
	err             error
	listByTradeDate func(context.Context, time.Time)
}

func (repository fakeComparisonRepository) ListByTradeDate(ctx context.Context, date time.Time) ([]model.DailyPriceComparison, error) {
	if repository.listByTradeDate != nil {
		repository.listByTradeDate(ctx, date)
	}
	return repository.comparisons, repository.err
}

func (fakeQuerier) ListDailyPricesBySymbol(context.Context, string) ([]db.DailyPrice, error) {
	return nil, nil
}

func TestGetDailyPriceChanges(t *testing.T) {
	tradeDate := time.Date(2026, 10, 19, 0, 0, 0, 0, time.UTC)
	comparison := model.DailyPriceComparison{
		Symbol:             "12340",
		TradeDate:          tradeDate,
		ClosePrice:         decimal.RequireFromString("110.5"),
		PreviousTradeDate:  tradeDate.AddDate(0, 0, -3),
		PreviousClosePrice: decimal.NewFromInt(100),
	}
	zeroPreviousClose := comparison
	zeroPreviousClose.PreviousClosePrice = decimal.Zero

	tests := []struct {
		name        string
		query       string
		comparisons []model.DailyPriceComparison
		err         error
		wantStatus  int
		wantBody    string
		wantCalled  bool
	}{
		{
			name:        "騰落率と比較元の日付を返す",
			query:       "?date=2026-10-19",
			comparisons: []model.DailyPriceComparison{comparison},
			wantStatus:  http.StatusOK,
			wantBody:    `[{"symbol":"12340","trade_date":"2026-10-19","close_price":"110.5","previous_trade_date":"2026-10-16","previous_close_price":"100.0","change_rate":"10.5"}]`,
			wantCalled:  true,
		},
		{
			name:       "データがなければ空配列を返す",
			query:      "?date=2026-10-19",
			wantStatus: http.StatusOK,
			wantBody:   `[]`,
			wantCalled: true,
		},
		{
			name:       "日付未指定は400を返す",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"message":"date must be in YYYY-MM-DD format"}`,
		},
		{
			name:       "日付形式が不正なら400を返す",
			query:      "?date=2026/10/19",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"message":"date must be in YYYY-MM-DD format"}`,
		},
		{
			name:       "存在しない日付は400を返す",
			query:      "?date=2026-02-30",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"message":"date must be in YYYY-MM-DD format"}`,
		},
		{
			name:       "取得失敗は500を返す",
			query:      "?date=2026-10-19",
			err:        errors.New("database error"),
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"message":"failed to calculate daily price changes"}`,
			wantCalled: true,
		},
		{
			name:        "計算失敗は500を返す",
			query:       "?date=2026-10-19",
			comparisons: []model.DailyPriceComparison{zeroPreviousClose},
			wantStatus:  http.StatusInternalServerError,
			wantBody:    `{"message":"failed to calculate daily price changes"}`,
			wantCalled:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/daily-price-changes"+test.query, nil)
			called := false
			repository := fakeComparisonRepository{
				comparisons: test.comparisons,
				err:         test.err,
				listByTradeDate: func(ctx context.Context, date time.Time) {
					called = true
					if !date.Equal(tradeDate) {
						t.Errorf("date = %v, want %v", date, tradeDate)
					}
					if ctx != request.Context() {
						t.Error("request context was not passed to repository")
					}
				},
			}
			calculator := dailyprice.NewDailyPriceChangeCalculator(repository)
			engine := NewRouter(handler.NewStockHandler(fakeQuerier{}), handler.NewDailyPriceChangeHandler(calculator))
			recorder := httptest.NewRecorder()
			engine.ServeHTTP(recorder, request)

			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d: %s", recorder.Code, test.wantStatus, recorder.Body.String())
			}
			if called != test.wantCalled {
				t.Errorf("repository called = %v, want %v", called, test.wantCalled)
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

func (fakeQuerier) UpsertDailyPrice(context.Context, db.UpsertDailyPriceParams) error {
	return nil
}

func (fakeQuerier) UpsertStock(context.Context, db.UpsertStockParams) error {
	return nil
}

func TestPing(t *testing.T) {

	stockHandler := handler.NewStockHandler(fakeQuerier{})
	calculator := dailyprice.NewDailyPriceChangeCalculator(fakeComparisonRepository{})
	changeHandler := handler.NewDailyPriceChangeHandler(calculator)
	engine := NewRouter(stockHandler, changeHandler)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/ping", nil)

	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fail()
	}

	want := `{"message":"pong"}`

	if got := recorder.Body.String(); got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}

}
