package integration

// Handler-level regression tests that drive handleRequest with real request frames, the way the
// remote's WebSocket connection does, rather than calling the handleXRequest methods directly.
// TestHandleSubscribeEventRequest_CallsSubscribeCallback and its unsubscribe counterpart exist
// specifically to guard the bug a shadow/errcheck-style lint pass can't catch: handleSubscribeEventRequest
// and handleUnsubscribeEventsRequest used to call GetEntityById and only invoke the subscribe/unsubscribe
// callback when it returned an error - i.e. only when the entity was *not* found - so the callback
// silently never ran for a real subscription. Fixed by inspection; these tests are what stops it from
// coming back unnoticed.

import (
	"encoding/json"
	"testing"

	"github.com/splattner/goucrt/pkg/entities"
)

func newTestIntegration(t *testing.T) *Integration {
	t.Helper()
	i, err := NewIntegration(Config{ListenPort: 0, WebsocketPath: "/ws", ConfigHome: t.TempDir() + "/"})
	if err != nil {
		t.Fatalf("NewIntegration: %v", err)
	}
	// deviceState defaults to DisconnectedDeviceState, so sendEntityAvailable/sendDeviceStateEvent
	// are no-ops and won't try to write to Remote.messageChannel underneath us.
	return i
}

// awaitOneMessage returns a channel that receives exactly one message sent on
// i.Remote.messageChannel, so a synchronous handleRequest call (which sends its response on that
// same unbuffered channel) doesn't deadlock the test.
func awaitOneMessage(i *Integration) <-chan []byte {
	out := make(chan []byte, 1)
	go func() { out <- <-i.Remote.messageChannel }()
	return out
}

func decodeResponse(t *testing.T, raw []byte) ResponseMessage {
	t.Helper()
	var res ResponseMessage
	if err := json.Unmarshal(raw, &res); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return res
}

func TestHandleSubscribeEventRequest_CallsSubscribeCallback(t *testing.T) {
	i := newTestIntegration(t)

	e := entities.NewSwitchEntity("switch-1", entities.LanguageText{En: "Switch"}, "")
	subscribed := false
	e.SetSubscribeCallbackFunc(func() { subscribed = true })
	i.Entities = append(i.Entities, e)

	raw, err := json.Marshal(SubscribeEventMessageReq{
		CommonReq: CommonReq{Kind: "req", Id: 1, Msg: "subscribe_events"},
		MsgData:   SubscribeEventMessageData{EntityIds: []string{"switch-1"}},
	})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	resp := awaitOneMessage(i)
	i.handleRequest(&RequestMessage{CommonReq: CommonReq{Kind: "req", Id: 1, Msg: "subscribe_events"}}, raw)

	if res := decodeResponse(t, <-resp); res.Code != 200 {
		t.Errorf("response code = %d, want 200", res.Code)
	}
	if !subscribed {
		t.Error("SubscribeCallbackFunc was not called for a named entity subscription")
	}
	if !containsString(i.SubscribedEntities, "switch-1") {
		t.Errorf("SubscribedEntities = %v, want it to contain \"switch-1\"", i.SubscribedEntities)
	}
}

func TestHandleUnsubscribeEventsRequest_CallsUnsubscribeCallback(t *testing.T) {
	i := newTestIntegration(t)

	e := entities.NewSwitchEntity("switch-1", entities.LanguageText{En: "Switch"}, "")
	unsubscribed := false
	e.SetUnsubscribeCallbackFunc(func() { unsubscribed = true })
	i.Entities = append(i.Entities, e)
	i.SubscribedEntities = append(i.SubscribedEntities, "switch-1")

	raw, err := json.Marshal(UnubscribeEventMessageReq{
		CommonReq: CommonReq{Kind: "req", Id: 2, Msg: "unsubscribe_events"},
		MsgData:   SubscribeEventMessageData{EntityIds: []string{"switch-1"}},
	})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	resp := awaitOneMessage(i)
	i.handleRequest(&RequestMessage{CommonReq: CommonReq{Kind: "req", Id: 2, Msg: "unsubscribe_events"}}, raw)

	if res := decodeResponse(t, <-resp); res.Code != 200 {
		t.Errorf("response code = %d, want 200", res.Code)
	}
	if !unsubscribed {
		t.Error("UnsubscribeCallbackFunc was not called for a named entity unsubscription")
	}
	if containsString(i.SubscribedEntities, "switch-1") {
		t.Errorf("SubscribedEntities = %v, want it to no longer contain \"switch-1\"", i.SubscribedEntities)
	}
}

