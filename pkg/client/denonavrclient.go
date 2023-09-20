package client

import (
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/splattner/goucrt/pkg/denonavr"
	"github.com/splattner/goucrt/pkg/entities"
	"github.com/splattner/goucrt/pkg/integration"
)

// Denon AVR Client Implementation
type DenonAVRClient struct {
	Client
	denon *denonavr.DenonAVR

	testButton   *entities.ButtonEntity
	volumeSensor *entities.SensorEntity
	mediaPlayer  *entities.MediaPlayerEntity
}

func NewDenonAVRClient(i *integration.Integration) *DenonAVRClient {
	client := DenonAVRClient{}

	client.IntegrationDriver = i
	// Start without a connection
	client.DeviceState = integration.DisconnectedDeviceState

	client.messages = make(chan string)
	client.setupData = make(map[string]string)

	inputSetting := integration.SetupDataSchemaSettings{
		Id: "ipaddr",
		Label: integration.LanguageText{
			En: "IP Address",
		},
		Field: integration.SettingTypeText{
			Text: integration.SettingTypeTextDefinition{
				Value: "192.168.10.153",
			},
		},
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
		Version: "0.0.1",
		SetupDataSchema: integration.SetupDataSchema{
			Title: integration.LanguageText{
				En: "Integration Settings",
			},
			Settings: []integration.SetupDataSchemaSettings{inputSetting},
		},
	}

	client.IntegrationDriver.SetMetadata(&metadata)

	// set the client specific functions
	client.initFunc = client.initDenonAVRClient
	client.setupFunc = client.denonHandleSetup
	client.clientLoopFunc = client.denonClientLoop

	return &client
}

func (c *DenonAVRClient) initDenonAVRClient() {
	// Some dummy test data
	c.testButton = entities.NewButtonEntity("mybutton", entities.LanguageText{En: "My Button", De: "Mein Button"}, "")
	c.testButton.AddCommand(entities.PushButtonEntityCommand, c.HandleButtonPressCommand)
	c.IntegrationDriver.AddEntity(c.testButton)

	// Volume Sensor
	c.volumeSensor = entities.NewSensorEntity("mastervolume", entities.LanguageText{En: "Master Volume", De: "Master Volume"}, "", entities.CustomSensorDeviceClass)
	c.IntegrationDriver.AddEntity(c.volumeSensor)

	// Media Player
	c.mediaPlayer = entities.NewMediaPlayerEntity("mediaplayer", entities.LanguageText{En: "Denon AVR"}, "", entities.ReceiverMediaPlayerDeviceClass)
	c.mediaPlayer.AddFeature(entities.OnOffMediaPlayerEntityFeatures)
	c.mediaPlayer.AddFeature(entities.VolumeMediaPlayerEntityyFeatures)
	c.IntegrationDriver.AddEntity(c.mediaPlayer)

	// Pass function to the integration driver that is called when the remote want to setup the driver
	c.IntegrationDriver.SetHandleSetupFunction(c.HandleSetup)
	// Pass function to the integration driver that is called when the remote want to setup the driver
	c.IntegrationDriver.SetHandleConnectionFunction(c.HandleConnection)
	c.IntegrationDriver.SetHandleSetDriverUserDataFunction(c.HandleSetDriverUserDataFunction)
}

func (c *DenonAVRClient) denonHandleSetup() {

	log.WithField("SetupData", c.setupData).Info("Handle setup_driver request in client")

	//event_type: SETUP with state: SETUP is a progress event to keep the process running,
	// If the setup process takes more than a few seconds,
	// the integration should send driver_setup_change events with state: SETUP to the Remote Two
	// to show a setup progress to the user and prevent an inactivity timeout.
	c.IntegrationDriver.SetDriverSetupState(integration.SetupEvent, integration.SetupState, "", nil)
	time.Sleep(1 * time.Second)

	// var userAction = integration.RequireUserAction{
	// 	Confirmation: integration.ConfirmationPage{
	// 		Title: integration.LanguageText{
	// 			En: "You are about to add this integration. Just confirm it",
	// 		},
	// 	},
	// }

	// Start the setup with some require user data
	//c.IntegrationDriver.SetDriverSetupState(integration.SetupEvent, integration.WaitUserActionState, "", &userAction)

	// // Finish the setup
	c.IntegrationDriver.SetDriverSetupState(integration.StopEvent, integration.OkState, "", nil)

}

func (c *DenonAVRClient) setupDenon() {
	if c.denon == nil {
		if c.setupData != nil && c.setupData["ipaddr"] != "" {
			c.denon = denonavr.NewDenonAVR(c.setupData["ipaddr"])
		} else {
			log.Error("Cannot setup Denon, missing setupData")
		}
	}
}

func (c *DenonAVRClient) configureDenon() {

	if c.denon != nil {
		c.denon.AddHandleEntityChangeFunc("MasterVolume", func(value string) {
			attributes := make(map[string]interface{})

			attributes[entities.ValueSensortEntityyAttribute] = value
			attributes[entities.UnitSSensorntityyAttribute] = "db"

			c.volumeSensor.SetAttributes(attributes)
		})

		c.denon.AddHandleEntityChangeFunc("MasterVolume", func(value string) {
			attributes := make(map[string]interface{})

			var volume float64
			if s, err := strconv.ParseFloat(value, 64); err == nil {
				volume = s
			}

			attributes[entities.VolumeMediaPlayerEntityAttribute] = volume + 80

			c.mediaPlayer.SetAttributes(attributes)
		})
	}
}

func (c *DenonAVRClient) denonClientLoop() {

	defer func() {
		c.setDeviceState(integration.DisconnectedDeviceState)
	}()

	if c.denon == nil {
		// Initialize Denon Client
		c.setupDenon()

	}

	// Start the Denon Liste Loop if already configured
	if c.denon != nil {

		// Configure Denon Client
		c.configureDenon()

		log.WithFields(log.Fields{
			"Denon IP": c.denon.Host}).Info("Start Denon Client Loop")
		go c.denon.StartListenLoop()
	}

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
