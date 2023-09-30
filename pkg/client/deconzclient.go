package client

import (
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/splattner/goucrt/pkg/deconz"
	"github.com/splattner/goucrt/pkg/entities"
	"github.com/splattner/goucrt/pkg/integration"
)

var mapOnState = map[bool]entities.LightEntityState{
	true:  entities.OnLightEntityState,
	false: entities.OffLightEntityState,
}

// Denon AVR Client Implementation
type DeconzClient struct {
	Client
	deconz *deconz.Deconz
}

func NewDeconzClient(i *integration.Integration) *DeconzClient {
	client := DeconzClient{}

	client.IntegrationDriver = i
	// Start without a connection
	client.DeviceState = integration.DisconnectedDeviceState

	client.messages = make(chan string)

	ipaddr := integration.SetupDataSchemaSettings{
		Id: "ipaddr",
		Label: integration.LanguageText{
			En: "IP Address of your deCONZ Client",
		},
		Field: integration.SettingTypeText{
			Text: integration.SettingTypeTextDefinition{
				Value: "",
			},
		},
	}

	port := integration.SetupDataSchemaSettings{
		Id: "port",
		Label: integration.LanguageText{
			En: "Port used by your deCONZ CLient",
		},
		Field: integration.SettingTypeText{
			Text: integration.SettingTypeTextDefinition{
				Value: "8080",
			},
		},
	}

	websocketport := integration.SetupDataSchemaSettings{
		Id: "websocketport",
		Label: integration.LanguageText{
			En: "Websocket Port used by your deCONZ CLient",
		},
		Field: integration.SettingTypeText{
			Text: integration.SettingTypeTextDefinition{
				Value: "8081",
			},
		},
	}

	metadata := integration.DriverMetadata{
		DriverId: "deCONZ",
		Developer: integration.Developer{
			Name: "Sebastian Plattner",
		},
		Name: integration.LanguageText{
			En: "DeCONZ",
		},
		Version: "0.0.1",
		SetupDataSchema: integration.SetupDataSchema{
			Title: integration.LanguageText{
				En: "Configuration",
				De: "Konfiguration",
			},
			Settings: []integration.SetupDataSchemaSettings{ipaddr, port, websocketport},
		},
		Icon: "custom:deconz.png",
	}

	client.IntegrationDriver.SetMetadata(&metadata)

	// set the client specific functions
	client.initFunc = client.initDeconzClient
	client.setupFunc = client.deconzHandleSetup
	client.clientLoopFunc = client.deconzClientLoop
	client.setDriverUserDataFunc = client.handleSetDriverUserData

	return &client
}

func (c *DeconzClient) handleSetDriverUserData(user_data map[string]string, confirm bool) {

	log.Debug("Deconz handle set driver user data")

	// confirm seems to be set to false always, maybe just the presence of the field tells me,
	// confirmation was sent?
	if len(user_data) == 0 {
		// Get a new Denon API Key

		ipaddr := c.IntegrationDriver.SetupData["ipaddr"]
		port, _ := strconv.Atoi(c.IntegrationDriver.SetupData["port"])
		websocketport, _ := strconv.Atoi(c.IntegrationDriver.SetupData["websocketport"])

		deconz := deconz.NewDeconz(ipaddr, port, websocketport, "")
		apikey, err := deconz.GetNewAPIKey(c.IntegrationDriver.DriverId)

		if err != nil {
			log.WithError(err).Debug("Failed to get new api Key")
			c.IntegrationDriver.SetDriverSetupState(integration.StopEvent, integration.ErrorState, integration.AuthErrorError, nil)
			return
		}

		c.IntegrationDriver.SetupData["apikey"] = apikey
		c.IntegrationDriver.PersistSetupData()

		c.IntegrationDriver.SetDriverSetupState(integration.StopEvent, integration.OkState, "", nil)

	}
}

func (c *DeconzClient) initDeconzClient() {

}

