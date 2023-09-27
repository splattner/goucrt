package client

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
	Client
	tasmota *tasmota.Tasmota
}

func NewTasmotaClient(i *integration.Integration) *TasmotaClient {
	tasmota := TasmotaClient{}

	tasmota.IntegrationDriver = i
	// Start without a connection
	tasmota.DeviceState = integration.DisconnectedDeviceState

	tasmota.messages = make(chan string)

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
		Version: "0.0.1",
		SetupDataSchema: integration.SetupDataSchema{
			Title: integration.LanguageText{
				En: "Configuration",
				De: "Konfiguration",
			},
			Settings: []integration.SetupDataSchemaSettings{ipaddr, port, username, password},
		},
		Icon: "",
	}

	tasmota.IntegrationDriver.SetMetadata(&metadata)

	// set the client specific functions
	tasmota.initFunc = tasmota.initTasmotaClient
	tasmota.setupFunc = tasmota.tasmotaHandleSetup
	tasmota.clientLoopFunc = tasmota.tasmotaClientLoop
	//client.setDriverUserDataFunc = client.handleSetDriverUserData

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
		c.setDeviceState(integration.ErrorDeviceState)
	}

	// Handle connection to device this integration shall control
	// Set Device state to connected when connection is established
	c.setDeviceState(integration.ConnectedDeviceState)

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
		switchEntity.AddFeature(entities.OnOffSwitchEntityyFeatures)
		switchEntity.AddFeature(entities.ToggleSwitchEntityCommand)

		switchEntity.MapCommand(entities.OnSwitchEntityCommand, device.TurnOn)
		switchEntity.MapCommand(entities.OffSwitchEntityCommand, device.TurnOff)
		switchEntity.MapCommand(entities.ToggleLightEntityCommand, device.Toggle)

		device.AddMsgReceivedFunc("RESULT", func(msg []byte) {

			attributes := make(map[string]interface{})

			switch string(msg) {
			case "on":
				attributes[string(entities.StateSwitchEntityyAttribute)] = entities.OnSwitchtEntityState
			case "off":
				attributes[string(entities.StateSwitchEntityyAttribute)] = entities.OffSwitchtEntityState
			}

			switchEntity.SetAttributes(attributes)
		})

		tasmotaDevice = switchEntity

	case 4:
		// RGBW
		lightEntity := entities.NewLightEntity(device.Topic, entities.LanguageText{En: "Tasmota " + device.FriendlyName[0]}, "")

		lightEntity.AddFeature(entities.OnOffLightEntityFeatures)
		lightEntity.AddFeature(entities.ToggleLightEntityFeatures)
		lightEntity.AddFeature(entities.DimLightEntityFeatures)
		lightEntity.AddFeature(entities.ColorLightEntityFeatures)

		lightEntity.MapCommand(entities.OnLightEntityCommand, device.TurnOn)
		lightEntity.MapCommand(entities.OffLightEntityCommand, device.TurnOff)
		lightEntity.MapCommand(entities.ToggleLightEntityCommand, device.Toggle)

		device.AddMsgReceivedFunc("RESULT", func(msg []byte) {

			attributes := make(map[string]interface{})

			switch string(msg) {
			case "on":
				attributes[string(entities.StateLightEntityAttribute)] = entities.OnLightEntityState
			case "off":
				attributes[string(entities.StateLightEntityAttribute)] = entities.OffLightEntityState
			}

			lightEntity.SetAttributes(attributes)
		})

		tasmotaDevice = lightEntity

	}

	if tasmotaDevice != nil {
		c.IntegrationDriver.AddEntity(tasmotaDevice)
	}

}

func (c *TasmotaClient) handleRemoveDevice(device *tasmota.TasmotaDevice) {
	log.WithFields(log.Fields{
		"Topic":       device.Topic,
		"IP Address":  device.IPAddress,
		"MAC Address": device.MACAddress,
	}).Debug("Tasmota Device not available anymore")
}

// Callen on RT connect
func (c *TasmotaClient) tasmotaClientLoop() {

	defer func() {
		c.tasmota.StopDiscovery()
		c.tasmota.Stop()
		c.setDeviceState(integration.DisconnectedDeviceState)
	}()

	if c.tasmota == nil {
		c.setupTasmota()
	} else {
		return
	}

	if c.tasmota != nil {
		c.startTasmota()
	} else {
		return
	}

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
