package integration

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Generic client
type Client struct {
	IntegrationDriver *Integration

	DeviceState DState

	Messages chan string

	// Client specific functions
	// Initialize the client
	// Here you can add entities if they are already known
	InitFunc func()
	// Called by RemoteTwo when the integration is added and setup started
	SetupFunc func(SetupData)
	// Handles connect/disconnect calls from RemoteTwo
	ClientLoopFunc        func()
	SetDriverUserDataFunc func(map[string]string, bool)
}

func NewClient(i *Integration) *Client {

	client := Client{}

	client.IntegrationDriver = i
	// Start without a connection
	client.DeviceState = DisconnectedDeviceState

	client.Messages = make(chan string)

	return &client

}

func (c *Client) InitClient() {

	// Pass function to the integration driver that is called when the remote want to setup the driver
	c.IntegrationDriver.SetHandleSetupFunction(c.HandleSetup)
	// Pass function to the integration driver that is called when the remote want to connect the driver
	c.IntegrationDriver.SetHandleConnectionFunction(c.HandleConnection)
	// Pass function to the integration driver that is called when the remote want to send data from required user input page
	c.IntegrationDriver.SetHandleSetDriverUserDataFunction(c.HandleSetDriverUserDataFunction)

	// Call setup Function if its set
	if c.InitFunc != nil {
		c.InitFunc()
	}
}

func (c *Client) HandleConnection(e *ConnectEvent) {
	switch e.Msg {
	case "connect":
		// Only start connecting if in disconnected state or error state
		if c.DeviceState == DisconnectedDeviceState || c.DeviceState == ErrorDeviceState {
			c.SetDeviceState(ConnectingDeviceState)

			// to make sure event is sent
			time.Sleep(1 * time.Second)

			// And then connect
			go c.Connect()
		} else {
			// Just send the current state
			c.SetDeviceState(c.DeviceState)
		}

	case "disconnect":

		if c.DeviceState == ConnectedDeviceState {
			// And disconnect
			go c.Disconnect()
		}

	}
}

// Handle Setup called by Remote Two to setup this integration
// the SetupData are passed to this function
func (c *Client) HandleSetup(setup_data SetupData) {

	if c.SetupFunc != nil {
		c.SetupFunc(setup_data)
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

	if c.SetDriverUserDataFunc != nil {
		c.SetDriverUserDataFunc(userdata, confirm)
	}

}

func (c *Client) Connect() {
	c.ClientLoop()
}

func (c *Client) Disconnect() {
	c.Messages <- "disconnect"
}

func (c *Client) SetDeviceState(state DState) {
	log.WithField("state", state).Debug("Set device state and send to integration")
	c.DeviceState = state
	c.IntegrationDriver.SetDeviceState(c.DeviceState)
}

func (c *Client) ClientLoop() {
	log.Info("Start Client Loop")

	if c.ClientLoopFunc != nil {
		go c.ClientLoopFunc()
	} else {
		log.Fatal("Client loop not implemented")
	}

}