func (c *DeconzClient) deconzHandleSetup(setup_data integration.SetupData) {
	//event_type: SETUP with state: SETUP is a progress event to keep the process running,
	// If the setup process takes more than a few seconds,
	// the integration should send driver_setup_change events with state: SETUP to the Remote Two
	// to show a setup progress to the user and prevent an inactivity timeout.
	//c.IntegrationDriver.SetDriverSetupState(integration.SetupEvent, integration.SetupState, "", nil)
	time.Sleep(1 * time.Second)

	var userAction = integration.RequireUserAction{
		Confirmation: integration.ConfirmationPage{
			Title: integration.LanguageText{
				En: "Gateway configuration",
			},
			Message1: integration.LanguageText{
				En: "Please unlock your DeCONZ Gateway to create a new API Key",
			},
		},
	}

	// Start the setup with some require user data
	c.IntegrationDriver.SetDriverSetupState(integration.SetupEvent, integration.WaitUserActionState, "", &userAction)

	// // Finish the setup
	//c.IntegrationDriver.SetDriverSetupState(integration.StopEvent, integration.OkState, "", nil)

}

func (c *DeconzClient) setupDeconz() {

	if c.IntegrationDriver.SetupData["apikey"] != "" {

		ipaddr := c.IntegrationDriver.SetupData["ipaddr"]
		port, _ := strconv.Atoi(c.IntegrationDriver.SetupData["port"])
		websocketport, _ := strconv.Atoi(c.IntegrationDriver.SetupData["websocketport"])

		log.WithFields(log.Fields{
			"ipaddr":        ipaddr,
			"port":          port,
			"websocketport": websocketport,
		}).Debug("Create DeCONZ Client")

		deconz := deconz.NewDeconz(ipaddr, port, websocketport, c.IntegrationDriver.SetupData["apikey"])
		c.deconz = deconz
	}

}

func (c *DeconzClient) configureDeconz() {

	log.Debug("Configure DeCONZ")

	c.deconz.SetDeviceDiscoveredHandler(c.handleNewDeviceDiscovered)
	c.deconz.SetDeviceRemoveHandler(c.handleRemoveDevice)

	// TODO, enable groups as setup_data
	c.deconz.StartDiscovery(true)

}

func (c *DeconzClient) handleNewSensorDeviceDiscovered(device *deconz.DeconzDevice) {

	var sensor *entities.SensorEntity
	if device.Sensor.State.Temperature != nil {
		sensor = entities.NewSensorEntity(fmt.Sprintf("sensor%d", device.GetID()), entities.LanguageText{En: device.GetName()}, "", entities.TemperaturSensorDeviceClass)
	}

	if device.Sensor.State.Humidity != nil {
		sensor = entities.NewSensorEntity(fmt.Sprintf("sensor%d", device.GetID()), entities.LanguageText{En: device.GetName()}, "", entities.HumiditySensorDeviceClass)
	}

	// Currently no other sensors are implemeted

	if sensor != nil {

		device.SetHandleChangeStateFunc(func(state *deconz.DeconzState) {
			log.WithFields(log.Fields{
				"ID":    device.GetID(),
				"State": state,
			}).Trace("Sensor changed")

			attributes := make(map[string]interface{})

			switch sensor.DeviceClass {
			case entities.TemperaturSensorDeviceClass:
				if state.Temperature != nil {
					attributes["value"] = *state.Temperature / int16(100.0)
				}

			case entities.HumiditySensorDeviceClass:
				if state.Humidity != nil {
					attributes["value"] = *state.Humidity / uint16(100)
				}

			}

			if attributes["value"] != nil {
				sensor.SetAttributes(attributes)
			}

		})

		if err := c.IntegrationDriver.AddEntity(sensor); err != nil {
			log.WithError(err).Error("Cannot add entity")
		}
	}
}