func TestHandleEntityCommandRequest_DispatchesToRegisteredCommand(t *testing.T) {
	i := newTestIntegration(t)

	e := entities.NewSwitchEntity("switch-1", entities.LanguageText{En: "Switch"}, "")
	turnedOn := false
	e.MapCommand(entities.OnSwitchEntityCommand, func() error {
		turnedOn = true
		return nil
	})
	i.Entities = append(i.Entities, e)

	raw, err := json.Marshal(EntityCommandReq{
		CommonReq: CommonReq{Kind: "req", Id: 3, Msg: "entity_command"},
		MsgData: EntityCommandData{
			EntityId: "switch-1",
			CmdId:    "on",
		},
	})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	resp := awaitOneMessage(i)
	i.handleRequest(&RequestMessage{CommonReq: CommonReq{Kind: "req", Id: 3, Msg: "entity_command"}}, raw)

	if res := decodeResponse(t, <-resp); res.Code != 200 {
		t.Errorf("response code = %d, want 200", res.Code)
	}
	if !turnedOn {
		t.Error("the registered \"on\" command handler was not invoked")
	}
}

func TestHandleEntityCommandRequest_UnknownEntityReturns404(t *testing.T) {
	i := newTestIntegration(t)

	raw, err := json.Marshal(EntityCommandReq{
		CommonReq: CommonReq{Kind: "req", Id: 4, Msg: "entity_command"},
		MsgData: EntityCommandData{
			EntityId: "does-not-exist",
			CmdId:    "on",
		},
	})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	resp := awaitOneMessage(i)
	i.handleRequest(&RequestMessage{CommonReq: CommonReq{Kind: "req", Id: 4, Msg: "entity_command"}}, raw)

	if res := decodeResponse(t, <-resp); res.Code != 404 {
		t.Errorf("response code = %d, want 404 for an unknown entity_id", res.Code)
	}
}

func containsString(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return false
}

func TestHandleRequest_UnknownMessageReturnsNotFound(t *testing.T) {
	i := newTestIntegration(t)

	raw, err := json.Marshal(RequestMessage{
		CommonReq: CommonReq{Kind: "req", Id: 5, Msg: "not_a_real_message"},
	})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	resp := awaitOneMessage(i)
	i.handleRequest(&RequestMessage{CommonReq: CommonReq{Kind: "req", Id: 5, Msg: "not_a_real_message"}}, raw)

	res := decodeResponse(t, <-resp)
	if res.Code != 404 {
		t.Errorf("response code = %d, want 404", res.Code)
	}
	if res.Msg != "result" {
		t.Errorf("response msg = %q, want \"result\"", res.Msg)
	}

	var errData ErrorData
	msgDataRaw, _ := json.Marshal(res.MsgData)
	if err := json.Unmarshal(msgDataRaw, &errData); err != nil {
		t.Fatalf("decode msg_data as ErrorData: %v", err)
	}
	if errData.Code != "NOT_FOUND" {
		t.Errorf("msg_data.code = %q, want \"NOT_FOUND\"", errData.Code)
	}
}

func TestHandleRequest_AuthReturnsSuccess(t *testing.T) {
	i := newTestIntegration(t)

	raw, err := json.Marshal(AuthRequestMessage{
		CommonReq: CommonReq{Kind: "req", Id: 6, Msg: "auth"},
		MsgData:   AuthRequestData{Token: "anything"},
	})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	resp := awaitOneMessage(i)
	i.handleRequest(&RequestMessage{CommonReq: CommonReq{Kind: "req", Id: 6, Msg: "auth"}}, raw)

	res := decodeResponse(t, <-resp)
	if res.Code != 200 {
		t.Errorf("response code = %d, want 200", res.Code)
	}
	if res.Id != 6 {
		t.Errorf("response req_id = %d, want 6 (must echo the request's id)", res.Id)
	}
}

