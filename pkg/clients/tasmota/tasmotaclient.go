package tasmotaclient

import (
	"fmt"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"github.com/splattner/goucrt/pkg/entities"
	"github.com/splattner/goucrt/pkg/integration"
	"github.com/splattner/goucrt/pkg/tasmota"
)

// Tasmota Implementation
type TasmotaClient struct {
	integration.Client
	tasmota *tasmota.Tasmota

	mapOnState map[string]entities.LightEntityState
}

func NewTasmotaClient(i *integration.Integration) *TasmotaClient {
	tasmota := TasmotaClient{}

	tasmota.IntegrationDriver = i
	// Start without a connection
	tasmota.DeviceState = integration.DisconnectedDeviceState

	tasmota.Messages = make(chan string)

	ipaddr := integration.SetupDataSchemaSettings{
		Id: "mqtt_ipaddr",
		Label: integration.LanguageText{
			En: "MQTT Broker Address",
		},
		Field: integration.SettingTypeText{
			Text: integration.SettingTypeTextDefinition{
				Value: "",
			},
		},
	}

	port := integration.SetupDataSchemaSettings{
		Id: "mqtt_port",
		Label: integration.LanguageText{
			En: "MQTT Broker Port",
		},
		Field: integration.SettingTypeText{
			Text: integration.SettingTypeTextDefinition{
				Value: "1883",
			},
		},
	}

	username := integration.SetupDataSchemaSettings{
		Id: "mqtt_username",
		Label: integration.LanguageText{
			En: "MQTT Broker Username",
		},
		Field: integration.SettingTypeText{
			Text: integration.SettingTypeTextDefinition{
				Value: "",
			},
		},
	}

	password := integration.SetupDataSchemaSettings{
		Id: "mqtt_password",
		Label: integration.LanguageText{
			En: "MQTT Broker Password",
		},
		Field: integration.SettingTypeText{
			Text: integration.SettingTypeTextDefinition{
				Value: "",
			},
		},
	}

	metadata := integration.DriverMetadata{
		DriverId: "tasmota",
		Developer: integration.Developer{
			Name: "Sebastian Plattner",
		},
		Name: integration.LanguageText{
			En: "Tasmota",
		},
		Version: "0.2.0",
		SetupDataSchema: integration.SetupDataSchema{
			Title: integration.LanguageText{
				En: "Configuration",
				De: "Konfiguration",
			},
			Settings: []integration.SetupDataSchemaSettings{ipaddr, port, username, password},
		},
		Icon: "custom:tasmota.png",
	}

	tasmota.IntegrationDriver.SetMetadata(&metadata)

	// set the client specific functions
	tasmota.InitFunc = tasmota.initTasmotaClient
	tasmota.SetupFunc = tasmota.tasmotaHandleSetup
	tasmota.ClientLoopFunc = tasmota.tasmotaClientLoop
	//client.setDriverUserDataFunc = client.handleSetDriverUserData

	tasmota.mapOnState = map[string]entities.LightEntityState{
		"ON":  entities.OnLightEntityState,
		"OFF": entities.OffLightEntityState,
	}

	return &tasmota
}

func (c *TasmotaClient) initTasmotaClient() {

}

func (c *TasmotaClient) tasmotaHandleSetup(setup_data integration.SetupData) {
	// Finish the setup
	// Nothing to configure
	// Setup Data already persistet by integration Driver
	c.IntegrationDriver.SetDriverSetupState(integration.StopEvent, integration.OkState, "", nil)
}

