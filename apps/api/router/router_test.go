package router

import (
	"api/db"
	"api/domain/model"
	"api/handler"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type stockQuerierStub struct{}

func (stockQuerierStub) ListDailyPricesBySymbol(context.Context, string) ([]db.DailyPrice, error) {
	return nil, nil
}

type priceChangeCalculatorStub struct {
	changesByDate map[time.Time][]model.DailyPriceChange
	latest        model.LatestDailyPriceChanges
	err           error
}

func (stub priceChangeCalculatorStub) Calculate(_ context.Context, date time.Time) ([]model.DailyPriceChange, error) {
	return stub.changesByDate[date], stub.err
}

func (stub priceChangeCalculatorStub) CalculateLatest(context.Context) (model.LatestDailyPriceChanges, error) {
	return stub.latest, stub.err
}

type volumeRatioCalculatorStub struct {
	ratiosByDate map[time.Time][]model.DailyVolumeRatio
	err          error
}

func (stub volumeRatioCalculatorStub) Calculate(_ context.Context, date time.Time) ([]model.DailyVolumeRatio, error) {
	return stub.ratiosByDate[date], stub.err
}

func newTestRouter(priceCalculator handler.DailyPriceChangeCalculator, volumeCalculator handler.DailyVolumeRatioCalculator) *gin.Engine {
	return NewRouter(
		handler.NewStockHandler(stockQuerierStub{}),
		handler.NewDailyPriceChangeHandler(priceCalculator),
		handler.NewDailyVolumeRatioHandler(volumeCalculator),
	)
}

func assertJSONResponse(t *testing.T, recorder *httptest.ResponseRecorder, wantStatus int, wantBody string) {
	t.Helper()
	if recorder.Code != wantStatus {
		t.Fatalf("status = %d, want %d: %s", recorder.Code, wantStatus, recorder.Body.String())
	}
	var got, want any
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	if err := json.Unmarshal([]byte(wantBody), &want); err != nil {
		t.Fatalf("invalid expected JSON: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("body = %s, want %s", recorder.Body.String(), wantBody)
	}
}

func TestPing_Success(t *testing.T) {
	// Arrange
	engine := newTestRouter(priceChangeCalculatorStub{}, volumeRatioCalculatorStub{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/ping", nil)

	// Act
	engine.ServeHTTP(recorder, request)

	// Assert
	assertJSONResponse(t, recorder, http.StatusOK, `{"message":"pong"}`)
}

func TestNewRouter_VersionedRoutes(t *testing.T) {
	tests := []struct {
		path       string
		wantStatus int
	}{
		{path: "/v1/stocks/12340/prices", wantStatus: http.StatusOK},
		{path: "/daily-price-changes?date=2026-10-19", wantStatus: http.StatusNotFound},
		{path: "/daily-price-changes/latest", wantStatus: http.StatusNotFound},
		{path: "/daily-volume-ratios?date=2026-10-19", wantStatus: http.StatusNotFound},
		{path: "/stocks/12340/prices", wantStatus: http.StatusNotFound},
	}
	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			// Arrange
			engine := newTestRouter(priceChangeCalculatorStub{}, volumeRatioCalculatorStub{})
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, test.path, nil)

			// Act
			engine.ServeHTTP(recorder, request)

			// Assert
			if recorder.Code != test.wantStatus {
				t.Errorf("status = %d, want %d", recorder.Code, test.wantStatus)
			}
		})
	}
}
