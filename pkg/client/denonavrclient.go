package client

import (
	"strconv"
	"strings"
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

	inputSetting := integration.SetupDataSchemaSettings{
		Id: "ipaddr",
		Label: integration.LanguageText{
			En: "IP Address of your Denon Receiver",
		},
		Field: integration.SettingTypeText{
			Text: integration.SettingTypeTextDefinition{
				Value: "192.168.10.153",
			},
		},
	}

	metadata := integration.DriverMetadata{
		DriverId: "denonavr",
		Developer: integration.Developer{
			Name: "Sebastian Plattner",
		},
		Name: integration.LanguageText{
			En: "Denon AVR",
		},
		Version: "0.0.1",
		SetupDataSchema: integration.SetupDataSchema{
			Title: integration.LanguageText{
				En: "Configuration",
				De: "KOnfiguration",
			},
			Settings: []integration.SetupDataSchemaSettings{inputSetting},
		},
		Icon: "custom:denon.png",
	}

	client.IntegrationDriver.SetMetadata(&metadata)

	// set the client specific functions
	client.initFunc = client.initDenonAVRClient
	client.setupFunc = client.denonHandleSetup
	client.clientLoopFunc = client.denonClientLoop

	return &client
}

func (c *DenonAVRClient) initDenonAVRClient() {

	// Media Player
	c.mediaPlayer = entities.NewMediaPlayerEntity("mediaplayer", entities.LanguageText{En: "Denon AVR"}, "", entities.ReceiverMediaPlayerDeviceClass)
	c.mediaPlayer.AddFeature(entities.OnOffMediaPlayerEntityFeatures)
	c.mediaPlayer.AddFeature(entities.ToggleMediaPlayerEntityyFeatures)
	c.mediaPlayer.AddFeature(entities.VolumeMediaPlayerEntityyFeatures)
	c.mediaPlayer.AddFeature(entities.MuteMediaPlayerEntityFeatures)
	c.mediaPlayer.AddFeature(entities.UnmuteMediaPlayerEntityFeatures)
	c.mediaPlayer.AddFeature(entities.MuteToggleMediaPlayerEntityFeatures)
	c.mediaPlayer.AddFeature(entities.SelectSourceMediaPlayerEntityFeatures)
	c.IntegrationDriver.AddEntity(c.mediaPlayer)

	// Pass function to the integration driver that is called when the remote want to setup the driver
	c.IntegrationDriver.SetHandleSetupFunction(c.HandleSetup)
	// Pass function to the integration driver that is called when the remote want to setup the driver
	c.IntegrationDriver.SetHandleConnectionFunction(c.HandleConnection)
	c.IntegrationDriver.SetHandleSetDriverUserDataFunction(c.HandleSetDriverUserDataFunction)
}

func (c *DenonAVRClient) denonHandleSetup() {

	log.WithField("SetupData", c.IntegrationDriver.SetupData).Info("Handle setup_driver request in client")

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
		if c.IntegrationDriver.SetupData != nil && c.IntegrationDriver.SetupData["ipaddr"] != "" {
			c.denon = denonavr.NewDenonAVR(c.IntegrationDriver.SetupData["ipaddr"])
		} else {
			log.Error("Cannot setup Denon, missing setupData")
		}
	}
}