func (c *TasmotaClient) setupTasmota() {

	if c.tasmota == nil {

		if c.IntegrationDriver.SetupData["mqtt_ipaddr"] != "" {

			ipaddr := c.IntegrationDriver.SetupData["mqtt_ipaddr"]
			port, _ := strconv.Atoi(c.IntegrationDriver.SetupData["mqtt_port"])
			mqttBroker := fmt.Sprintf("tcp://%s:%d", ipaddr, port)

			log.WithFields(log.Fields{
				"MQTT Host":   ipaddr,
				"MQTT Port":   port,
				"MQTT Broker": mqttBroker,
				"ClientID":    c.IntegrationDriver.Metadata.DriverId}).Info("Connecting to MQTT Host")

			opts := mqtt.NewClientOptions().AddBroker(mqttBroker).SetClientID(c.IntegrationDriver.Metadata.DriverId)

			opts.SetKeepAlive(60 * time.Second)
			opts.SetPingTimeout(1 * time.Second)
			opts.SetProtocolVersion(3)
			opts.SetOrderMatters(false)
			if c.IntegrationDriver.SetupData["mqtt_username"] != "" && c.IntegrationDriver.SetupData["mqtt_password"] != "" {
				opts.SetUsername(c.IntegrationDriver.SetupData["mqtt_username"])
				opts.SetPassword(c.IntegrationDriver.SetupData["mqtt_password"])
			}

			mqttClient := mqtt.NewClient(opts)
			c.tasmota = tasmota.NewTasmota(mqttClient)
			c.tasmota.SetDeviceDiscoveredHandler(c.handleNewDeviceDiscovered)

		} else {
			log.Error("Cannot setup Tasmota Client, missing setupData")
		}
	}

}

func (c *TasmotaClient) startTasmota() {

	log.Debug("Start and connect Tamota")

	if err := c.tasmota.Start(); err != nil {
		c.SetDeviceState(integration.ErrorDeviceState)
	}

	// Handle connection to device this integration shall control
	// Set Device state to connected when connection is established
	c.SetDeviceState(integration.ConnectedDeviceState)

	c.tasmota.StartDiscovery()

}

