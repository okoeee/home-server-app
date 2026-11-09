package prometheus

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

type trackingBody struct {
	io.Reader
	closed bool
}

func (body *trackingBody) Close() error {
	body.closed = true
	return nil
}

func TestQuery_RequestAndResponse(t *testing.T) {
	// Arrange
	at := time.Date(2026, 11, 7, 0, 0, 0, 123000000, time.UTC)
	body := &trackingBody{Reader: strings.NewReader(`{"status":"success","data":{"resultType":"vector","result":[{"metric":{},"value":[1794009600.123,"12.5"]}]}}`)}
	client := NewClient("http://prometheus:9090/")
	client.httpClient.Transport = roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Method != http.MethodGet || request.URL.Host != "prometheus:9090" || request.URL.Path != "/api/v1/query" {
			t.Errorf("unexpected request: %s %s", request.Method, request.URL)
		}
		if request.URL.Query().Get("query") != cpuUsageQuery || request.URL.Query().Get("time") != at.Format(time.RFC3339Nano) {
			t.Errorf("unexpected query: %s", request.URL.RawQuery)
		}
		return &http.Response{StatusCode: http.StatusOK, Body: body}, nil
	})

	// Act
	result, err := client.query(context.Background(), cpuUsageQuery, at)

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 || result[0] != 12.5 {
		t.Errorf("result = %v, want one value of 12.5", result)
	}
	if !body.closed {
		t.Error("response body was not closed")
	}
}

func TestQuery_ResponseErrors(t *testing.T) {
	for _, test := range []struct {
		name   string
		status int
		body   string
	}{
		{"HTTPError", 503, `unavailable`},
		{"InvalidJSON", 200, `{`},
		{"TrailingData", 200, `{"status":"success","data":{"resultType":"vector","result":[]}}{}`},
		{"QueryError", 200, `{"status":"error","error":"bad query"}`},
		{"MissingResult", 200, `{"status":"success","data":{"resultType":"vector"}}`},
		{"WrongResultType", 200, `{"status":"success","data":{"resultType":"scalar","result":[1,"1"]}}`},
		{"MissingValue", 200, `{"status":"success","data":{"resultType":"vector","result":[{}]}}`},
		{"NullTimestamp", 200, `{"status":"success","data":{"resultType":"vector","result":[{"value":[null,"1"]}]}}`},
		{"NullValue", 200, `{"status":"success","data":{"resultType":"vector","result":[{"value":[1,null]}]}}`},
		{"NaN", 200, `{"status":"success","data":{"resultType":"vector","result":[{"value":[1,"NaN"]}]}}`},
		{"Infinity", 200, `{"status":"success","data":{"resultType":"vector","result":[{"value":[1,"+Inf"]}]}}`},
		{"InvalidNumber", 200, `{"status":"success","data":{"resultType":"vector","result":[{"value":[1,"invalid"]}]}}`},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			body := &trackingBody{Reader: strings.NewReader(test.body)}
			client := NewClient("http://prometheus:9090")
			client.httpClient.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: test.status, Body: body}, nil
			})

			// Act
			_, err := client.query(context.Background(), upQuery, time.Now())

			// Assert
			if !errors.Is(err, ErrQueryFailed) {
				t.Errorf("error = %v, want ErrQueryFailed", err)
			}
			if !body.closed {
				t.Error("response body was not closed")
			}
		})
	}
}

func TestQuery_CancelledRequest(t *testing.T) {
	// Arrange
	client := NewClient("http://prometheus:9090")
	client.httpClient.Transport = roundTripFunc(func(request *http.Request) (*http.Response, error) {
		return nil, request.Context().Err()
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	_, err := client.query(ctx, upQuery, time.Now())

	// Assert
	if !errors.Is(err, ErrQueryFailed) || !errors.Is(err, context.Canceled) {
		t.Errorf("error = %v, want wrapped context.Canceled and ErrQueryFailed", err)
	}
}