func (c *DeconzClient) handleNewLightDeviceDiscovered(device *deconz.DeconzDevice) {
	light := entities.NewLightEntity(fmt.Sprintf("light%d", device.GetID()), entities.LanguageText{En: device.GetName()}, "")

	// Add Features and initial values
	light.AddFeature(entities.OnOffLightEntityFeatures)
	light.AddFeature(entities.ToggleLightEntityFeatures)
	light.UpdateAttribute(entities.StateLightEntityAttribute, mapOnState[device.IsOn()])

	light.AddFeature(entities.DimLightEntityFeatures)
	light.UpdateAttribute(entities.BrightnessLightEntityAttribute, device.GetBrightness())

	if device.HasColor() {
		switch device.Light.State.ColorMode {
		case "ct":
			light.AddFeature(entities.ColorTemperatureLightEntityFeatures)
			light.UpdateAttribute(entities.ColorTemperatureLightEntityAttribute, device.GetColorTempInPercent())
		case "hs":
			light.AddFeature(entities.ColorLightEntityFeatures)
			light.UpdateAttribute(entities.HueLightEntityAttribute, device.GetHueConverted())
			light.UpdateAttribute(entities.SaturationLightEntityAttribute, device.GetSaturation())
		}
	}

	// Set initial attribute

	// Add commands
	light.AddCommand(entities.OnLightEntityCommand, func(entity entities.LightEntity, params map[string]interface{}) int {

		// NO param set, so just turn on
		if len(params) == 0 {
			if err := device.TurnOn(); err != nil {
				return 404
			}
		} else {

			if params["brightness"] != nil {

				if err := device.SetBrightness(float32(params["brightness"].(uint))); err != nil {
					return 404
				}
			}

			if params["hue"] != nil {
				hue := params["hue"].(float64) / 360 * 65535
				if err := device.SetHue(float32(hue)); err != nil {
					return 404
				}
			}

			if params["saturation"] != nil {
				if err := device.SetSaturation(float32(params["saturation"].(uint))); err != nil {
					return 404
				}
			}

			if params["color_temperature"] != nil {

				raw_ct := params["color_temperature"].(float64)
				ct := raw_ct/100*(500-153) + 153

				if err := device.SetColorTemp(float32(ct)); err != nil {
					return 404
				}
			}
		}

		return 200
	})

	light.MapCommand(entities.OffLightEntityCommand, device.TurnOff)
	light.MapCommand(entities.ToggleLightEntityCommand, device.Toggle)

	// Set the Handle State Change function
	device.SetHandleChangeStateFunc(func(state *deconz.DeconzState) {
		log.WithFields(log.Fields{
			"ID":   device.GetID(),
			"Type": device.Type,
		}).Debug("Handle State Change for this Light")

		attributes := make(map[string]interface{})

		if light.HasAttribute(entities.StateLightEntityAttribute) {
			if state.Reachable != nil && !*state.Reachable {
				attributes[string(entities.StateLightEntityAttribute)] = entities.UnavailableLightEntityState
			} else {
				attributes[string(entities.StateLightEntityAttribute)] = mapOnState[*state.On]
			}

		}

		if light.HasAttribute(entities.HueLightEntityAttribute) {
			if state.Hue != nil {
				attributes[string(entities.HueLightEntityAttribute)] = device.GetHueConverted()
			}
		}

		if light.HasAttribute(entities.SaturationLightEntityAttribute) {
			if state.Sat != nil {
				attributes[string(entities.SaturationLightEntityAttribute)] = *state.Sat
			}

		}

		if light.HasAttribute(entities.BrightnessLightEntityAttribute) {
			if state.Bri != nil {
				attributes[string(entities.BrightnessLightEntityAttribute)] = *state.Bri
			}

		}

		if light.HasAttribute(entities.ColorTemperatureLightEntityAttribute) {
			if state.CT != nil {
				attributes[string(entities.ColorTemperatureLightEntityAttribute)] = device.GetColorTempInPercent()
			}

		}

		light.SetAttributes(attributes)

	})

	if err := c.IntegrationDriver.AddEntity(light); err != nil {
		log.WithError(err).Error("Cannot add entity")
	}
}

