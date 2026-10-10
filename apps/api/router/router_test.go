package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {

	engine := NewRouter()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/ping", nil)

	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fail()
	}

	want := `{"message":"pong"}`

	fmt.Print(recorder.Body.String())

	if got := recorder.Body.String(); got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}

}
