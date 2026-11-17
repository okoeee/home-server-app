package fetch

import (
	"crypto/tls"
	"fmt"
	"net/http/httptrace"
	"strings"
	"sync"
	"time"
)

// 各イベントにはリクエスト開始からの経過時間を記録する。
// 接続試行は並行することがあるため、各区間の時間を合算しない。
type requestTrace struct {
	startedAt time.Time
	mu        sync.Mutex
	events    []string
}

func newRequestTrace() *requestTrace {
	return &requestTrace{startedAt: time.Now()}
}

func (trace *requestTrace) record(event string) {
	// httptraceのコールバックは、並行実行やリクエスト終了後の実行があり得る。
	trace.mu.Lock()
	defer trace.mu.Unlock()
	trace.events = append(trace.events, fmt.Sprintf("%s@%dms", event, time.Since(trace.startedAt).Milliseconds()))
}

func (trace *requestTrace) recordResult(event string, err error) {
	if err != nil {
		trace.record(fmt.Sprintf("%s(error=%q)", event, err.Error()))
		return
	}
	trace.record(event + "(ok)")
}

func (trace *requestTrace) String() string {
	trace.mu.Lock()
	defer trace.mu.Unlock()
	return fmt.Sprintf("total=%dms events=[%s]", time.Since(trace.startedAt).Milliseconds(), strings.Join(trace.events, ", "))
}

func (trace *requestTrace) hooks() *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			trace.record("get_conn(" + hostPort + ")")
		},
		DNSStart: func(info httptrace.DNSStartInfo) {
			trace.record("dns_start(" + info.Host + ")")
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			trace.recordResult("dns_done", info.Err)
		},
		ConnectStart: func(network, address string) {
			trace.record("connect_start(" + network + " " + address + ")")
		},
		ConnectDone: func(network, address string, err error) {
			trace.recordResult("connect_done("+network+" "+address+")", err)
		},
		TLSHandshakeStart: func() {
			trace.record("tls_start")
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			trace.recordResult("tls_done("+tls.VersionName(state.Version)+")", err)
		},
		GotConn: func(info httptrace.GotConnInfo) {
			trace.record(fmt.Sprintf("got_conn(reused=%t)", info.Reused))
		},
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			trace.recordResult("request_written", info.Err)
		},
		GotFirstResponseByte: func() {
			trace.record("first_response_byte")
		},
	}
}
