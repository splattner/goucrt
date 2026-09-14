package entities

// The voice_assistant entity type: receives a microphone input stream to forward voice commands to
// a device or cloud voice assistant (e.g. Amazon Alexa) or a two-way system like Home Assistant's
// assist pipeline. See doc/entities/entity_voice_assistant.md.
//
// Deliberately NOT implemented here: the actual audio stream. Per the spec, voice_start only starts
// it - the audio itself (RemoteVoiceBegin/RemoteVoiceData/RemoteVoiceEnd) travels as binary
// Protobuf WebSocket frames defined in a separate .proto file the Core-API repo doesn't vendor into
// doc/entities' worked examples, over a message framing pkg/integration's WebSocket loop doesn't
// handle at all today (wsReader only ever reads text/JSON frames). Wiring that up is a separate,
// substantially larger piece of work (vendoring the proto, generating Go bindings, teaching wsReader
// binary frame dispatch) than the JSON-only entity_command/entity_change/assistant_event surface
// implemented here, which is otherwise everything a driver needs to declare a voice_assistant entity
// and report on a voice command's outcome.

import "fmt"

type VoiceAssistantEntityState EntityState
type VoiceAssistantEntityFeatures EntityFeature
type VoiceAssistantEntityAttribute EntityAttribute
type VoiceAssistantEntityCommand EntityCommand
type VoiceAssistantEntityOption EntityOption

const (
	OffVoiceAssistantEntityState VoiceAssistantEntityState = "OFF"
	OnVoiceAssistantEntityState  VoiceAssistantEntityState = "ON"
)

const (
	// TranscriptionVoiceAssistantEntityFeatures: the voice assistant can provide a speech-to-text
	// transcription of the voice command, via a stt_response assistant_event.
	TranscriptionVoiceAssistantEntityFeatures VoiceAssistantEntityFeatures = "transcription"
	// ResponseTextVoiceAssistantEntityFeatures: the voice assistant can provide a textual response
	// about the performed action, via a text_response assistant_event.
	ResponseTextVoiceAssistantEntityFeatures VoiceAssistantEntityFeatures = "response_text"
	// ResponseSpeechVoiceAssistantEntityFeatures: the voice assistant can provide a speech response
	// about the performed action, via a speech_response assistant_event.
	ResponseSpeechVoiceAssistantEntityFeatures VoiceAssistantEntityFeatures = "response_speech"
)

const (
	// StateVoiceAssistantEntityAttribute is the entity's only attribute, always present regardless
	// of which features are declared (the spec's Attributes table lists no enabling feature for
	// it, unlike every other entity type's state attribute).
	StateVoiceAssistantEntityAttribute VoiceAssistantEntityAttribute = "state"
)

const (
	// VoiceStartVoiceAssistantEntityCommand starts a voice command's audio stream. Params:
	// session_id (int, required) identifies the stream in the binary voice messages that follow;
	// audio_cfg (AudioConfiguration, required) is the format that stream will use; speech_response
	// (bool, optional) requests a speech response for this command; timeout (int, optional)
	// processing timeout in seconds; profile_id (string, optional) selects a VoiceAssistantProfile.
	//
	// Must be confirmed within 2 seconds. The integration must follow up with a "ready"
	// assistant_event (see Integration.SendVoiceAssistantReady) once ready to receive the audio
	// stream identified by session_id.
	VoiceStartVoiceAssistantEntityCommand VoiceAssistantEntityCommand = "voice_start"
)

const (
	AudioCfgVoiceAssistantEntityOption         VoiceAssistantEntityOption = "audio_cfg"
	ProfilesVoiceAssistantEntityOption         VoiceAssistantEntityOption = "profiles"
	PreferredProfileVoiceAssistantEntityOption VoiceAssistantEntityOption = "preferred_profile"
)

// AudioSampleFormat is an audio stream's sample format, for AudioConfiguration.SampleFormat.
type AudioSampleFormat string

const (
	I16AudioSampleFormat AudioSampleFormat = "I16"
	I32AudioSampleFormat AudioSampleFormat = "I32"
	U16AudioSampleFormat AudioSampleFormat = "U16"
	U32AudioSampleFormat AudioSampleFormat = "U32"
	F32AudioSampleFormat AudioSampleFormat = "F32"
)

