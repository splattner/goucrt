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

	moni1Button    *entities.ButtonEntity
	moni2Button    *entities.ButtonEntity
	moniAutoButton *entities.ButtonEntity

	mediaPlayer *entities.MediaPlayerEntity
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
				Value: "",
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
	c.mediaPlayer.AddFeature(entities.SelectSoundModeMediaPlayerEntityCommand)
	c.mediaPlayer.AddFeature(entities.DPadMediaPlayerEntityFeatures)
	c.IntegrationDriver.AddEntity(c.mediaPlayer)

	// Butons
	c.moni1Button = entities.NewButtonEntity("moni1", entities.LanguageText{En: "Monitor Out 1"}, "")
	c.IntegrationDriver.AddEntity(c.moni1Button)
	c.moni2Button = entities.NewButtonEntity("moni2", entities.LanguageText{En: "Monitor Out 2"}, "")
	c.IntegrationDriver.AddEntity(c.moni2Button)
	c.moniAutoButton = entities.NewButtonEntity("moniauto", entities.LanguageText{En: "Monitor Out Auto"}, "")
	c.IntegrationDriver.AddEntity(c.moniAutoButton)

	// Pass function to the integration driver that is called when the remote want to setup the driver
	c.IntegrationDriver.SetHandleSetupFunction(c.HandleSetup)
	// Pass function to the integration driver that is called when the remote want to setup the driver
	c.IntegrationDriver.SetHandleConnectionFunction(c.HandleConnection)
	c.IntegrationDriver.SetHandleSetDriverUserDataFunction(c.HandleSetDriverUserDataFunction)
}

