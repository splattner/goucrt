package integration

type DState string

const (
	ConnectedDeviceState    DState = "CONNECTED"
	ConnectingDeviceState          = "CONNECTING"
	DisconnectedDeviceState        = "DISCONNECTED"
	ErrorDeviceState               = "ERROR"
)