// AudioConfiguration is the audio_cfg option/voice_start param: the format of the audio stream sent
// to the integration. The Remote's default (used for any zero field here, or if the requested
// format isn't supported) is PCM 16 kHz mono 16-bit signed, sent in 100-200ms chunks.
type AudioConfiguration struct {
	// Channels: 1 or 2. Default 1.
	Channels int `json:"channels,omitempty"`
	// SampleRate in Hz: one of 8000, 11025, 16000, 22050, 24000, 44100, 48000. Default 16000.
	SampleRate int `json:"sample_rate,omitempty"`
	// SampleFormat: default I16AudioSampleFormat.
	SampleFormat AudioSampleFormat `json:"sample_format,omitempty"`
}

// VoiceAssistantProfile is one entry in the profiles option: a named, selectable voice input
// configuration (e.g. a language or a local-vs-cloud processing choice), selectable via
// voice_start's profile_id param.
type VoiceAssistantProfile struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	// Language: optional language code if this profile is for a specific speech-recognition
	// language.
	Language string `json:"language,omitempty"`
	// Features overrides the entity's own Features for this profile - e.g. a profile with fewer or
	// more capabilities than the entity's default. An empty, non-nil slice means "no features" for
	// this profile; nil means "use the entity's own features".
	Features []VoiceAssistantEntityFeatures `json:"features,omitempty"`
}

// VoiceAssistantEntity interacts with a voice assistant or voice-capable device: it can request a
// microphone audio stream (see voice_start/SendVoiceAssistantReady) and report the outcome via
// assistant_event (see Integration.SendVoiceAssistant{SttResponse,TextResponse,SpeechResponse,
// Finished,Error}).
type VoiceAssistantEntity struct {
	BaseEntity
	Commands map[VoiceAssistantEntityCommand]func(VoiceAssistantEntity, map[string]interface{}) int `json:"-"`
	Options  map[VoiceAssistantEntityOption]interface{}                                             `json:"options,omitempty"`
}

func NewVoiceAssistantEntity(id string, name LanguageText, area string) *VoiceAssistantEntity {

	voiceAssistantEntity := VoiceAssistantEntity{}
	voiceAssistantEntity.Id = id
	voiceAssistantEntity.Name = name
	voiceAssistantEntity.Area = area

	voiceAssistantEntity.Type = "voice_assistant"

	voiceAssistantEntity.Commands = make(map[VoiceAssistantEntityCommand]func(VoiceAssistantEntity, map[string]interface{}) int)
	voiceAssistantEntity.Attributes = make(map[string]interface{})
	voiceAssistantEntity.Options = make(map[VoiceAssistantEntityOption]interface{})

	voiceAssistantEntity.AddAttribute(string(StateVoiceAssistantEntityAttribute), OffVoiceAssistantEntityState)

	return &voiceAssistantEntity
}

func (e *VoiceAssistantEntity) UpdateEntity(newEntity interface{}) error {
	updated, ok := newEntity.(VoiceAssistantEntity)
	if !ok {
		return fmt.Errorf("cannot update VoiceAssistantEntity from %T", newEntity)
	}

	e.Name = updated.Name
	e.Area = updated.Area
	e.Commands = updated.Commands
	e.Features = updated.Features
	e.Attributes = updated.Attributes
	e.Options = updated.Options

	return nil
}

// AddFeature declares a capability of this voice assistant. Unlike most other entity types, no
// feature here gates an entity attribute: transcription/response_text/response_speech are all
// reported live through assistant_event messages (stt_response/text_response/speech_response), not
// through a steady-state attribute - the spec's Attributes table lists only "state", enabled by
// none of these features.
func (e *VoiceAssistantEntity) AddFeature(feature VoiceAssistantEntityFeatures) {
	e.Features = append(e.Features, feature)
}

// Add an option to the VoiceAssistant Entity
func (e *VoiceAssistantEntity) AddOption(option VoiceAssistantEntityOption, value interface{}) {
	e.Options[option] = value
}

// Register a function for the Entity command
func (e *VoiceAssistantEntity) AddCommand(command VoiceAssistantEntityCommand, function func(VoiceAssistantEntity, map[string]interface{}) int) {
	e.Commands[command] = function
}

// Map a VoiceAssistantEntityCommand to a function call with params (voice_start)
func (e *VoiceAssistantEntity) MapCommandWithParams(command VoiceAssistantEntityCommand, f func(map[string]interface{}) error) {

	e.AddCommand(command, func(entity VoiceAssistantEntity, params map[string]interface{}) int {

		if err := f(params); err != nil {
			return 404
		}
		return 200
	})
}

// Call the registred function for this entity_command
func (e *VoiceAssistantEntity) HandleCommand(cmd_id string, params map[string]interface{}) int {
	if e.Commands[VoiceAssistantEntityCommand(cmd_id)] != nil {
		return e.Commands[VoiceAssistantEntityCommand(cmd_id)](*e, params)
	}

	return 404
}