func (c *DenonAVRClient) denonHandleSetup() {
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

		// Buttons

		c.moni1Button.AddCommand(entities.PushButtonEntityCommand, func(button entities.ButtonEntity) int {
			return c.denon.SetMoni1Out()
		})

		c.moni2Button.AddCommand(entities.PushButtonEntityCommand, func(button entities.ButtonEntity) int {
			return c.denon.SetMoni2Out()
		})

		c.moniAutoButton.AddCommand(entities.PushButtonEntityCommand, func(button entities.ButtonEntity) int {
			return c.denon.SetMoniAutoOut()
		})

		// Media Player

		c.denon.AddHandleEntityChangeFunc("Power", func(value interface{}) {

			attributes := make(map[string]interface{})

			switch value.(string) {
			case "ON":
				attributes[string(entities.StateMediaPlayerEntityAttribute)] = entities.OnMediaPlayerEntityState
			case "OFF":
				attributes[string(entities.StateMediaPlayerEntityAttribute)] = entities.OffMediaPlayerEntityState
			}

			c.mediaPlayer.SetAttributes(attributes)
		})

		c.denon.AddHandleEntityChangeFunc("MainZoneVolume", func(value interface{}) {
			attributes := make(map[string]interface{})

			var volume float64
			if s, err := strconv.ParseFloat(value.(string), 64); err == nil {
				volume = s
			}

			attributes[entities.VolumeMediaPlayerEntityAttribute] = volume + 80

			c.mediaPlayer.SetAttributes(attributes)
		})

		c.denon.AddHandleEntityChangeFunc("MainZoneMute", func(value interface{}) {

			attributes := make(map[string]interface{})

			switch value.(string) {
			case "on":
				attributes[string(entities.MutedMediaPlayeEntityAttribute)] = true
			case "off":
				attributes[string(entities.MutedMediaPlayeEntityAttribute)] = false
			}

			c.mediaPlayer.SetAttributes(attributes)
		})

		c.denon.AddHandleEntityChangeFunc("MainZoneInputFuncList", func(value interface{}) {

			attributes := make(map[string]interface{})

			var sourceList []string
			mainZoneInputFuncSelectList := c.denon.GetZoneInputFuncList(denonavr.MainZone)
			for _, renamedSource := range mainZoneInputFuncSelectList {
				sourceList = append(sourceList, renamedSource)
			}
			attributes["source_list"] = sourceList

			c.mediaPlayer.SetAttributes(attributes)

		})

		c.denon.AddHandleEntityChangeFunc("MainZoneInputFuncSelect", func(value interface{}) {

			attributes := make(map[string]interface{})

			// We use the renamed Name
			mainZoneInputFuncSelectList := c.denon.GetZoneInputFuncList(denonavr.MainZone)

			attributes["source"] = mainZoneInputFuncSelectList[value.(string)]

			c.mediaPlayer.SetAttributes(attributes)

		})

		c.denon.AddHandleEntityChangeFunc("MainZoneSurroundMode", func(value interface{}) {

			attributes := make(map[string]interface{})

			attributes["sound_mode"] = value.(string)

			c.mediaPlayer.SetAttributes(attributes)

		})

		// We can set the sound_mode_list without change handler. Its static
		func() {

			attributes := make(map[string]interface{})

			attributes["sound_mode_list"] = c.denon.GetSoundModeList()

			c.mediaPlayer.SetAttributes(attributes)

		}()

		// Add Commands
		c.mediaPlayer.AddCommand(entities.OnMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) int {
			log.WithField("entityId", mediaPlayer.Id).Debug("OnMediaPlayerEntityCommand called")
			return c.denon.TurnOn()

		})

		c.mediaPlayer.AddCommand(entities.OffMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) int {
			log.WithField("entityId", mediaPlayer.Id).Debug("OffMediaPlayerEntityCommand called")
			return c.denon.TurnOff()

		})

		c.mediaPlayer.AddCommand(entities.ToggleMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) int {
			log.WithField("entityId", mediaPlayer.Id).Debug("ToggleMediaPlayerEntityCommand called")

			return c.denon.TogglePower()

		})

		c.mediaPlayer.AddCommand(entities.VolumeMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) int {
			log.WithField("entityId", mediaPlayer.Id).Debug("VolumeMediaPlayerEntityCommand called")

			var volume float64
			if v, err := strconv.ParseFloat(params["volume"].(string), 64); err == nil {
				volume = v
			}
			return c.denon.SetVolume(volume)
		})

		// Volume commands
		c.mediaPlayer.AddCommand(entities.VolumeUpMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) int {
			log.WithField("entityId", mediaPlayer.Id).Debug("VolumeUpMediaPlayerEntityCommand called")
			return c.denon.SetVolumeUp()
		})

		c.mediaPlayer.AddCommand(entities.VolumeDownMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) int {
			log.WithField("entityId", mediaPlayer.Id).Debug("VolumeDownMediaPlayerEntityCommand called")
			return c.denon.SetVolumeDown()
		})

		c.mediaPlayer.AddCommand(entities.MuteMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) int {
			log.WithField("entityId", mediaPlayer.Id).Debug("MuteMediaPlayerEntityCommand called")
			return c.denon.MainZoneMute()
		})

		c.mediaPlayer.AddCommand(entities.UnmuteMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) int {
			log.WithField("entityId", mediaPlayer.Id).Debug("UnmuteMediaPlayerEntityCommand called")
			return c.denon.MainZoneUnMute()
		})

		c.mediaPlayer.AddCommand(entities.MuteToggleMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) int {
			log.WithField("entityId", mediaPlayer.Id).Debug("MuteToggleMediaPlayerEntityCommand called")
			return c.denon.MainZoneMuteToggle()
		})

		// Source commands
		c.mediaPlayer.AddCommand(entities.SelectSourcMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) int {
			log.WithField("entityId", mediaPlayer.Id).Debug("SelectSourcMediaPlayerEntityCommand called")
			return c.denon.SetSelectSourceMainZone(params["source"].(string))
		})

		// Cursor commands
		c.mediaPlayer.AddCommand(entities.CursorUpMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) int {
			log.WithField("entityId", mediaPlayer.Id).Debug("CursorUpMediaPlayerEntityCommand called")
			return c.denon.CursorControl(denonavr.DenonCursorControlUp)
		})
		c.mediaPlayer.AddCommand(entities.CursorDownMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) int {
			log.WithField("entityId", mediaPlayer.Id).Debug("CursorDownMediaPlayerEntityCommand called")
			return c.denon.CursorControl(denonavr.DenonCursorControlDown)
		})
		c.mediaPlayer.AddCommand(entities.CursorLeftMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) int {
			log.WithField("entityId", mediaPlayer.Id).Debug("CursorUpMediaPlayerEntityCommand called")
			return c.denon.CursorControl(denonavr.DenonCursorControlLeft)
		})
		c.mediaPlayer.AddCommand(entities.CursorRightMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) int {
			log.WithField("entityId", mediaPlayer.Id).Debug("CursorRightMediaPlayerEntityCommand called")
			return c.denon.CursorControl(denonavr.DenonCursorControlRight)
		})
		c.mediaPlayer.AddCommand(entities.CursorEnterMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) int {
			log.WithField("entityId", mediaPlayer.Id).Debug("CursorEnterMediaPlayerEntityCommand called")
			return c.denon.CursorControl(denonavr.DenonCursorControlEnter)
		})

		// Sound Mode
		c.mediaPlayer.AddCommand(entities.SelectSoundModeMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) int {
			log.WithField("entityId", mediaPlayer.Id).Debug("SelectSoundModeMediaPlayerEntityCommand called")
			return c.denon.SetSoundModeMainZone(params["mode"].(string))
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
	} else {
		return
	}

	// Start the Denon Liste Loop if already configured
	if c.denon != nil {

		// Configure Denon Client
		c.configureDenon()

		log.WithFields(log.Fields{
			"Denon IP": c.denon.Host}).Info("Start Denon AVR Client Loop")
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
