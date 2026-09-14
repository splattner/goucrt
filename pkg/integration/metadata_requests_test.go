package integration

// Round-trip tests for the driver-initiated metadata requests (metadata_requests.go): GetVersion
// and friends send a "req" frame and block waiting for a correlated "resp" frame delivered through
// handleResponse - the first thing in this package that ever both sends a request AND consumes its
// own response, rather than either only responding (handleRequest) or only firing events
// (sendEventMessage). These tests drive that whole path without a real WebSocket: read the request
// GetVersion put on Remote.messageChannel, build the matching response frame by hand, and feed it
// back through handleResponse exactly as wsReader would for a "resp" kind message.

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// awaitOneRequest returns a channel that receives exactly one outbound frame sent on
// i.Remote.messageChannel, decoded far enough to read its req id - enough to build a matching
// response. Mirrors awaitOneMessage in handlers_test.go, which does the same for a synchronous
// handleRequest call's response.
func awaitOneRequest(t *testing.T, i *Integration) <-chan RequestMessage {
	t.Helper()
	out := make(chan RequestMessage, 1)
	go func() {
		raw := <-i.Remote.messageChannel
		var req RequestMessage
		if err := json.Unmarshal(raw, &req); err != nil {
			t.Errorf("decode outbound request: %v", err)
			return
		}
		out <- req
	}()
	return out
}

func TestGetVersion_RoundTrip(t *testing.T) {
	i := newTestIntegration(t)

	reqCh := awaitOneRequest(t, i)
	resultCh := make(chan *VersionInfo, 1)
	errCh := make(chan error, 1)
	go func() {
		v, err := i.GetVersion()
		resultCh <- v
		errCh <- err
	}()

	req := <-reqCh
	if req.Msg != "get_version" {
		t.Fatalf("outbound request msg = %q, want %q", req.Msg, "get_version")
	}
	if req.Kind != "req" {
		t.Errorf("outbound request kind = %q, want %q", req.Kind, "req")
	}

	resp, err := json.Marshal(VersionMessage{
		CommonResp: CommonResp{Kind: "resp", Id: req.Id, Msg: "version", Code: 200},
		MsgData:    VersionInfo{Model: "UCR3", DeviceName: "Living Room", Api: "1.0", Core: "2.0"},
	})
	if err != nil {
		t.Fatalf("marshal response: %v", err)
	}
	i.handleResponse(resp)

	if err := <-errCh; err != nil {
		t.Fatalf("GetVersion() error = %v", err)
	}
	got := <-resultCh
	if got == nil || got.Model != "UCR3" || got.DeviceName != "Living Room" {
		t.Errorf("GetVersion() = %+v, want Model=UCR3 DeviceName=\"Living Room\"", got)
	}
}

func TestGetSupportedEntityTypes_RoundTrip(t *testing.T) {
	i := newTestIntegration(t)

	reqCh := awaitOneRequest(t, i)
	resultCh := make(chan []string, 1)
	errCh := make(chan error, 1)
	go func() {
		types, err := i.GetSupportedEntityTypes()
		resultCh <- types
		errCh <- err
	}()

	req := <-reqCh
	if req.Msg != "get_supported_entity_types" {
		t.Fatalf("outbound request msg = %q, want %q", req.Msg, "get_supported_entity_types")
	}

	resp, err := json.Marshal(SupportedEntityTypesMessage{
		CommonResp: CommonResp{Kind: "resp", Id: req.Id, Msg: "supported_entity_types", Code: 200},
		MsgData:    []string{"light", "switch", "media_player"},
	})
	if err != nil {
		t.Fatalf("marshal response: %v", err)
	}
	i.handleResponse(resp)

	if err := <-errCh; err != nil {
		t.Fatalf("GetSupportedEntityTypes() error = %v", err)
	}
	got := <-resultCh
	if len(got) != 3 || got[0] != "light" {
		t.Errorf("GetSupportedEntityTypes() = %v, want [light switch media_player]", got)
	}
}

