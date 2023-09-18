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

	infoSetting := integration.SetupDataSchemaSettings{
		Id: "info",
		Label: integration.LanguageText{
			En: "Integration",
		},
		Field: integration.SettingTypeLabel{
			Label: integration.SettingTypeLabelDefinition{
				Value: integration.LanguageText{
					En: "Hi",
				},
			},
		},
	}

	inputSetting := integration.SetupDataSchemaSettings{
		Id: "data",
		Label: integration.LanguageText{
			En: "We need some data",
		},
		Field: integration.SettingTypeText{
			Text: integration.SettingTypeTextDefinition{
				Value: "some preset text",
			},
		},
	}

	setupdataschema := integration.SetupDataSchema{
		Title: integration.LanguageText{
			En: "Integration Settings",
		},
		Settings: []integration.SetupDataSchemaSettings{infoSetting, inputSetting},
	}

	metadata := integration.DriverMetadata{
		DriverId: "myintegration",
		Developer: integration.Developer{
			Name: "Sebastian Plattner",
		},
		Name: integration.LanguageText{
			En: "My UCRT Integration",
			De: "Meine UCRT Integration",
		},
		Version:         "0.0.1",
		SetupDataSchema: setupdataschema,
	}

	c.IntegrationDriver.SetMetadata(&metadata)

	// Some dummy test data
	button := entities.NewButtonEntity("mybutton", entities.LanguageText{En: "My Button", De: "Mein Button"}, "")
	button.AddCommand(entities.PushButtonEntityCommand, c.HandleButtonPressCommand)
	c.IntegrationDriver.AddEntity(button)

	// Pass function to the integration driver that is called when the remote want to setup the driver
	c.IntegrationDriver.SetHandleSetupFunction(c.HandleSetup)
	// Pass function to the integration driver that is called when the remote want to setup the driver
	c.IntegrationDriver.SetHandleConnectionFunction(c.HandleConnection)
	c.IntegrationDriver.SetHandleSetDriverUserDataFunction(c.HandleSetDriverUserDataFunction)

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

func (c *Client) HandleSetup(setup_data map[string]string) {

	log.WithField("SetupData", setup_data).Info("Handle setup_driver request in client")

	var userAction = integration.RequireUserAction{
		Confirmation: integration.ConfirmationPage{
			Title: integration.LanguageText{
				En: "You are about to add this integration. Just confirm it",
			},
		},
	}

	//event_type: SETUP with state: SETUP is a progress event to keep the process running,
	// If the setup process takes more than a few seconds,
	// the integration should send driver_setup_change events with state: SETUP to the Remote Two
	// to show a setup progress to the user and prevent an inactivity timeout.
	c.IntegrationDriver.SetDriverSetupState(integration.SetupEvent, integration.SetupState, "", nil)
	time.Sleep(1 * time.Second)

	// Start the setup
	c.IntegrationDriver.SetDriverSetupState(integration.SetupEvent, integration.WaitUserActionState, "", &userAction)

	// // For Testing, just Wait a bit
	// time.Sleep(1 * time.Second)

	// // Finish the setup
	// c.IntegrationDriver.SetDriverSetupState(integration.StopEvent, integration.OkState, "", nil)

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

func (c *Client) HandleButtonPressCommand(button entities.ButtonEntity) {
	log.Println("Button " + button.Id + "pressed")
}

func (c *Client) connect() {
	go c.clientLoop()
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