func TestHandleAbortDriverSetupEvent_UpdatesStateAndCallsCallback(t *testing.T) {
	i := newTestIntegration(t)

	called := false
	i.SetHandleAbortSetupFunction(func() { called = true })

	raw, err := json.Marshal(AbortDriverSetupEvent{
		CommonEvent: CommonEvent{Kind: "event", Msg: "abort_driver_setup", Cat: "DEVICE"},
		MsgData:     AbortDriverSetupData{Error: TimeoutError},
	})
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}

	i.handleEvent(&RequestMessage{CommonReq: CommonReq{Kind: "event", Msg: "abort_driver_setup"}}, raw)

	if !called {
		t.Error("the registered abort-setup callback was not invoked")
	}
	if i.SetupState != ErrorState {
		t.Errorf("SetupState = %q, want %q", i.SetupState, ErrorState)
	}
}

func TestSetDriverSetupState_UpdatesSetupState(t *testing.T) {
	i := newTestIntegration(t)

	i.SetDriverSetupState(SetupEvent, OkState, NoneError, nil)

	if i.SetupState != OkState {
		t.Errorf("SetupState = %q, want %q", i.SetupState, OkState)
	}
}

// TestHandleEntityCommandRequest_NewEntityTypesDispatchCorrectly exists specifically to prove the
// Phase 3 interface refactor's whole point: adding select and ir_emitter (Phase 6) required zero
// changes anywhere in this package - both satisfy entities.Entity purely by embedding BaseEntity
// and providing their own HandleCommand/UpdateEntity, so the dispatch path below (unchanged since
// before either type existed) already handles them correctly.
func TestHandleEntityCommandRequest_NewEntityTypesDispatchCorrectly(t *testing.T) {
	i := newTestIntegration(t)

	sel := entities.NewSelectEntity("select-1", entities.LanguageText{En: "Input"}, "", []string{"HDMI1", "HDMI2"})
	selected := ""
	sel.MapCommandWithParams(entities.SelectOptionSelectEntityCommand, func(params map[string]interface{}) error {
		selected, _ = params["option"].(string)
		return nil
	})

	ir := entities.NewIrEmitterEntity("ir-1", entities.LanguageText{En: "Blaster"}, "")
	sentCode := ""
	ir.MapCommandWithParams(entities.SendIrEmitterEntityCommand, func(params map[string]interface{}) error {
		sentCode, _ = params["code"].(string)
		return nil
	})

	i.Entities = append(i.Entities, sel, ir)

	selRaw, err := json.Marshal(EntityCommandReq{
		CommonReq: CommonReq{Kind: "req", Id: 10, Msg: "entity_command"},
		MsgData:   EntityCommandData{EntityId: "select-1", CmdId: "select_option", Params: map[string]interface{}{"option": "HDMI2"}},
	})
	if err != nil {
		t.Fatalf("marshal select request: %v", err)
	}
	resp := awaitOneMessage(i)
	i.handleRequest(&RequestMessage{CommonReq: CommonReq{Kind: "req", Id: 10, Msg: "entity_command"}}, selRaw)
	if res := decodeResponse(t, <-resp); res.Code != 200 {
		t.Errorf("select_option response code = %d, want 200", res.Code)
	}
	if selected != "HDMI2" {
		t.Errorf("selected option = %q, want \"HDMI2\"", selected)
	}

	irRaw, err := json.Marshal(EntityCommandReq{
		CommonReq: CommonReq{Kind: "req", Id: 11, Msg: "entity_command"},
		MsgData:   EntityCommandData{EntityId: "ir-1", CmdId: "send_ir", Params: map[string]interface{}{"code": "0000 006D"}},
	})
	if err != nil {
		t.Fatalf("marshal send_ir request: %v", err)
	}
	resp = awaitOneMessage(i)
	i.handleRequest(&RequestMessage{CommonReq: CommonReq{Kind: "req", Id: 11, Msg: "entity_command"}}, irRaw)
	if res := decodeResponse(t, <-resp); res.Code != 200 {
		t.Errorf("send_ir response code = %d, want 200", res.Code)
	}
	if sentCode != "0000 006D" {
		t.Errorf("sent IR code = %q, want \"0000 006D\"", sentCode)
	}
}

