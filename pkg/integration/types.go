package integration

import (
	"github.com/splattner/goucrt/pkg/entities"
)

// Common
type CommonReq struct {
	Kind string `json:"kind"`
	Id   int    `json:"id"`
	Msg  string `json:"msg"`
}

type CommonResp struct {
	Kind string `json:"kind"`
	Id   int    `json:"req_id"`
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}
type CommonEvent struct {
	Kind string `json:"kind"`
	Msg  string `json:"msg"`
	Cat  string `json:"cat"`
	Ts   string `json:"ts,omitempty"`
}

type LanguageText struct {
	En string `json:"en,omitempty"`
	De string `json:"de,omitempty"`
}

type DeviceId struct {
	DeviceId string `json:"device_id,omitempty"`
}

type DeviceState struct {
	DeviceId
	State string `json:"state"`
}

type DriverMetadata struct {
	DriverId        string          `json:"driver_id"`
	Name            LanguageText    `json:"name"`
	DriverUrl       string          `json:"driver_url,omitempty"`
	AuthMethod      string          `json:"auth_method,omitempty"`
	Version         string          `json:"version"`
	MinCoreAPI      string          `json:"min_core_api,omitempty"`
	Icon            string          `json:"icon,omitempty"`
	Description     LanguageText    `json:"description"`
	Developer       Developer       `json:"developer,omitempty"`
	HomePage        string          `json:"home_page,omitempty"`
	DeviceDiscovery bool            `json:"device_discovery,omitempty"`
	SetupDataSchema SetupDataSchema `json:"setup_data_schema,omitempty"`
	ReleaseData     string          `json:"release_date,omitempty"`
}

// Other ?

type DriverSetupState string
type DriverSetupEventType string
type DriverSetupError string

const (
	SetupState          DriverSetupState = "SETUP"
	WaitUserActionState DriverSetupState = "WAIT_USER_ACTION"
	OkState             DriverSetupState = "OK"
	ErrorState          DriverSetupState = "ERROR"
)

const (
	StartEvent DriverSetupEventType = "START"
	SetupEvent DriverSetupEventType = "SETUP"
	StopEvent  DriverSetupEventType = "STOP"
)

const (
	NoneError              DriverSetupError = "NONE"
	NotFoundError          DriverSetupError = "NOT_FOUND"
	ConnectionRefusedError DriverSetupError = "CONNECTION_REFUSED"
	AuthErrorError         DriverSetupError = "AUTHORIZATION_ERROR"
	TimeoutError           DriverSetupError = "TIMEOUT"
	OtherError             DriverSetupError = "OTHER"
)

type AvailableEntityFilter struct {
	DeviceId
	entities.EntityType
}

type Version struct {
	Api    string `json:"api"`
	Driver string `json:"driver"`
}

type Token struct {
	Token string `json:"token"`
}

type ConfirmationPage struct {
	Title    LanguageText `json:"title"`
	Message1 interface{}  `json:"message1,omitempty"`
	Image    string       `json:"image,omitempty"`
	Message2 interface{}  `json:"message2,omitempty"`
}

type SettigsPage struct {
	Title    LanguageText `json:"title"`
	Settings []Setting    `json:"settings"`
}

type Setting struct {
	Id    string       `json:"id"`
	Label LanguageText `json:"label"`
	Field interface{}  `json:"field"`
}

// Requests
type RequestMessage struct {
	CommonReq
	MsgData interface{} `json:"msg_data,omitempty"`
}

type AuthRequestMessage struct {
	CommonReq
	MsgData AuthRequestData `json:"msg_data"`
}

type AuthRequestData struct {
	Token string `json:"token"`
}

type DriverVersionReq struct {
	CommonReq
}

type DriverVersionData struct {
	Name    string  `json:"name"`
	Version Version `json:"version"`
}

type DriverMetadataReq struct {
	CommonReq
}

type AvailableEntityMessageReq struct {
	CommonReq
	MsgData AvailableEntityMessageData `json:"msg_data,omitempty"`
}

type AvailableEntityMessageData struct {
	Filter AvailableEntityFilter `json:"filter,omitempty"`
}

type DeviceStateMessageReq struct {
	CommonReq
	MsgData DeviceId
}

type SubscribeEventMessageReq struct {
	CommonReq
	MsgData SubscribeEventMessageData `json:"msg_data,omitempty"`
}

type SubscribeEventMessageData struct {
	DeviceId  string   `json:"device_id"`
	EntityIds []string `json:"entity_ids"`
}

