package fetch

import (
	"net/http/httptrace"
	"strings"
	"sync"
	"testing"
)

func TestRequestTrace_ConcurrentCallbacksAndSnapshots(t *testing.T) {
	// Arrange
	trace := newRequestTrace()
	hooks := trace.hooks()
	var workers sync.WaitGroup
	const count = 32

	// Act
	for range count {
		workers.Go(func() {
			hooks.WroteRequest(httptrace.WroteRequestInfo{})
			_ = trace.String()
		})
	}
	workers.Wait()
	got := trace.String()

	// Assert
	if events := strings.Count(got, "request_written(ok)@"); events != count {
		t.Errorf("recorded events = %d, want %d: %s", events, count, got)
	}
}
