package client

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/splattner/goucrt/pkg/integration"
)

// Generic client
type Client struct {
	IntegrationDriver *integration.Integration

	DeviceState integration.DState

	messages chan string

	// Client specific functions
	initFunc       func()
	setupFunc      func()
	clientLoopFunc func()
}

func NewClient(i *integration.Integration) *Client {

	client := Client{}

	client.IntegrationDriver = i
	// Start without a connection
	client.DeviceState = integration.DisconnectedDeviceState

	client.messages = make(chan string)

	return &client

}

func (c *Client) InitClient() {
	// Call setup Function if its set
	if c.initFunc != nil {
		c.initFunc()
	}
}

func (c *Client) HandleConnection(e *integration.ConnectEvent) {
	log.Println("Client, Handle connection")
	switch e.Msg {
	case "connect":
		// Only start connecting if in disconnected state or error state
		if c.DeviceState == integration.DisconnectedDeviceState || c.DeviceState == integration.ErrorDeviceState {
			log.Info("Start connecting")
			c.setDeviceState(integration.ConnectingDeviceState)

			// to make sure event is sent
			time.Sleep(1 * time.Second)

			// And then connect
			go c.connect()
		} else {
			// Just send the current state
			c.setDeviceState(c.DeviceState)
		}

	case "disconnect":

		if c.DeviceState == integration.ConnectedDeviceState {
			log.Info("Disconnecting")

			// And disconnect
			go c.disconnect()
		}

	}
}

func (c *Client) HandleSetup(setup_data integration.SetupData) {

	if c.setupFunc != nil {
		c.setupFunc()
	}

}

// User input result of a SettingsPage as key values.
// key: id of the field
// value: entered user value as string. This is either the entered text or number, selected checkbox state or the selected dropdown item id.
// ⚠️ Non native string values as numbers or booleans are represented as string values!
func (c *Client) HandleSetDriverUserDataFunction(userdata map[string]string, confirm bool) {

	log.WithFields(log.Fields{
		"Userdata": userdata,
		"Confim":   confirm,
	}).Debug(("Handle SetDriverUserData"))

	if confirm {
		c.IntegrationDriver.SetDriverSetupState(integration.StopEvent, integration.OkState, "", nil)
	} else {
		c.IntegrationDriver.SetDriverSetupState(integration.StopEvent, integration.OkState, "", nil)
		// Confirm is not set.. Bug?
		//c.IntegrationDriver.SetDriverSetupState(integration.SetupEvent, integration.WaitUserActionState, "", nil)
	}

}

func (c *Client) connect() {
	c.clientLoop()
}

func (c *Client) disconnect() {
	c.messages <- "disconnect"
}

func (c *Client) setDeviceState(state integration.DState) {
	log.Println("Set device state and send to integration driver: " + state)
	c.DeviceState = state
	c.IntegrationDriver.SetDeviceState(c.DeviceState)
}

func (c *Client) clientLoop() {
	log.Debug("Start Client Loop")

	if c.clientLoopFunc != nil {
		go c.clientLoopFunc()
	}

}
