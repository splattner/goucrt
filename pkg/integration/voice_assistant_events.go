package integration

// Sends the assistant_event message a voice_assistant entity uses to report on a voice command
// started by a voice_start entity_command: from the initial "ready" signal through to a
// transcription, a textual or spoken response, and finally "finished" (or "error" at any point).
// See doc/entities/entity_voice_assistant.md's "Assistant event" section and
// pkg/entities/voice_assistant.go's package doc for what's out of scope (the actual audio stream).

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

// sendAssistantEvent sends one assistant_event frame for entityId/sessionId. data is the
// event-specific payload (AssistantSttResponseData, ...), or nil for ready/finished.
func (i *Integration) sendAssistantEvent(entityId string, sessionId int, eventType AssistantEventType, data interface{}) {
	var res interface{}
	now := time.Now()

	res = AssistantEventMessage{
		CommonEvent{Kind: "event", Msg: "assistant_event", Cat: "ENTITY", Ts: now.Format(time.RFC3339)},
		AssistantEventData{Type: eventType, EntityId: entityId, SessionId: sessionId, Data: data},
	}

	if err := i.sendEventMessage(&res, websocket.TextMessage); err != nil {
		log.WithError(err).Error("Cannot send Assistant Event Message")
	}
}

// SendVoiceAssistantReady signals that the integration is ready to receive the voice audio stream
// for sessionId, after confirming its voice_start command. Must be sent for every voice_start the
// entity accepts - the Remote won't start the audio stream without it.
func (i *Integration) SendVoiceAssistantReady(entityId string, sessionId int) {
	i.sendAssistantEvent(entityId, sessionId, ReadyAssistantEventType, nil)
}

// SendVoiceAssistantSttResponse reports the transcribed text of the voice command. Requires the
// TranscriptionVoiceAssistantEntityFeatures feature.
func (i *Integration) SendVoiceAssistantSttResponse(entityId string, sessionId int, text string) {
	i.sendAssistantEvent(entityId, sessionId, SttResponseAssistantEventType, AssistantSttResponseData{Text: text})
}

// SendVoiceAssistantTextResponse reports a textual description of the action the voice command
// performed. Requires the ResponseTextVoiceAssistantEntityFeatures feature.
func (i *Integration) SendVoiceAssistantTextResponse(entityId string, sessionId int, success bool, text string) {
	i.sendAssistantEvent(entityId, sessionId, TextResponseAssistantEventType, AssistantTextResponseData{Success: success, Text: text})
}

// SendVoiceAssistantSpeechResponse reports a spoken audio response to the voice command, as a URL
// the Remote fetches. Requires the ResponseSpeechVoiceAssistantEntityFeatures feature.
func (i *Integration) SendVoiceAssistantSpeechResponse(entityId string, sessionId int, url string, mimeType string) {
	i.sendAssistantEvent(entityId, sessionId, SpeechResponseAssistantEventType, AssistantSpeechResponseData{Url: url, MimeType: mimeType})
}

// SendVoiceAssistantFinished signals that voice command processing for sessionId is complete. Must
// be sent to close out every session that got a SendVoiceAssistantReady, whether or not an error
// occurred along the way.
func (i *Integration) SendVoiceAssistantFinished(entityId string, sessionId int) {
	i.sendAssistantEvent(entityId, sessionId, FinishedAssistantEventType, nil)
}

// SendVoiceAssistantError reports a failure processing or sending the voice command's audio stream.
func (i *Integration) SendVoiceAssistantError(entityId string, sessionId int, code AssistantErrorCode, message string) {
	i.sendAssistantEvent(entityId, sessionId, ErrorAssistantEventType, AssistantErrorData{Code: code, Message: message})
}
