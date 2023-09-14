package integration

import "github.com/splattner/goucrt/pkg/integration/entities"

// Common
type CommonReq struct {
	Kind string `json:"kind"`
	Id   int    `json:"id"`
	Msg  string `json:"msg"`
}

type CommonResp struct {
	Kind string `json:"kind"`
	Id   int    `json:"id"`
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
	En string `json:"en"`
	De string `json:"de"`
}

type DeviceId struct {
	DeviceId string `json:"device_id"`
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
	Developer       Developer       `json:"description,omitempty"`
	HomePage        string          `json:"home_page,omitempty`
	DeviceDiscovery bool            `json:"device_discovery,omitempty"`
	SetupDataSchema SetupDataSchema `json:"setup_data_schema,omitempty"`
	ReleaseData     string          `json:"release_date,omitempty`
}

// Other ?

type DriverSetupState string
type DriverSetupEventType string
type DriverSetupError string

const (
	SetupState              DriverSetupState = "SETUP"
	WaitUserActionState                      = "WAIT_USER_ACTION"
	RequiredUserActionState                  = "WAIT_USER_ACTION"
	OkState                                  = "OK"
	ErrorState                               = "ERROR"
)

const (
	StartEvent DriverSetupEventType = "START"
	SetupEvent                      = "SETUP"
	StopEvent                       = "STOP"
)

const (
	NoneError              DriverSetupError = "NONE"
	NotFoundError                           = "NOT_FOUND"
	ConnectionRefusedError                  = "CONNECTION_REFUSED"
	AuthErrorError                          = "AUTHORIZATION_ERROR"
	TimeoutError                            = "TIMEOUT"
	OtherError                              = "OTHER"
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
	Message1 LanguageText `json:"message1,omitempty"`
	Image    string       `json:"image,omitempty"`
	Message2 LanguageText `json:"message2,omitempty"`
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

type GetEntityStateMessageReq struct {
	CommonReq
	MsgData GetEntityStateMessageData `json:"msg_data,omitempty"`
}

type GetEntityStateMessageData struct {
	DeviceId string `json:"device_id"`
}

type EntityCommandReq struct {
	CommonReq
	MsgData EntityCommandData `json:"msg_data,omitempty"`
}

type EntityCommandData struct {
	DeviceId string      `json:"device_id"`
	EntityId string      `json:"entity_id"`
	CmdId    string      `json:"cmd_id"`
	Params   interface{} `json:"params"`
}

type SetupDriverMessageReq struct {
	CommonReq
	MsgData SettingsVaulues `json:"msg_data"`
}

type SettingsVaulues struct {
	Value map[string]string `json:"setup_data"`
}

// Set required data to configure the integration driver or continue the setup process.
type SetDriverUserDataReq struct {
	CommonReq
	MsgData interface{} `json:"msg_data`
}

// Response
type ResponseMessage struct {
	CommonResp
	MsgData interface{} `json:"msg_data,omitempty"`
}

type AvailableEntityData struct {
	Filter            AvailableEntityFilter `json:"filter"`
	AvailableEntities []entities.Entity     `json:"available_entities`
}

type AvailableEntityMessage struct {
	CommonResp
	MsgData AvailableEntityData `json:"msg_data"`
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
	Title map[string]string `json:"title"`
}

type SubscribeEventMessage struct {
	CommonResp
}

type UnubscribeEventMessage struct {
	CommonResp
}

type GetEntityStateMessage struct {
	CommonResp
	MsgData []entities.EntityStateData `json:"msg_data,omitempty"`
}

type EntityCommandResponse struct {
	CommonResp
}

// Events
type EventMessage struct {
	CommonEvent
	MsgData interface{} `json:"msg_data"`
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
	MsgData entities.Entity `json:"msg_data"`
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
	EventType          DriverSetupEventType `json:"event_type"`
	State              DriverSetupState     `json:"state"`
	Error              DriverSetupError     `json:"error"`
	RequiredUserAction interface{}          `json:"required_user_action"`
}