func (c *DenonAVRClient) configureDenon() {

	if c.denon != nil {
		// Configure the Entity Change Func
		c.denon.AddHandleEntityChangeFunc("MasterVolume", func(value interface{}) {
			attributes := make(map[string]interface{})

			attributes[entities.ValueSensortEntityyAttribute] = value.(string)
			attributes[entities.UnitSSensorntityyAttribute] = "db"

			c.volumeSensor.SetAttributes(attributes)
		})

		c.denon.AddHandleEntityChangeFunc("MasterVolume", func(value interface{}) {
			attributes := make(map[string]interface{})

			var volume float64
			if s, err := strconv.ParseFloat(value.(string), 64); err == nil {
				volume = s
			}

			attributes[entities.VolumeMediaPlayerEntityAttribute] = volume + 80

			c.mediaPlayer.SetAttributes(attributes)
		})

		c.denon.AddHandleEntityChangeFunc("ZonePower", func(value interface{}) {

			attributes := make(map[string]interface{})

			switch value.(string) {
			case "ON":
				attributes[string(entities.StateMediaPlayerEntityAttribute)] = entities.OnMediaPlayerEntityState
			case "OFF":
				attributes[string(entities.StateMediaPlayerEntityAttribute)] = entities.OffMediaPlayerEntityState
			}

			c.mediaPlayer.SetAttributes(attributes)
		})

		c.denon.AddHandleEntityChangeFunc("Mute", func(value interface{}) {

			attributes := make(map[string]interface{})

			switch value.(string) {
			case "on":
				attributes[string(entities.MutedMediaPlayeEntityAttribute)] = true
			case "off":
				attributes[string(entities.MutedMediaPlayeEntityAttribute)] = false
			}

			c.mediaPlayer.SetAttributes(attributes)
		})

		c.denon.AddHandleEntityChangeFunc("VideoSelectLists", func(value interface{}) {

			attributes := make(map[string]interface{})

			videoSelectList := value.([]denonavr.ValueLists)

			var tabelEntries []string

			for _, i := range videoSelectList {
				if i.Index == "ON" || i.Index == "OFF" {
					continue
				}
				tabelEntries = append(tabelEntries, strings.TrimRight(i.Table, " "))
			}

			attributes["source_list"] = tabelEntries

			c.mediaPlayer.SetAttributes(attributes)

		})

		c.denon.AddHandleEntityChangeFunc("VideoSelect", func(value interface{}) {

			attributes := make(map[string]interface{})

			videoSelect := value.(string)

			videoSelectList := make(map[string]string)

			for _, i := range c.denon.GetVideoSelectList() {
				videoSelectList[i.Index] = strings.TrimRight(i.Table, " ")
			}

			attributes["source"] = videoSelectList[videoSelect]

			c.mediaPlayer.SetAttributes(attributes)

		})

		// Add Commands
		c.mediaPlayer.AddCommand(entities.OnMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) {
			log.WithField("entityId", mediaPlayer.Id).Info("OnMediaPlayerEntityCommand called")
			c.denon.TurnOn()

		})

		c.mediaPlayer.AddCommand(entities.OffMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) {
			log.WithField("entityId", mediaPlayer.Id).Info("OffMediaPlayerEntityCommand called")
			c.denon.TurnOff()

		})

		c.mediaPlayer.AddCommand(entities.ToggleMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) {
			log.WithField("entityId", mediaPlayer.Id).Info("ToggleMediaPlayerEntityCommand called")

			c.denon.TogglePower()

		})

		c.mediaPlayer.AddCommand(entities.VolumeMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) {
			log.WithField("entityId", mediaPlayer.Id).Info("VolumeMediaPlayerEntityCommand called")

			var volume float64
			if v, err := strconv.ParseFloat(params["volume"].(string), 64); err == nil {
				volume = v
			}
			c.denon.SetVolume(volume)
		})

		c.mediaPlayer.AddCommand(entities.VolumeUpMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) {
			log.WithField("entityId", mediaPlayer.Id).Info("VolumeUpMediaPlayerEntityCommand called")
			c.denon.SetVolumeUp()
		})

		c.mediaPlayer.AddCommand(entities.VolumeDownMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) {
			log.WithField("entityId", mediaPlayer.Id).Info("VolumeDownMediaPlayerEntityCommand called")
			c.denon.SetVolumeDown()
		})

		c.mediaPlayer.AddCommand(entities.MuteMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) {
			log.WithField("entityId", mediaPlayer.Id).Info("MuteMediaPlayerEntityCommand called")
			c.denon.Mute()
		})

		c.mediaPlayer.AddCommand(entities.UnmuteMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) {
			log.WithField("entityId", mediaPlayer.Id).Info("UnmuteMediaPlayerEntityCommand called")
			c.denon.UnMute()
		})

		c.mediaPlayer.AddCommand(entities.MuteToggleMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) {
			log.WithField("entityId", mediaPlayer.Id).Info("MuteToggleMediaPlayerEntityCommand called")
			c.denon.MuteToggle()
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
