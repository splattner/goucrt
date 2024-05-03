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