func (c *DeconzClient) handleNewGroupDeviceDiscovered(device *deconz.DeconzDevice) {
	group := entities.NewLightEntity(fmt.Sprintf("group%d", device.GetID()), entities.LanguageText{En: device.GetName()}, "")

	// Add Features and initial values
	group.AddFeature(entities.OnOffLightEntityFeatures)
	group.AddFeature(entities.ToggleLightEntityFeatures)
	group.UpdateAttribute(entities.StateLightEntityAttribute, mapOnState[device.IsOn()])

	group.AddFeature(entities.DimLightEntityFeatures)
	group.UpdateAttribute(entities.BrightnessLightEntityAttribute, device.GetBrightness())

	if device.HasColor() {
		for _, light := range device.Group.Lights {
			if light.HasColor() {
				switch light.Light.State.ColorMode {
				case "ct":
					group.AddFeature(entities.ColorTemperatureLightEntityFeatures)
					group.UpdateAttribute(entities.ColorTemperatureLightEntityAttribute, light.GetColorTempInPercent())
				case "hs":
					group.AddFeature(entities.ColorLightEntityFeatures)
					group.UpdateAttribute(entities.HueLightEntityAttribute, light.GetHueConverted())
					group.UpdateAttribute(entities.SaturationLightEntityAttribute, light.GetSaturation())
				}
			}
		}

	}

	// Commands
	group.AddCommand(entities.OnLightEntityCommand, func(entity entities.LightEntity, params map[string]interface{}) int {

		// NO param set, so just turn on
		if len(params) == 0 {
			if err := device.TurnOn(); err != nil {
				return 404
			}
		} else {

			if params["brightness"] != nil {
				//bri, _ := strconv.ParseFloat(params["brightness"].(string), 32)
				if err := device.SetBrightness(float32(params["brightness"].(float64))); err != nil {
					return 404
				}
			}

			if params["hue"] != nil {
				hue := params["hue"].(float64) / 360 * 65535
				if err := device.SetHue(float32(hue)); err != nil {
					return 404
				}
			}

			if params["saturation"] != nil {
				if err := device.SetSaturation(float32(params["saturation"].(float64))); err != nil {
					return 404
				}
			}

			if params["color_temperature"] != nil {
				raw_ct := params["color_temperature"].(float64)
				ct := raw_ct/100*(500-153) + 153

				if err := device.SetColorTemp(float32(ct)); err != nil {
					return 404
				}
			}
		}
		return 200
	})

	group.AddCommand(entities.OffLightEntityCommand, func(entity entities.LightEntity, params map[string]interface{}) int {

		if err := device.TurnOff(); err != nil {
			return 404
		}
		return 200
	})

	group.AddCommand(entities.ToggleLightEntityCommand, func(entity entities.LightEntity, params map[string]interface{}) int {
		if device.IsOn() {
			if err := device.TurnOff(); err != nil {
				return 404
			}
		} else {
			if err := device.TurnOn(); err != nil {
				return 404
			}
		}
		return 200
	})

	device.SetHandleChangeStateFunc(func(state *deconz.DeconzState) {
		log.WithFields(log.Fields{
			"ID":   device.GetID(),
			"Type": device.Type,
		}).Debug("Handle State Change for this Group")

		attributes := make(map[string]interface{})

		if group.HasAttribute(entities.StateLightEntityAttribute) {
			attributes[string(entities.StateLightEntityAttribute)] = mapOnState[*device.Group.State.AnyOn]

			// if off, then just update this and leave the rest
			if !*device.Group.State.AnyOn {
				group.SetAttributes(attributes)
				return
			}

		}

		if group.HasAttribute(entities.BrightnessLightEntityAttribute) {
			if state.Bri != nil {
				attributes[string(entities.BrightnessLightEntityAttribute)] = *state.Bri
			}
		}

		if group.HasAttribute(entities.HueLightEntityAttribute) {
			if state.Hue != nil {
				attributes[string(entities.HueLightEntityAttribute)] = device.GetHueConverted()
			}
		}

		if group.HasAttribute(entities.SaturationLightEntityAttribute) {
			if state.Sat != nil {
				// Todo mapping
				attributes[string(entities.SaturationLightEntityAttribute)] = *state.Sat
			}
		}

		if group.HasAttribute(entities.ColorTemperatureLightEntityAttribute) {
			if state.CT != nil {
				attributes[string(entities.ColorTemperatureLightEntityAttribute)] = device.GetColorTempInPercent()
			}
		}

		group.SetAttributes(attributes)

	})

	if err := c.IntegrationDriver.AddEntity(group); err != nil {
		log.WithError(err).Error("Cannot add entity")
	}
}

