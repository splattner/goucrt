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