func (c *TasmotaClient) handleNewDeviceDiscovered(device *tasmota.TasmotaDevice) {
	log.WithFields(log.Fields{
		"Topic":       device.Topic,
		"IP Address":  device.IPAddress,
		"MAC Address": device.MACAddress,
	}).Debug("New Tasmota Device discovered")

	var tasmotaDevice interface{}

	switch device.LightSubtype {
	case 0:
		// Sonoff Basic
		switchEntity := entities.NewSwitchEntity(device.Topic, entities.LanguageText{En: "Tasmota " + device.FriendlyName[0]}, "")

		switchEntity.SubscribeCallbackFunc = device.Subscribe
		switchEntity.UnsubscribeCallbackFunc = device.Unsubscribe

		switchEntity.AddFeature(entities.OnOffSwitchEntityyFeatures)
		switchEntity.AddFeature(entities.ToggleSwitchEntityyFeatures)

		switchEntity.MapCommand(entities.OnSwitchEntityCommand, device.TurnOn)
		switchEntity.MapCommand(entities.OffSwitchEntityCommand, device.TurnOff)
		switchEntity.MapCommand(entities.ToggleSwitchEntityCommand, device.Toggle)

		device.AddMsgReceivedFunc("RESULT", func(msg interface{}) {

			res := msg.(tasmota.TasmotaResultMsg)

			attributes := make(map[string]interface{})

			if res.Power == "ON" || res.Power1 == "ON" {
				attributes[string(entities.StateLightEntityAttribute)] = entities.OnLightEntityState
			} else {
				attributes[string(entities.StateLightEntityAttribute)] = entities.OffLightEntityState
			}

			switchEntity.SetAttributes(attributes)
		})

		tasmotaDevice = switchEntity

	case 4:
		// RGBW
		lightEntity_rgb := entities.NewLightEntity(device.Topic, entities.LanguageText{En: "Tasmota " + device.FriendlyName[0]}, "")

		lightEntity_rgb.SubscribeCallbackFunc = device.Subscribe
		lightEntity_rgb.UnsubscribeCallbackFunc = device.Unsubscribe

		lightEntity_rgb.AddFeature(entities.OnOffLightEntityFeatures)
		lightEntity_rgb.AddFeature(entities.ToggleLightEntityFeatures)
		lightEntity_rgb.AddFeature(entities.DimLightEntityFeatures)
		lightEntity_rgb.AddFeature(entities.ColorLightEntityFeatures)

		// Add commands
		lightEntity_rgb.AddCommand(entities.OnLightEntityCommand, func(entity entities.LightEntity, params map[string]interface{}) int {

			// NO param set, so just turn on
			if len(params) == 0 {
				if err := device.TurnOn(); err != nil {
					return 404
				}
			} else {
				if params["saturation"] != nil && params["hue"] != nil {

					hue := float32(params["hue"].(float64))
					sat := float32(params["saturation"].(float64) / 255 * 100)

					// Color Light
					if err := device.SetHue(hue); err != nil {
						return 404
					}
					if err := device.SetSaturation(sat); err != nil {
						return 404
					}

				}

				if params["brightness"] != nil {
					bri := int(params["brightness"].(float64) / 255 * 100)
					if bri > 0 && device.LocalState.White == 0 {
						// Set Brightness if not in White mode
						if err := device.SetBrightness(bri); err != nil {
							return 404
						}
					} else {

						// When in color mode, and bri is 0, set/turn on white mode
						// Setting white to 0 turns off the white mode and return to color mode
						if err := device.SetWhite(bri); err != nil {
							return 404
						}
					}

				}

			}

			return 200
		})

		lightEntity_rgb.MapCommand(entities.OffLightEntityCommand, device.TurnOff)
		lightEntity_rgb.MapCommand(entities.ToggleLightEntityCommand, device.Toggle)

		device.AddMsgReceivedFunc("RESULT", func(msg interface{}) {

			res := msg.(tasmota.TasmotaResultMsg)

			log.WithFields(log.Fields{"res": res, "Device": device.FriendlyName}).Debug("Result msg received")

			attributes := make(map[string]interface{})

			if res.Power == "ON" || res.Power1 == "ON" {
				attributes[string(entities.StateLightEntityAttribute)] = entities.OnLightEntityState
			} else {
				attributes[string(entities.StateLightEntityAttribute)] = entities.OffLightEntityState
			}

			// Only White light
			if res.White > 0 {
				attributes[string(entities.SaturationLightEntityAttribute)] = 0
				attributes[string(entities.BrightnessLightEntityAttribute)] = int(float32(res.White) / 100 * 255)
			} else {
				if res.HSBCOlor != "" {
					// Handle COlor Part of light
					hue, sat, bri := device.GetHSB(res.HSBCOlor)

					attributes[string(entities.HueLightEntityAttribute)] = int(hue)
					attributes[string(entities.SaturationLightEntityAttribute)] = int(float64(sat) / 100 * 255)
					attributes[string(entities.BrightnessLightEntityAttribute)] = int(float64(bri) / 100 * 255)

				}
			}

			lightEntity_rgb.SetAttributes(attributes)
		})

		tasmotaDevice = lightEntity_rgb

	}

	if tasmotaDevice != nil {
		if err := c.IntegrationDriver.AddEntity(tasmotaDevice); err != nil {
			log.WithError(err).Error("Cannot add Entity")
		}
	}

}

// func (c *TasmotaClient) handleRemoveDevice(device *tasmota.TasmotaDevice) {
// 	log.WithFields(log.Fields{
// 		"Topic":       device.Topic,
// 		"IP Address":  device.IPAddress,
// 		"MAC Address": device.MACAddress,
// 	}).Debug("Tasmota Device not available anymore")
// }

// Callen on RT connect
func (c *TasmotaClient) tasmotaClientLoop() {

	ticker := time.NewTicker(5 * time.Minute)

	defer func() {
		if c.tasmota != nil {
			c.tasmota.StopDiscovery()
			c.tasmota.Stop()
		}
		ticker.Stop()
		c.SetDeviceState(integration.DisconnectedDeviceState)
	}()

	if c.tasmota == nil {
		c.setupTasmota()
	}

	if c.tasmota != nil {
		c.startTasmota()
	} else {
		return
	}

	// Run Client Loop to handle entity changes from device
	for {
		msg := <-c.Messages

		switch msg {
		case "disconnect":
			return
		}

	}

}
