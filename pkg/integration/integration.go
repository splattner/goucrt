package integration

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/grandcat/zeroconf"
)

const API_VERSION = "0.8.1-alpha"

type Integration struct {
	DeviceId string

	Metadata *DriverMetadata

	authToken string

	deviceState DState

	config Config

	Remote remote

	Entities []interface{}

	SubscribedEntities []string

	handleSetupFunction             func(SetupData)
	handleConnectionFunction        func(*ConnectEvent)
	handleSetDriverUserDataFunction func(map[string]string, bool)

	SetupState DriverSetupState

	mdns *zeroconf.Server
}

func NewIntegration(config Config) (*Integration, error) {

	i := Integration{
		config:      config,
		deviceState: DisconnectedDeviceState,
		DeviceId:    "", // I think device_id is not yet implemented in Remote TV, used for multi-device integrati

	}

	i.Remote.messageChannel = make(chan []byte)

	return &i, nil

}

func (i *Integration) SetMetadata(metadata *DriverMetadata) {
	log.WithField("Metadata", metadata).Debug("Set Metadata")
	i.Metadata = metadata
}

func (i *Integration) Run() error {

	if i.Metadata == nil {
		log.Error("Metadata not set, cannot run")
		return fmt.Errorf("Metadata not set")
	}

	http.HandleFunc("/ws", i.wsEndpoint)

	listenAddress := fmt.Sprintf(":%d", i.config["listenport"].(int))

	// MDNS
	i.startAdvertising()

	log.Fatal(http.ListenAndServe(listenAddress, nil))

	return nil

}

// Set the function which is called when the setup_driver request was sent by the remote
func (i *Integration) SetHandleSetupFunction(f func(SetupData)) {
	i.handleSetupFunction = f
}

// Set the function which is called when the connect/disconnect request was sent by the remote
func (i *Integration) SetHandleConnectionFunction(f func(*ConnectEvent)) {
	i.handleConnectionFunction = f
}

// Set the function which is called when the connect/disconnect request was sent by the remote
func (i *Integration) SetHandleSetDriverUserDataFunction(f func(map[string]string, bool)) {
	i.handleSetDriverUserDataFunction = f
}

func (i *Integration) SetDriverSetupState(event_Type DriverSetupEventType, state DriverSetupState, err DriverSetupError, requireUserAction *RequireUserAction) {

	log.WithFields(log.Fields{
		"EventType": event_Type,
		"State":     state,
		"Error":     err,
	}).Info("Set DriverSetup State from Client")

	// Overwrite state if requireUserAction is set
	if requireUserAction != nil {
		state = WaitUserActionState
	}

	i.sendDriverSetupChangeEvent(event_Type, state, err, requireUserAction)

}
