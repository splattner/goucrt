package integration

import "log"

type DState string

const (
	ConnectedDeviceState    DState = "CONNECTED"
	ConnectingDeviceState          = "CONNECTING"
	DisconnectedDeviceState        = "DISCONNECTED"
	ErrorDeviceState               = "ERROR"
)

func (i *Integration) SetDeviceState(state DState) {
	log.Println("Set Device state to:" + state)
	i.deviceState = state

	// Notify remote about new state
	i.sendDeviceStateEvent()
}
