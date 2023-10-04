package client

import (
	"fmt"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"github.com/splattner/goucrt/pkg/entities"
	"github.com/splattner/goucrt/pkg/integration"
	"github.com/splattner/goucrt/pkg/shelly"
)

// Shelly Implementation
type ShellyClient struct {
	Client
	shelly *shelly.Shelly
}

func NewShellyClient(i *integration.Integration) *ShellyClient {
	client := ShellyClient{}

	client.IntegrationDriver = i
	// Start without a connection
	client.DeviceState = integration.DisconnectedDeviceState

	client.messages = make(chan string)

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
		DriverId: "shelly",
		Developer: integration.Developer{
			Name: "Sebastian Plattner",
		},
		Name: integration.LanguageText{
			En: "Shelly",
		},
		Version: "0.0.1",
		SetupDataSchema: integration.SetupDataSchema{
			Title: integration.LanguageText{
				En: "Configuration",
				De: "Konfiguration",
			},
			Settings: []integration.SetupDataSchemaSettings{ipaddr, port, username, password},
		},
		Icon: "custom:shelly.png",
	}

	client.IntegrationDriver.SetMetadata(&metadata)

	// set the client specific functions
	client.initFunc = client.initShellyClient
	client.setupFunc = client.shellyHandleSetup
	client.clientLoopFunc = client.shellyClientLoop
	//client.setDriverUserDataFunc = client.handleSetDriverUserData

	return &client
}

func (c *ShellyClient) initShellyClient() {

}

func (c *ShellyClient) shellyHandleSetup(setup_data integration.SetupData) {
	// Finish the setup
	// Nothing to configure
	// Setup Data already persistet by integration Driver
	c.IntegrationDriver.SetDriverSetupState(integration.StopEvent, integration.OkState, "", nil)
}

func (c *ShellyClient) setupShelly() {

	if c.shelly == nil {

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
			c.shelly = shelly.NewShelly(mqttClient)
			c.shelly.SetDeviceDiscoveredHandler(c.handleNewDeviceDiscovered)
		} else {
			log.Error("Cannot setup Shelly Client, missing setupData")
		}
	}

}

func (c *ShellyClient) startShelly() {

	log.Debug("Start and connect Shelly")

	if err := c.shelly.Start(); err != nil {
		c.setDeviceState(integration.ErrorDeviceState)
	}

	// Handle connection to device this integration shall control
	// Set Device state to connected when connection is established
	c.setDeviceState(integration.ConnectedDeviceState)

	c.shelly.StartDiscovery()

}

func (c *ShellyClient) handleNewDeviceDiscovered(device *shelly.ShellyDevice) {
	log.WithFields(log.Fields{
		"ID":          device.Id,
		"IP Address":  device.IPAddress,
		"MAC Address": device.MACAddress,
	}).Debug("New Shelly Device discovered")

	shellySwitch := entities.NewSwitchEntity(device.Id, entities.LanguageText{En: "Shelly " + device.Id}, "")

	shellySwitch.SubscribeCallbackFunc = device.Subscribe
	shellySwitch.UnsubscribeCallbackFunc = device.Unsubscribe

	shellySwitch.AddFeature(entities.OnOffSwitchEntityyFeatures)
	shellySwitch.AddFeature(entities.ToggleSwitchEntityyFeatures)

	shellySwitch.MapCommand(entities.OnSwitchEntityCommand, device.TurnOn)
	shellySwitch.MapCommand(entities.OffSwitchEntityCommand, device.TurnOff)
	shellySwitch.MapCommand(entities.ToggleSwitchEntityCommand, device.Toggle)

	device.AddMsgReceivedFunc("relay/0", func(msg []byte) {

		attributes := make(map[string]interface{})

		switch string(msg) {
		case "on":
			attributes[string(entities.StateSwitchEntityyAttribute)] = entities.OnSwitchtEntityState
		case "off":
			attributes[string(entities.StateSwitchEntityyAttribute)] = entities.OffSwitchtEntityState
		}

		shellySwitch.SetAttributes(attributes)
	})

	if err := c.IntegrationDriver.AddEntity(shellySwitch); err != nil {
		log.WithError(err).Error("Cannot add Entity")
	}

}

// func (c *ShellyClient) handleRemoveDevice(device *shelly.ShellyDevice) {
// 	log.WithFields(log.Fields{
// 		"ID":          device.Id,
// 		"IP Address":  device.IPAddress,
// 		"MAC Address": device.MACAddress,
// 	}).Debug("New Shelly Device not available anymore")
// }

// Callen on RT connect
func (c *ShellyClient) shellyClientLoop() {

	defer func() {
		c.shelly.StopDiscovery()
		c.shelly.Stop()
		c.setDeviceState(integration.DisconnectedDeviceState)
	}()

	if c.shelly == nil {
		c.setupShelly()
	} else {
		return
	}

	if c.shelly != nil {
		c.startShelly()
	} else {
		return
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