// TestMetadataRequests_ConcurrentCallsGetTheirOwnResponse is the regression test for the
// pendingRequests correlation itself: fires several GetRuntimeInfo-style requests concurrently and
// replies to them out of order, keyed only by the req_id each request actually got assigned.
// Routing by array position or a single shared channel (rather than a map keyed by id) would
// silently hand one caller another caller's response; run with -race, a missing requestsMu lock
// around the map would also be flagged here.
func TestMetadataRequests_ConcurrentCallsGetTheirOwnResponse(t *testing.T) {
	i := newTestIntegration(t)

	const n = 10
	type result struct {
		info *RuntimeInfo
		err  error
	}
	resultCh := make(chan result, n)

	for k := 0; k < n; k++ {
		go func() {
			info, err := i.GetRuntimeInfo()
			resultCh <- result{info, err}
		}()
	}

	// Drain n outbound requests and reply to each with a driver_id derived from its own req_id, out
	// of arrival order, so a mis-routed response would produce duplicate or missing driver_id
	// values below instead of coincidentally passing.
	go func() {
		pending := make([]RequestMessage, 0, n)
		for len(pending) < n {
			raw := <-i.Remote.messageChannel
			var req RequestMessage
			if err := json.Unmarshal(raw, &req); err != nil {
				t.Errorf("decode outbound request: %v", err)
				return
			}
			pending = append(pending, req)
		}
		// Reply in reverse order of arrival.
		for k := len(pending) - 1; k >= 0; k-- {
			req := pending[k]
			resp, err := json.Marshal(RuntimeInfoMessage{
				CommonResp: CommonResp{Kind: "resp", Id: req.Id, Msg: "runtime_info", Code: 200},
				MsgData:    RuntimeInfo{DriverId: fmt.Sprintf("driver-%d", req.Id)},
			})
			if err != nil {
				t.Errorf("marshal response: %v", err)
				return
			}
			i.handleResponse(resp)
		}
	}()

	// Each request's driver_id is unique (derived from its own req_id), so every distinct value
	// seen across all n calls, with none repeated, is only possible if each call received the
	// response for exactly the request it sent - a response delivered to the wrong waiter (or two
	// waiters sharing one channel) would surface here as a duplicate and a missing value.
	seen := make(map[string]bool, n)
	for k := 0; k < n; k++ {
		r := <-resultCh
		if r.err != nil {
			t.Errorf("GetRuntimeInfo() error = %v", r.err)
			continue
		}
		if seen[r.info.DriverId] {
			t.Errorf("DriverId %q was returned to more than one caller", r.info.DriverId)
		}
		seen[r.info.DriverId] = true
	}
	if len(seen) != n {
		t.Errorf("got %d distinct DriverId values across %d concurrent calls, want %d", len(seen), n, n)
	}
}

func TestSendMetadataRequest_TimesOutWithoutResponse(t *testing.T) {
	i := newTestIntegration(t)

	orig := requestTimeout
	requestTimeout = 20 * time.Millisecond
	t.Cleanup(func() { requestTimeout = orig })

	// Drain the outbound request but never reply, so the wait for a response times out.
	go func() { <-i.Remote.messageChannel }()

	if _, err := i.GetVersion(); err == nil {
		t.Error("GetVersion() with no response = nil error, want a timeout error")
	}

	// The abandoned request must not linger in pendingRequests after it times out - otherwise a
	// late, spurious "resp" frame reusing that id would be delivered to nothing (fine) but the
	// entry would leak for the life of the process.
	i.requestsMu.Lock()
	n := len(i.pendingRequests)
	i.requestsMu.Unlock()
	if n != 0 {
		t.Errorf("pendingRequests has %d entries after timeout, want 0", n)
	}
}

func TestSendMetadataRequest_ErrorCodeIsReturnedAsError(t *testing.T) {
	i := newTestIntegration(t)

	reqCh := awaitOneRequest(t, i)
	errCh := make(chan error, 1)
	go func() {
		_, err := i.GetVersion()
		errCh <- err
	}()

	req := <-reqCh
	resp, err := json.Marshal(ResponseMessage{
		CommonResp: CommonResp{Kind: "resp", Id: req.Id, Msg: "result", Code: 501},
		MsgData:    ErrorData{Code: "OTHER", Message: "not supported"},
	})
	if err != nil {
		t.Fatalf("marshal response: %v", err)
	}
	i.handleResponse(resp)

	if err := <-errCh; err == nil {
		t.Error("GetVersion() with a non-200 response code = nil error, want an error")
	}
}

func TestHandleResponse_UnknownReqIdIsDroppedNotPanicked(t *testing.T) {
	i := newTestIntegration(t)

	resp, err := json.Marshal(ResponseMessage{
		CommonResp: CommonResp{Kind: "resp", Id: 999, Msg: "version", Code: 200},
	})
	if err != nil {
		t.Fatalf("marshal response: %v", err)
	}

	// Must return promptly rather than blocking, and must not panic.
	done := make(chan struct{})
	go func() {
		i.handleResponse(resp)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("handleResponse blocked on a response with no matching pending request")
	}
}
