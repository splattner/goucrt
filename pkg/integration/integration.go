package integration

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/grandcat/zeroconf"
	"github.com/splattner/goucrt/pkg/entities"
)

// Environment variables the remote itself sets when running a custom-installed driver in its
// sandbox. See doc/integration-driver/driver-installation.md's "Runtime environment" section.
// Nothing else sets these, so UCIntegrationHTTPPortEnv's presence also doubles as the signal that
// this process is running in that sandbox at all.
const (
	UCIntegrationHTTPPortEnv  = "UC_INTEGRATION_HTTP_PORT"
	UCIntegrationInterfaceEnv = "UC_INTEGRATION_INTERFACE"
)

// isSandboxed reports whether this process is running as a custom-installed driver in the
// remote's sandbox, as opposed to standalone (container, bare process) with mDNS discovery and/or
// registration. The remote sets UCIntegrationHTTPPortEnv itself for this purpose.
func isSandboxed() bool {
	_, set := os.LookupEnv(UCIntegrationHTTPPortEnv)
	return set
}

const API_VERSION = "0.10.0"

type Integration struct {
	DeviceId string
	DriverId string

	Metadata *DriverMetadata

	//authToken string

	deviceState DState

	Config        Config
	listenAddress string

	Remote remote

	// entitiesMu guards Entities and SubscribedEntities: the WebSocket read loop (handleRequest/
	// handleEvent) and a driver's own goroutines (discovery callbacks, MQTT/WS handlers calling
	// AddEntity, RemoveEntity, SendEntityChangeEvent) both touch them concurrently.
	entitiesMu sync.RWMutex

	Entities []entities.Entity

	SubscribedEntities []string

	// requestsMu guards pendingRequests and nextReqID: sendMetadataRequest (called from a driver's
	// own goroutine, e.g. calling GetVersion) registers a channel before sending, and the WebSocket
	// read loop's handleResponse delivers a matching "resp" frame to it by req_id - both run
	// concurrently with each other and with further sendMetadataRequest calls.
	requestsMu      sync.Mutex
	pendingRequests map[int]chan []byte
	nextReqID       int

	handleSetupFunction             func(SetupData)
	handleConnectionFunction        func(*ConnectEvent)
	handleSetDriverUserDataFunction func(map[string]string, bool)
	handleAbortSetupFunction        func()

	SetupState DriverSetupState

	SetupData SetupData

	mdns *zeroconf.Server
}

func NewIntegration(config Config) (*Integration, error) {

	// UC_INTEGRATION_HTTP_PORT/UC_INTEGRATION_INTERFACE are set by the remote itself when running
	// as a custom-installed driver, and are authoritative when present - the remote is telling the
	// driver exactly what to bind to, overriding any configured port/interface. Handled directly
	// here (rather than only via CLI flag/env binding) so this is correct for any entry point that
	// constructs an Integration, not just ones that remember to wire the binding themselves.
	if portEnv, ok := os.LookupEnv(UCIntegrationHTTPPortEnv); ok {
		if port, err := strconv.Atoi(portEnv); err == nil {
			config.ListenPort = port
		} else {
			log.WithError(err).WithField("value", portEnv).Errorf("Cannot parse %s, ignoring it", UCIntegrationHTTPPortEnv)
		}
	}
	if iface, ok := os.LookupEnv(UCIntegrationInterfaceEnv); ok {
		config.BindInterface = iface
	}

	if config.DataHome == "" {
		config.DataHome = config.ConfigHome
	}

	if isSandboxed() {
		// mDNS discovery and self-registration are meaningless here: the remote already knows
		// about this driver from the installation archive and connects to it directly. There's
		// also no way for a driver.json-installed driver to pass CLI flags to disable these itself
		// (the remote runs the binary with no arguments), so this can only be handled here.
		log.Info("Running as a custom-installed driver (UC_INTEGRATION_HTTP_PORT is set); disabling mDNS advertisement and self-registration")
		config.DisableMDNS = true
		config.EnableRegistration = false
	}

	bindHost := config.BindInterface
	if bindHost == "" {
		// TODO: for the moment, only IPv4, as somehow the behaviour seems strange when both.. not investigated though
		bindHost = "0.0.0.0"
	}

	i := Integration{
		Config:        config,
		listenAddress: net.JoinHostPort(bindHost, strconv.Itoa(config.ListenPort)),
		deviceState:   DisconnectedDeviceState,
		DeviceId:      "", // I think device_id is not yet implemented in Remote TV, used for multi-device integrati

	}

	i.Remote.messageChannel = make(chan []byte)
	i.Remote.controlChannel = make(chan string)
	i.pendingRequests = make(map[int]chan []byte)

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
		return fmt.Errorf("metadata not set, cannot start Remote Two integration")
	}

	http.HandleFunc(i.Config.WebsocketPath, i.wsEndpoint)

	//MDNS
	if !i.Config.DisableMDNS {
		go i.startAdvertising()
	}

	// Register the integration
	if i.Config.EnableRegistration && i.Config.RegistrationPin != "" {
		go func() {
			if err := i.registerIntegration(); err != nil {
				log.WithError(err).Error("Cannot register integration with the Remote")
			}
		}()
	}

	log.Debug("Listen for new Websocket connection")

	return http.ListenAndServe(i.listenAddress, nil)

}