type UnubscribeEventMessageReq struct {
	CommonReq
	MsgData SubscribeEventMessageData `json:"msg_data,omitempty"`
}

type UnubscribeEventMessageData struct {
	DeviceId  string   `json:"device_id"`
	EntityIds []string `json:"entity_ids"`
}

type GetEntityStatesMessageReq struct {
	CommonReq
	MsgData GetEntityStatesMessageData `json:"msg_data,omitempty"`
}

type GetEntityStatesMessageData struct {
	DeviceId string `json:"device_id"`
}

type EntityCommandReq struct {
	CommonReq
	MsgData EntityCommandData `json:"msg_data,omitempty"`
}

type EntityCommandData struct {
	DeviceId string                 `json:"device_id"`
	EntityId string                 `json:"entity_id"`
	CmdId    string                 `json:"cmd_id"`
	Params   map[string]interface{} `json:"params"`
}

type BrowseMediaMessageReq struct {
	CommonReq
	MsgData BrowseMediaData `json:"msg_data"`
}

type BrowseMediaData struct {
	EntityId  string                    `json:"entity_id"`
	MediaId   string                    `json:"media_id,omitempty"`
	MediaType entities.MediaContentType `json:"media_type,omitempty"`
	StableIds bool                      `json:"stable_ids,omitempty"`
	Paging    *entities.MediaPaging     `json:"paging,omitempty"`
}

type SearchMediaMessageReq struct {
	CommonReq
	MsgData SearchMediaData `json:"msg_data"`
}

type SearchMediaData struct {
	EntityId  string                      `json:"entity_id"`
	Query     string                      `json:"query"`
	MediaId   string                      `json:"media_id,omitempty"`
	MediaType entities.MediaContentType   `json:"media_type,omitempty"`
	StableIds bool                        `json:"stable_ids,omitempty"`
	Filter    *entities.MediaSearchFilter `json:"filter,omitempty"`
	Paging    *entities.MediaPaging       `json:"paging,omitempty"`
}

type MediaBrowseMessage struct {
	CommonResp
	MsgData MediaBrowseResponseData `json:"msg_data"`
}

type MediaBrowseResponseData struct {
	Media      *entities.BrowseMediaItem `json:"media,omitempty"`
	Pagination entities.MediaPagination  `json:"pagination"`
}

type MediaSearchMessage struct {
	CommonResp
	MsgData MediaSearchResponseData `json:"msg_data"`
}

type MediaSearchResponseData struct {
	Media      []entities.BrowseMediaItem `json:"media"`
	Pagination entities.MediaPagination   `json:"pagination"`
}

type SetupDriverMessageReq struct {
	CommonReq
	MsgData SetupDataValue `json:"msg_data"`
}

type SetupDataValue struct {
	Reconfigure bool      `json:"reconfigure,omitempty"`
	Value       SetupData `json:"setup_data"`
}

type SetupData map[string]string

// Set required data to configure the integration driver or continue the setup process.
type SetDriverUserDataRequest struct {
	CommonReq
	MsgData SetDriverUserData `json:"msg_data"`
}

type SetDriverUserData struct {
	InputValues map[string]string `json:"input_values,omitempty"`
	Confirm     bool              `json:"confirm,omitempty"`
}

// Response

type AuthenticationResponse struct {
	CommonResp
	MsgData DriverVersionData `json:"msg_data"`
}

type ResponseMessage struct {
	CommonResp
	MsgData interface{} `json:"msg_data,omitempty"`
}

// ErrorData is a result message's optional error payload: a short machine-readable code plus a
// human-readable message, e.g. {"code": "NOT_FOUND", "message": "message not known: foo"}.
type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type AvailableEntityData struct {
	Filter            AvailableEntityFilter `json:"filter,omitempty"`
	AvailableEntities []interface{}         `json:"available_entities"`
}

type AvailableEntityNoFilterData struct {
	AvailableEntities []interface{} `json:"available_entities"`
}

type AvailableEntityMessage struct {
	CommonResp
	MsgData AvailableEntityData `json:"msg_data"`
}

type AvailableEntityNoFilterMessage struct {
	CommonResp
	MsgData AvailableEntityNoFilterData `json:"msg_data"`
}

type DeviceStateEventMessage struct {
	CommonEvent
	MsgData DeviceState `json:"msg_data"`
}

type DriverMetadataReponse struct {
	CommonResp
	MsgData DriverMetadata `json:"msg_data"`
}

type Developer struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	EMail string `json:"email,omitempty"`
}

type SetupDataSchema struct {
	Title    LanguageText              `json:"title"`
	Settings []SetupDataSchemaSettings `json:"settings"`
}

