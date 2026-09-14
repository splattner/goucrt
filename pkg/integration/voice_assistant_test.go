package integration

// Tests for the voice_assistant entity's protocol surface: the voice_start entity_command (which,
// being a plain JSON command with params, needs no special-case handling in requests.go - it's
// covered by the same handleEntityCommandRequest path as every other entity) and the
// assistant_event senders in voice_assistant_events.go.

import (
	"encoding/json"
	"testing"

	"github.com/splattner/goucrt/pkg/entities"
)

func TestHandleEntityCommandRequest_VoiceAssistantVoiceStart(t *testing.T) {
	i := newTestIntegration(t)

	va := entities.NewVoiceAssistantEntity("va-1", entities.LanguageText{En: "Assistant"}, "")
	var gotSessionId float64
	var gotProfileId string
	va.MapCommandWithParams(entities.VoiceStartVoiceAssistantEntityCommand, func(params map[string]interface{}) error {
		gotSessionId, _ = params["session_id"].(float64)
		gotProfileId, _ = params["profile_id"].(string)
		return nil
	})
	i.Entities = append(i.Entities, va)

	raw, err := json.Marshal(EntityCommandReq{
		CommonReq: CommonReq{Kind: "req", Id: 30, Msg: "entity_command"},
		MsgData: EntityCommandData{
			EntityId: "va-1",
			CmdId:    "voice_start",
			Params: map[string]interface{}{
				"session_id": float64(8),
				"audio_cfg": map[string]interface{}{
					"channels":      float64(1),
					"sample_rate":   float64(8000),
					"sample_format": "I16",
				},
				"profile_id": "en-local",
			},
		},
	})
	if err != nil {
		t.Fatalf("marshal voice_start request: %v", err)
	}

	resp := awaitOneMessage(i)
	i.handleRequest(&RequestMessage{CommonReq: CommonReq{Kind: "req", Id: 30, Msg: "entity_command"}}, raw)
	if res := decodeResponse(t, <-resp); res.Code != 200 {
		t.Errorf("voice_start response code = %d, want 200", res.Code)
	}
	if gotSessionId != 8 {
		t.Errorf("session_id = %v, want 8", gotSessionId)
	}
	if gotProfileId != "en-local" {
		t.Errorf("profile_id = %q, want \"en-local\"", gotProfileId)
	}
}

func TestSendVoiceAssistantEvents(t *testing.T) {
	i := newTestIntegration(t)

	cases := []struct {
		name     string
		send     func()
		wantType AssistantEventType
	}{
		{"Ready", func() { i.SendVoiceAssistantReady("va-1", 8) }, ReadyAssistantEventType},
		{"SttResponse", func() { i.SendVoiceAssistantSttResponse("va-1", 8, "turn off the lights") }, SttResponseAssistantEventType},
		{"TextResponse", func() { i.SendVoiceAssistantTextResponse("va-1", 8, true, "lights off") }, TextResponseAssistantEventType},
		{"SpeechResponse", func() { i.SendVoiceAssistantSpeechResponse("va-1", 8, "https://example.com/a.mp3", "audio/mpeg") }, SpeechResponseAssistantEventType},
		{"Finished", func() { i.SendVoiceAssistantFinished("va-1", 8) }, FinishedAssistantEventType},
		{"Error", func() { i.SendVoiceAssistantError("va-1", 8, TimeoutAssistantErrorCode, "took too long") }, ErrorAssistantEventType},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			done := make(chan []byte, 1)
			go func() { done <- <-i.Remote.messageChannel }()

			c.send()

			var msg AssistantEventMessage
			if err := json.Unmarshal(<-done, &msg); err != nil {
				t.Fatalf("decode assistant_event: %v", err)
			}

			if msg.Kind != "event" || msg.Msg != "assistant_event" || msg.Cat != "ENTITY" {
				t.Errorf("envelope = {kind: %q, msg: %q, cat: %q}, want {event, assistant_event, ENTITY}", msg.Kind, msg.Msg, msg.Cat)
			}
			if msg.MsgData.Type != c.wantType {
				t.Errorf("msg_data.type = %q, want %q", msg.MsgData.Type, c.wantType)
			}
			if msg.MsgData.EntityId != "va-1" || msg.MsgData.SessionId != 8 {
				t.Errorf("msg_data = {entity_id: %q, session_id: %d}, want {va-1, 8}", msg.MsgData.EntityId, msg.MsgData.SessionId)
			}
		})
	}
}

// TestSendVoiceAssistantSttResponse_DataMatchesSpecShape checks the actual data payload, not just
// the envelope - specifically that it round-trips as the spec's {"text": "..."} shape rather than,
// say, a bare string (an easy mistake since AssistantEventData.Data is an interface{}).
func TestSendVoiceAssistantSttResponse_DataMatchesSpecShape(t *testing.T) {
	i := newTestIntegration(t)

	done := make(chan []byte, 1)
	go func() { done <- <-i.Remote.messageChannel }()

	i.SendVoiceAssistantSttResponse("va-1", 8, "turn off the lights")

	var decoded struct {
		MsgData struct {
			Data struct {
				Text string `json:"text"`
			} `json:"data"`
		} `json:"msg_data"`
	}
	if err := json.Unmarshal(<-done, &decoded); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if decoded.MsgData.Data.Text != "turn off the lights" {
		t.Errorf("msg_data.data.text = %q, want \"turn off the lights\"", decoded.MsgData.Data.Text)
	}
}
