package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/grandcat/zeroconf"
)

const API_VERSION = "0.8.1-alpha"

type Integration struct {
	DeviceId string
	DriverId string

	Metadata *DriverMetadata

	authToken string

	deviceState DState

	Config        Config
	listenAddress string

	Remote remote

	Entities []interface{}

	SubscribedEntities []string

	handleSetupFunction             func(SetupData)
	handleConnectionFunction        func(*ConnectEvent)
	handleSetDriverUserDataFunction func(map[string]string, bool)

	SetupState DriverSetupState

	SetupData SetupData

	mdns *zeroconf.Server
}

func NewIntegration(config Config) (*Integration, error) {

	i := Integration{
		Config:        config,
		listenAddress: fmt.Sprintf(":%d", config.ListenPort),
		deviceState:   DisconnectedDeviceState,
		DeviceId:      "", // I think device_id is not yet implemented in Remote TV, used for multi-device integrati

	}

	i.Remote.messageChannel = make(chan []byte)
	i.Remote.controlChannel = make(chan string)

	return &i, nil

}

func (i *Integration) SetMetadata(metadata *DriverMetadata) {
	log.WithField("Metadata", metadata).Debug("Set Metadata")
	i.Metadata = metadata

	i.LoadSetupData()
}

func (i *Integration) Run() error {
	log.Info("Start Remote Two integration")

	defer func() {
		i.stopAdvertising()
	}()

	if i.Metadata == nil {
		log.Panic("Metadata not set, cannot start Remote Two integration")
		return fmt.Errorf("Metadata not set")
	}

	http.HandleFunc(i.Config.WebsocketPath, i.wsEndpoint)

	//MDNS
	if !i.Config.DisableMDNS {
		go i.startAdvertising()
	}

	// Register the integration
	if i.Config.EnableRegistration && i.Config.RegistrationPin != "" {
		go i.registerIntegration()
	}

	log.Debug("Listen for new Websocket connection")
	log.Fatal(http.ListenAndServe(i.listenAddress, nil))

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

// Set and then Send the Driver Setup State to Remote two
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

// Load persist setupData File
// TODO: handle location via ENV's
func (i *Integration) LoadSetupData() {

	file, err := os.ReadFile(i.Config.ConfigHome + i.Metadata.DriverId + ".json")
	if err != nil {
		log.WithError(err).Info("Cannot read setupDataFile")
		i.SetupData = make(SetupData)
	} else {
		json.Unmarshal(file, &i.SetupData)
		log.WithField("SetupData", i.SetupData).Info("Read persisted setup data")
	}
}

// Persist File
// TODO: handle location via ENV's
func (i *Integration) PersistSetupData() {

	log.WithField("SetupData", i.SetupData).Info("Persist setup data")
	file, _ := json.MarshalIndent(i.SetupData, "", " ")
	_ = os.WriteFile(i.Config.ConfigHome+i.Metadata.DriverId+".json", file, 0644)
}