type SetupDataSchemaSettings struct {
	Id    string       `json:"id"`
	Label LanguageText `json:"label"`
	Field interface{}  `json:"field"`
}

type SubscribeEventMessage struct {
	CommonResp
}

type UnubscribeEventMessage struct {
	CommonResp
}

type GetEntityStatesMessage struct {
	CommonResp
	MsgData []entities.EntityStateData `json:"msg_data,omitempty"`
}

type EntityCommandResponse struct {
	CommonResp
}

// Driver-initiated metadata requests (see GetVersion, GetSupportedEntityTypes,
// GetConfiguredEntities, GetLocalizationCfg, GetRuntimeInfo in metadata_requests.go). Unlike every
// other request in this file, these are sent BY the driver TO the Remote - the driver asking the
// Remote about itself, rather than the Remote asking the driver about its entities. The requests
// carry no msg_data (see the AsyncAPI spec's getVersionMsg/getSupportedEntityTypesMsg/... schemas),
// so a plain RequestMessage with a nil MsgData is enough to send one; only the responses need their
// own types.

type VersionInfo struct {
	// Model: short model identifier of the remote (UCR2 for Remote Two, UCR3 for Remote 3).
	Model      string `json:"model,omitempty"`
	DeviceName string `json:"device_name,omitempty"`
	Hostname   string `json:"hostname,omitempty"`
	// Address: MAC address of the remote.
	Address string `json:"address,omitempty"`
	Api     string `json:"api,omitempty"`
	Core    string `json:"core,omitempty"`
	Ui      string `json:"ui,omitempty"`
	Os      string `json:"os,omitempty"`
}

type VersionMessage struct {
	CommonResp
	MsgData VersionInfo `json:"msg_data"`
}

type SupportedEntityTypesMessage struct {
	CommonResp
	MsgData []string `json:"msg_data"`
}

type ConfiguredEntitiesMessage struct {
	CommonResp
	MsgData []string `json:"msg_data"`
}

type MeasurementUnit string

const (
	MetricMeasurementUnit MeasurementUnit = "METRIC"
	USMeasurementUnit     MeasurementUnit = "US"
	UKMeasurementUnit     MeasurementUnit = "UK"
)

type LocalizationSettings struct {
	// LanguageCode: language culture code, e.g. "en", "en_UK", "de_CH".
	LanguageCode string `json:"language_code,omitempty"`
	// CountryCode: two-letter ISO-3166-1-alpha-2 country code.
	CountryCode string `json:"country_code,omitempty"`
	// TimeZone: IANA time zone name, e.g. "Europe/Copenhagen".
	TimeZone        string          `json:"time_zone,omitempty"`
	TimeFormat24h   bool            `json:"time_format_24h"`
	MeasurementUnit MeasurementUnit `json:"measurement_unit,omitempty"`
}

type LocalizationCfgMessage struct {
	CommonResp
	MsgData LocalizationSettings `json:"msg_data"`
}

type RuntimeInfo struct {
	// DriverId: identifier under which this integration driver is accessible in the Remote.
	DriverId string `json:"driver_id"`
	// IntgIds: this driver's configured integration instance identifiers in the Remote. Only
	// single-device integration drivers are supported at the moment, so there's at most one.
	IntgIds []string `json:"intg_ids,omitempty"`
	// LogId: log service identifier for a custom integration driver running on the Remote.
	LogId string `json:"log_id,omitempty"`
}

type RuntimeInfoMessage struct {
	CommonResp
	MsgData RuntimeInfo `json:"msg_data"`
}

// Events

type AbortDriverSetupEvent struct {
	CommonEvent
	MsgData AbortDriverSetupData `json:"msg_data"`
}

type AbortDriverSetupData struct {
	Error DriverSetupError `json:"error"`
}

type EventMessage struct {
	CommonEvent
	MsgData interface{} `json:"msg_data"`
}

type EntityChangeEvent struct {
	CommonEvent
	MsgData EntityChangeData `json:"msg_data"`
}

type EntityChangeData struct {
	DeviceId   string                 `json:"device_id,omitempty"`
	EntityType string                 `json:"entity_type"`
	EntityId   string                 `json:"entity_id"`
	Attributes map[string]interface{} `json:"attributes"`
}
type EntityRemovedEvent struct {
	CommonEvent
	MsgData EntityRemovedEventData `json:"msg_data"`
}

type EntityRemovedEventData struct {
	DeviceId   string `json:"device_id,omitempty"`
	EntityType string `json:"entity_type"`
	EntityId   string `json:"entity_id"`
}

