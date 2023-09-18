package integration

import (
	log "github.com/sirupsen/logrus"
)

type DState string

const (
	ConnectedDeviceState    DState = "CONNECTED"
	ConnectingDeviceState          = "CONNECTING"
	DisconnectedDeviceState        = "DISCONNECTED"
	ErrorDeviceState               = "ERROR"
)

func (i *Integration) SetDeviceState(state DState) {
	log.WithField("DeviceState", state).Info("Set Device State")
	i.deviceState = state

	// Notify remote about new state
	i.sendDeviceStateEvent()
}
