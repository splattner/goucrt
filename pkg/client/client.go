package client

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/splattner/goucrt/pkg/entities"
	"github.com/splattner/goucrt/pkg/integration"
)

type Client struct {
	IntegrationDriver *integration.Integration

	DeviceState integration.DState

	messages chan string
}

func NewClient(i *integration.Integration) *Client {

	client := Client{}

	client.IntegrationDriver = i
	// Start without a connection
	client.DeviceState = integration.DisconnectedDeviceState

	client.messages = make(chan string)

	return &client

}

func (c *Client) SetupClient() {

	// Some dummy test data
	button := entities.NewButtonEntity("mybutton", entities.LanguageText{En: "My Button", De: "Mein Button"}, "")
	button.AddCommand(entities.PushButtonEntityCommand, c.HandleButtonPressCommand)
	c.IntegrationDriver.AddEntity(button)

	// Pass function to the integration driver that is called when the remote want to setup the driver
	c.IntegrationDriver.SetHandleSetupFunction(c.HandleSetup)
	// Pass function to the integration driver that is called when the remote want to setup the driver
	c.IntegrationDriver.SetHandleConnectionFunction(c.HandleConnection)

}

func (c *Client) HandleConnection(e *integration.ConnectEvent) {
	log.Println("Client, Handle connection")
	switch e.Msg {
	case "connect":
		// Only start connecting if in disconnected state or error state
		if c.DeviceState == integration.DisconnectedDeviceState || c.DeviceState == integration.ErrorDeviceState {
			log.Println("start connecting")
			c.setDeviceState(integration.ConnectingDeviceState)

			// And then connect
			go c.connect()
		} else {
			// Just send the current state
			c.setDeviceState(c.DeviceState)
		}

	case "disconnect":

		if c.DeviceState == integration.ConnectedDeviceState {
			log.Println("disconnect")

			// And disconnect
			go c.disconnect()
		}

	}
}

func (c *Client) HandleSetup() {

	var userAction = integration.RequiredUserAction{
		Confirmation: integration.ConfirmationPage{
			Title: integration.LanguageText{
				En: "You are about to add this integration",
			},
		},
	}

	// Start the setup
	c.IntegrationDriver.SetDriverSetupState(integration.StartEvent, integration.SetupState, integration.NoneError, &userAction)

	// For Testing, just Wait a bit
	time.Sleep(1 * time.Second)

	// Finish the setup
	c.IntegrationDriver.SetDriverSetupState(integration.StopEvent, integration.OkState, integration.NoneError, nil)

}

func (c *Client) HandleButtonPressCommand(button entities.ButtonEntity) {
	log.Println("Button " + button.Id + "pressed")
}

func (c *Client) connect() {
	go c.clientLoop()
}

func (c *Client) disconnect() {
	c.messages <- "dissconnect"
}

func (c *Client) setDeviceState(state integration.DState) {
	log.Println("Set device state and send to integration driver: " + state)
	c.DeviceState = state
	c.IntegrationDriver.SetDeviceState(c.DeviceState)
}

func (c *Client) clientLoop() {

	defer func() {
		c.setDeviceState(integration.DisconnectedDeviceState)
	}()

	// Handle connection to device this integration shall control
	// Set Device state to connected when connection is established
	c.setDeviceState(integration.ConnectedDeviceState)

	// Run Client Loop to handle entity changes from device
	for {

		select {
		case msg := <-c.messages:

			switch msg {
			case "disconnect":
				return
			}

		}
	}

}