type EntityAvailableEvent struct {
	CommonEvent
	MsgData interface{} `json:"msg_data"`
}

type ConnectEvent struct {
	CommonEvent
	MsgData ConnectEventData `json:"msg_data,omitempty"`
}

type ConnectEventData struct {
	DeviceId string `json:"device_id,omitempty"`
}

type DriverSetupChangeEvent struct {
	CommonEvent
	MsgData DriverSetupChangeData `json:"msg_data,omitempty"`
}

type DriverSetupChangeData struct {
	EventType         DriverSetupEventType `json:"event_type"`
	State             DriverSetupState     `json:"state"`
	Error             DriverSetupError     `json:"error,omitempty"`
	RequireUserAction interface{}          `json:"require_user_action,omitempty"`
}

type RequireUserAction struct {
	Input        interface{} `json:"input,omitempty"`
	Confirmation interface{} `json:"confirmation,omitempty"`
}

// AssistantEventType is the "type" discriminator of an assistant_event message's msg_data - see
// Integration.sendAssistantEvent and its SendVoiceAssistant* wrappers in voice_assistant_events.go.
type AssistantEventType string

const (
	ReadyAssistantEventType          AssistantEventType = "ready"
	SttResponseAssistantEventType    AssistantEventType = "stt_response"
	TextResponseAssistantEventType   AssistantEventType = "text_response"
	SpeechResponseAssistantEventType AssistantEventType = "speech_response"
	FinishedAssistantEventType       AssistantEventType = "finished"
	ErrorAssistantEventType          AssistantEventType = "error"
)

// AssistantErrorCode is the assistant_event error event's "code" field.
type AssistantErrorCode string

const (
	ServiceUnavailableAssistantErrorCode AssistantErrorCode = "SERVICE_UNAVAILABLE"
	InvalidAudioAssistantErrorCode       AssistantErrorCode = "INVALID_AUDIO"
	NoTextRecognizedAssistantErrorCode   AssistantErrorCode = "NO_TEXT_RECOGNIZED"
	IntentFailedAssistantErrorCode       AssistantErrorCode = "INTENT_FAILED"
	TtsFailedAssistantErrorCode          AssistantErrorCode = "TTS_FAILED"
	TimeoutAssistantErrorCode            AssistantErrorCode = "TIMEOUT"
	UnexpectedErrorAssistantErrorCode    AssistantErrorCode = "UNEXPECTED_ERROR"
)

type AssistantEventMessage struct {
	CommonEvent
	MsgData AssistantEventData `json:"msg_data"`
}

type AssistantEventData struct {
	Type      AssistantEventType `json:"type"`
	EntityId  string             `json:"entity_id"`
	SessionId int                `json:"session_id"`
	// Data carries the event-specific payload: AssistantSttResponseData, AssistantTextResponseData,
	// AssistantSpeechResponseData or AssistantErrorData, depending on Type. Omitted for ready/finished.
	Data interface{} `json:"data,omitempty"`
}

type AssistantSttResponseData struct {
	Text string `json:"text"`
}

type AssistantTextResponseData struct {
	Success bool   `json:"success"`
	Text    string `json:"text"`
}

type AssistantSpeechResponseData struct {
	Url string `json:"url"`
	// MimeType: one of audio/mpeg, audio/mp3, audio/wav, audio/x-wav, audio/ogg, audio/opus,
	// audio/webm, audio/flac, audio/aac. Other types are ignored by the Remote's UI.
	MimeType string `json:"mime_type,omitempty"`
}

type AssistantErrorData struct {
	Code    AssistantErrorCode `json:"code"`
	Message string             `json:"message"`
}

type ButtonMapping struct {
	Button     string  `json:"string"`
	ShortPress Command `json:"short_press,omitempty"`
	LongPress  Command `json:"long_press,omitempty"`
}

type Command struct {
	CmdId  string                 `json:"cmd_id"`
	Params map[string]interface{} `json:"params,omitempty"`
}

type UserInterface struct {
	Pages []UserInterfacePage `json:"pages"`
}

type UserInterfacePage struct {
	PageID string `json:"page_id"`
	Name   string `json:"name,omitempty"`
	Grid   struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"grid"`
	Items []UserInterfaceItem `json:"items"`
}

type UserInterfaceItem struct {
	Type     string  `json:"type"`
	Icon     string  `json:"icon"`
	Text     string  `json:"text"`
	Command  Command `json:"command"`
	Location struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"location"`
	Size struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"size"`
}
