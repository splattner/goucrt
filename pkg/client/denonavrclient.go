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

	log.Debug("Initialize DenonAVR CLient")

	// Media Player
	c.mediaPlayer = entities.NewMediaPlayerEntity("mediaplayer", entities.LanguageText{En: "Denon AVR"}, "", entities.ReceiverMediaPlayerDeviceClass)
	c.mediaPlayer.AddFeature(entities.OnOffMediaPlayerEntityFeatures)
	c.mediaPlayer.AddFeature(entities.ToggleMediaPlayerEntityyFeatures)
	c.mediaPlayer.AddFeature(entities.VolumeMediaPlayerEntityyFeatures)
	c.mediaPlayer.AddFeature(entities.VolumeUpDownMediaPlayerEntityFeatures)
	c.mediaPlayer.AddFeature(entities.MuteMediaPlayerEntityFeatures)
	c.mediaPlayer.AddFeature(entities.UnmuteMediaPlayerEntityFeatures)
	c.mediaPlayer.AddFeature(entities.MuteToggleMediaPlayerEntityFeatures)
	c.mediaPlayer.AddFeature(entities.SelectSourceMediaPlayerEntityFeatures)
	c.mediaPlayer.AddFeature(entities.SelectSoundModeMediaPlayerEntityFeatures)
	c.mediaPlayer.AddFeature(entities.DPadMediaPlayerEntityFeatures)

	if err := c.IntegrationDriver.AddEntity(c.mediaPlayer); err != nil {
		log.WithError(err).Error("Cannot add Entity")
	}

	// Butons
	c.moni1Button = entities.NewButtonEntity("moni1", entities.LanguageText{En: "Monitor Out 1"}, "")
	if err := c.IntegrationDriver.AddEntity(c.moni1Button); err != nil {
		log.WithError(err).Error("Cannot add Entity")
	}

	c.moni2Button = entities.NewButtonEntity("moni2", entities.LanguageText{En: "Monitor Out 2"}, "")
	if err := c.IntegrationDriver.AddEntity(c.moni2Button); err != nil {
		log.WithError(err).Error("Cannot add Entity")
	}

	c.moniAutoButton = entities.NewButtonEntity("moniauto", entities.LanguageText{En: "Monitor Out Auto"}, "")
	if err := c.IntegrationDriver.AddEntity(c.moniAutoButton); err != nil {
		log.WithError(err).Error("Cannot add Entity")
	}

}

func (c *DenonAVRClient) denonHandleSetup(setup_data integration.SetupData) {
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
		c.moni1Button.MapCommand(entities.PushButtonEntityCommand, c.denon.SetMoni1Out)
		c.moni2Button.MapCommand(entities.PushButtonEntityCommand, c.denon.SetMoni2Out)
		c.moniAutoButton.MapCommand(entities.PushButtonEntityCommand, c.denon.SetMoniAutoOut)

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

			attributes[string(entities.VolumeMediaPlayerEntityAttribute)] = volume + 80

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

		// Media Title
		c.denon.AddHandleEntityChangeFunc("media_title", func(value interface{}) {
			attributes := make(map[string]interface{})
			attributes["media_title"] = value.(string)
			c.mediaPlayer.SetAttributes(attributes)
		})

		// Media Image URL
		c.denon.AddHandleEntityChangeFunc("media_image_url", func(value interface{}) {
			attributes := make(map[string]interface{})
			attributes["media_image_url"] = value.(string)
			c.mediaPlayer.SetAttributes(attributes)
		})

		// Add Commands
		c.mediaPlayer.MapCommand(entities.OnMediaPlayerEntityCommand, c.denon.TurnOn)
		c.mediaPlayer.MapCommand(entities.OffMediaPlayerEntityCommand, c.denon.TurnOff)
		c.mediaPlayer.MapCommand(entities.ToggleMediaPlayerEntityCommand, c.denon.TogglePower)

		c.mediaPlayer.AddCommand(entities.VolumeMediaPlayerEntityCommand, func(mediaPlayer entities.MediaPlayerEntity, params map[string]interface{}) int {
			log.WithField("entityId", mediaPlayer.Id).Debug("VolumeMediaPlayerEntityCommand called")

			var volume float64
			if v, err := strconv.ParseFloat(params["volume"].(string), 64); err == nil {
				volume = v
			}
			if err := c.denon.SetVolume(volume); err != nil {
				return 404
			}
			return 200
		})

		// Volume commands
		c.mediaPlayer.MapCommand(entities.VolumeUpMediaPlayerEntityCommand, c.denon.SetVolumeUp)
		c.mediaPlayer.MapCommand(entities.VolumeDownMediaPlayerEntityCommand, c.denon.SetVolumeDown)
		c.mediaPlayer.MapCommand(entities.MuteMediaPlayerEntityCommand, c.denon.MainZoneMute)
		c.mediaPlayer.MapCommand(entities.UnmuteMediaPlayerEntityCommand, c.denon.MainZoneUnMute)
		c.mediaPlayer.MapCommand(entities.MuteToggleMediaPlayerEntityCommand, c.denon.MainZoneMuteToggle)

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

		// Handle connection to device this integration shall control
		// Set Device state to connected when connection is established
		c.setDeviceState(integration.ConnectedDeviceState)
	}

	// Run Client Loop to handle entity changes from device
	for {
		msg := <-c.messages
		switch msg {
		case "disconnect":
			return
		}
	}

}
