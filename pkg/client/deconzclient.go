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
		Icon: "",
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

	// confirm seens to be set to false always, maybe just the presence of the field tells be,
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
			c.IntegrationDriver.SetDriverSetupState(integration.StopEvent, integration.ErrorDeviceState, integration.AuthErrorError, nil)
			return
		}

		c.IntegrationDriver.SetupData["apikey"] = apikey
		c.IntegrationDriver.PersistSetupData()

		c.IntegrationDriver.SetDriverSetupState(integration.StopEvent, integration.OkState, "", nil)

	}
}

func (c *DeconzClient) initDeconzClient() {

}

func (c *DeconzClient) deconzHandleSetup() {
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
	go c.deconz.StartDiscovery(true)

}

func (c *DeconzClient) handleNewDeviceDiscovered(device *deconz.DeconzDevice) {
	log.WithFields(log.Fields{
		"id":   device.GetID(),
		"type": device.Type,
		"name": device.GetName(),
	}).Debug("New Deconz Device discovered")

	switch device.Type {
	case deconz.SensorDeconzDeviceType:
		//sensor := entities.NewSensorEntity(fmt.Sprintf("light%d", device.GetID()), entities.LanguageText{En: device.GetName()}, "")

		//c.IntegrationDriver.AddEntity(sensor)

	case deconz.LightDeconzDeviceType:
		light := entities.NewLightEntity(fmt.Sprintf("light%d", device.GetID()), entities.LanguageText{En: device.GetName()}, "")

		// All correct Attributes
		light.AddFeature(entities.OnOffLightEntityFeatures)
		light.AddFeature(entities.ToggleLightEntityFeatures)
		light.AddFeature(entities.DimLightEntityFeatures)

		if device.Light.HasColor {
			switch device.Light.State.ColorMode {
			case "ct":
				light.AddFeature(entities.ColorTemperatureLightEntityFeatures)
			case "hs":
				light.AddFeature(entities.ColorLightEntityFeatures)
			}
		}

		// Set initial attribute
		if light.HasAttribute(entities.StateLightEntityAttribute) {
			if *device.Light.State.On {
				light.Attributes[string(entities.StateLightEntityAttribute)] = entities.OnLightEntityState
			} else {
				light.Attributes[string(entities.StateLightEntityAttribute)] = entities.OffLightEntityState
			}
		}

		if light.HasAttribute(entities.BrightnessLightEntityAttribute) {
			light.Attributes[string(entities.BrightnessLightEntityAttribute)] = device.Light.State.Bri

		}

		// Add commands
		light.AddCommand(entities.OnLightEntityCommand, func(entity entities.LightEntity, params map[string]interface{}) int {

			if err := device.TurnOn(); err != nil {
				return 404
			}
			return 200
		})

		light.AddCommand(entities.OffLightEntityCommand, func(entity entities.LightEntity, params map[string]interface{}) int {

			if err := device.TurnOff(); err != nil {
				return 404
			}
			return 200
		})

		light.AddCommand(entities.ToggleLightEntityCommand, func(entity entities.LightEntity, params map[string]interface{}) int {
			if device.IsOn() {
				device.TurnOff()
			} else {
				device.TurnOff()
			}
			return 200
		})

		// Set the Handle State Change function
		device.SetHandleChangeStateFunc(func(state *deconz.DeconzState) {
			log.WithFields(log.Fields{
				"ID":   device.GetID(),
				"Type": device.Type,
			}).Debug("Handle State Change")

			attributes := make(map[string]interface{})

			if light.HasAttribute(entities.StateLightEntityAttribute) {
				if *state.On {
					attributes[string(entities.StateLightEntityAttribute)] = entities.OnLightEntityState
				} else {
					attributes[string(entities.StateLightEntityAttribute)] = entities.OffLightEntityState
				}
			}

			if light.HasAttribute(entities.BrightnessLightEntityAttribute) {
				attributes[string(entities.BrightnessLightEntityAttribute)] = state.Bri

			}

			light.SetAttributes(attributes)

		})

		c.IntegrationDriver.AddEntity(light)

	case deconz.GroupDeconzDeviceType:

		group := entities.NewLightEntity(fmt.Sprintf("group%d", device.GetID()), entities.LanguageText{En: device.GetName()}, "")
		// Group only allows for on/off -> basic switch, no dimming
		group.AddFeature(entities.OnOffLightEntityFeatures)
		group.AddFeature(entities.ToggleLightEntityFeatures)

		// Set initial attribute
		// if group.HasAttribute(entities.StateLightEntityAttribute) {
		// 	if *&device.Group.State.AnyOn {
		// 		light.Attributes[string(entities.StateLightEntityAttribute)] = entities.OnLightEntityState
		// 	} else {
		// 		light.Attributes[string(entities.StateLightEntityAttribute)] = entities.OffLightEntityState
		// 	}
		// }

		// Commands

		group.AddCommand(entities.OnLightEntityCommand, func(entity entities.LightEntity, params map[string]interface{}) int {

			if err := device.TurnOn(); err != nil {
				return 404
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
				device.TurnOff()
			} else {
				device.TurnOff()
			}
			return 200
		})

		device.SetHandleChangeStateFunc(func(state *deconz.DeconzState) {
			log.WithFields(log.Fields{
				"ID":   device.GetID(),
				"Type": device.Type,
			}).Debug("Handle State Change")

			attributes := make(map[string]interface{})

			if group.HasAttribute(entities.StateLightEntityAttribute) {
				if *&state.AnyOn {
					attributes[string(entities.StateLightEntityAttribute)] = entities.OnLightEntityState
				} else {
					attributes[string(entities.StateLightEntityAttribute)] = entities.OffLightEntityState
				}
			}

			group.SetAttributes(attributes)

		})

		c.IntegrationDriver.AddEntity(group)
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

	case deconz.LightDeconzDeviceType:
		light := entities.NewLightEntity(fmt.Sprintf("light%d", device.GetID()), entities.LanguageText{En: device.GetName()}, "")
		c.IntegrationDriver.RemoveEntity(light)

	case deconz.GroupDeconzDeviceType:
		group := entities.NewLightEntity(fmt.Sprintf("group%d", device.GetID()), entities.LanguageText{En: device.GetName()}, "")
		c.IntegrationDriver.RemoveEntity(group)
	}

}

// Callen on RT connect
func (c *DeconzClient) deconzClientLoop() {

	defer func() {
		c.deconz.Stop()
		c.setDeviceState(integration.DisconnectedDeviceState)
	}()

	if c.deconz == nil {
		c.setupDeconz()
	} else {
		return
	}

	if c.deconz != nil {
		c.configureDeconz()

		go c.deconz.StartandListenLoop()

	} else {
		return
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
