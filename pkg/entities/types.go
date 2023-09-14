package entities

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

type EntityCommandResponse struct {
	CommonResp
}