// Set the function which is called when the setup_driver request was sent by the remote
func (i *Integration) SetHandleSetupFunction(f func(SetupData)) {
	i.handleSetupFunction = f
}

// Set the function which is called when the connect/disconnect request was sent by the remote
func (i *Integration) SetHandleConnectionFunction(f func(*ConnectEvent)) {
	i.handleConnectionFunction = f
}

// Set the function which is called when the setDriverUsaerData request was sent by the remote
func (i *Integration) SetHandleSetDriverUserDataFunction(f func(map[string]string, bool)) {
	i.handleSetDriverUserDataFunction = f
}

// Set the function which is called when the remote sends abort_driver_setup, i.e. the user
// cancelled the setup flow. Use it to stop any in-progress setup work (discovery, open
// connections, etc.) - per the spec, further messages from the driver about this setup attempt
// are ignored by the remote after this point.
func (i *Integration) SetHandleAbortSetupFunction(f func()) {
	i.handleAbortSetupFunction = f
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

	i.SetupState = state

	i.sendDriverSetupChangeEvent(event_Type, state, err, requireUserAction)

}

// setupDataPath returns the path of the persisted setup data file: <ConfigHome>/<driver_id>.json.
func (i *Integration) setupDataPath() string {
	return filepath.Join(i.Config.ConfigHome, i.Metadata.DriverId+".json")
}

// Load persist setupData File
func (i *Integration) LoadSetupData() {

	file, err := os.ReadFile(i.setupDataPath())
	if err != nil {
		log.WithError(err).Info("Cannot read setupDataFile")
		i.SetupData = make(SetupData)
	} else {
		if err := json.Unmarshal(file, &i.SetupData); err != nil {
			log.WithError(err).Error("Cannot unmarshall setSetupDataupdata")
		}
		log.WithField("SetupData", i.SetupData).Info("Read persisted setup data")
	}
}

// Persist File
func (i *Integration) PersistSetupData() error {

	log.WithField("SetupData", i.SetupData).Info("Persist setup data")

	if err := os.MkdirAll(i.Config.ConfigHome, 0o755); err != nil {
		return fmt.Errorf("cannot create config directory %q: %w", i.Config.ConfigHome, err)
	}

	file, err := json.MarshalIndent(i.SetupData, "", " ")
	if err != nil {
		return fmt.Errorf("cannot marshal setup data: %w", err)
	}

	if err := os.WriteFile(i.setupDataPath(), file, 0o644); err != nil {
		return fmt.Errorf("cannot write setup data file: %w", err)
	}

	return nil
}
