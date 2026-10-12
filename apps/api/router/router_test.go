package router

import (
	"api/db"
	"api/handler"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeQuerier struct{}

func (fakeQuerier) ListDailyPricesBySymbol(context.Context, string) ([]db.DailyPrice, error) {
	return nil, nil
}

func TestPing(t *testing.T) {

	stockHandler := handler.NewStockHandler(fakeQuerier{})
	engine := NewRouter(stockHandler)

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