func (c *DeconzClient) handleNewDeviceDiscovered(device *deconz.DeconzDevice) {
	log.WithFields(log.Fields{
		"id":   device.GetID(),
		"type": device.Type,
		"name": device.GetName(),
	}).Debug("New Deconz Device discovered")

	switch device.Type {
	case deconz.SensorDeconzDeviceType:
		c.handleNewSensorDeviceDiscovered(device)

	case deconz.LightDeconzDeviceType:
		c.handleNewLightDeviceDiscovered(device)

	case deconz.GroupDeconzDeviceType:
		c.handleNewGroupDeviceDiscovered(device)
	}

}

func (c *DeconzClient) handleRemoveDevice(device *deconz.DeconzDevice) {
	log.WithFields(log.Fields{
		"ID":   device.GetID(),
		"Name": device.GetName(),
		"Type": device.Type,
	}).Debug("Deconz Device not available anymore")

	switch device.Type {
	case deconz.SensorDeconzDeviceType:
		if err := c.IntegrationDriver.RemoveEntityByID(fmt.Sprintf("sensor%d", device.GetID())); err != nil {
			log.WithError(err).Error("Cannot remove Entity")
		}
	case deconz.LightDeconzDeviceType:
		if err := c.IntegrationDriver.RemoveEntityByID(fmt.Sprintf("light%d", device.GetID())); err != nil {
			log.WithError(err).Error("Cannot remove Entity")
		}
	case deconz.GroupDeconzDeviceType:
		if err := c.IntegrationDriver.RemoveEntityByID(fmt.Sprintf("group%d", device.GetID())); err != nil {
			log.WithError(err).Error("Cannot remove Entity")
		}
	}

}

// Start the Denon Listen Loop
// disconnect when finished
func (c *DeconzClient) startDenonListenLoop() {
	defer func() {
		// disconnect and let RT make a new connection again
		c.messages <- "disconnect"
	}()

	c.deconz.StartandListenLoop()
}

// Callen on RT connect
func (c *DeconzClient) deconzClientLoop() {

	ticker := time.NewTicker(5 * time.Minute)

	defer func() {
		if c.deconz != nil {
			c.deconz.Stop()
		}
		ticker.Stop()
		c.setDeviceState(integration.DisconnectedDeviceState)
	}()

	if c.deconz == nil {
		c.setupDeconz()
	}

	if c.deconz != nil {
		c.configureDeconz()

		go c.startDenonListenLoop()

	} else {
		return
	}

	// Handle connection to device this integration shall control
	// Set Device state to connected when connection is established
	c.setDeviceState(integration.ConnectedDeviceState)

	// Run Client Loop to handle entity changes from device
	for {
		select {
		case <-ticker.C:
			// Run Discovery again
			c.deconz.StartDiscovery(true)
		case msg := <-c.messages:

			switch msg {
			case "disconnect":
				return
			case "discovery":
				// Run Discovery again
				c.deconz.StartDiscovery(true)
			}
		}
	}

}
