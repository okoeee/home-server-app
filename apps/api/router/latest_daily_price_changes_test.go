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

func TestGetLatestDailyPriceChanges_Scenarios(t *testing.T) {
	date := time.Date(2026, 10, 19, 0, 0, 0, 0, time.UTC)
	companyName := "サンプル株式会社"
	changes := []model.DailyPriceChange{{
		Symbol: "12340", TradeDate: date, ClosePrice: decimal.RequireFromString("110.5"),
		CompanyName:       &companyName,
		PreviousTradeDate: date.AddDate(0, 0, -3), PreviousClosePrice: decimal.NewFromInt(100),
		ChangeRate: decimal.RequireFromString("10.5"),
	}}
	missingName := changes[0]
	missingName.CompanyName = nil
	tests := []struct {
		name       string
		result     model.LatestDailyPriceChanges
		err        error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Usecaseの対象日と結果をJSONで返す",
			result:     model.LatestDailyPriceChanges{TradeDate: &date, Changes: changes},
			wantStatus: http.StatusOK,
			wantBody:   `{"trade_date":"2026-10-19","changes":[{"symbol":"12340","company_name":"サンプル株式会社","trade_date":"2026-10-19","close_price":"110.5","previous_trade_date":"2026-10-16","previous_close_price":"100.0","change_rate":"10.5"}]}`,
		},
		{
			name:       "会社名がなくても銘柄を残しnullを返す",
			result:     model.LatestDailyPriceChanges{TradeDate: &date, Changes: []model.DailyPriceChange{missingName}},
			wantStatus: http.StatusOK,
			wantBody:   `{"trade_date":"2026-10-19","changes":[{"symbol":"12340","company_name":null,"trade_date":"2026-10-19","close_price":"110.5","previous_trade_date":"2026-10-16","previous_close_price":"100.0","change_rate":"10.5"}]}`,
		},
		{
			name:       "対象日なしはnullと空配列を返す",
			wantStatus: http.StatusOK,
			wantBody:   `{"trade_date":null,"changes":[]}`,
		},
		{
			name:       "対象日があれば空結果でも日付を返す",
			result:     model.LatestDailyPriceChanges{TradeDate: &date},
			wantStatus: http.StatusOK,
			wantBody:   `{"trade_date":"2026-10-19","changes":[]}`,
		},
		{
			name:       "Usecase失敗は500",
			err:        errors.New("calculation failed"),
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"message":"failed to calculate latest daily price changes"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			calculator := priceChangeCalculatorStub{latest: test.result, err: test.err}
			engine := newTestRouter(calculator, volumeRatioCalculatorStub{})
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/v1/daily-price-changes/latest", nil)

			// Act
			engine.ServeHTTP(recorder, request)

			// Assert
			assertJSONResponse(t, recorder, test.wantStatus, test.wantBody)
		})
	}
}
