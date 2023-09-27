package shelly

import (
	"encoding/json"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type Shelly struct {
	mqttClient mqtt.Client

	handleDeviceDiscoveredFunc func(*ShellyDevice)
}

func NewShelly(mqttClient mqtt.Client) *Shelly {

	shelly := Shelly{}
	shelly.mqttClient = mqttClient

	return &shelly

}

// Set the function that get called when a new Deconz Device is discovered
func (s *Shelly) SetDeviceDiscoveredHandler(f func(*ShellyDevice)) {
	s.handleDeviceDiscoveredFunc = f
}

func (s *Shelly) Start() error {
	// Connect to MQTT Broker
	if token := s.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.WithError(token.Error()).Error("MQTT connect failed")

		return token.Error()
	}

	return nil
}

func (s *Shelly) Stop() {
	if s.mqttClient.IsConnected() {
		s.mqttClient.Disconnect(0)
	}
}

func (s *Shelly) StartDiscovery() {
	log.Info(("Starting Shelly Device discovery"))

	s.subscribeMqttTopic("shellies/announce", s.mqttDiscoverCallback())
	s.subscribeMqttTopic("shellies/+/info", s.mqttDiscoverCallback())
	s.publishMqttCommand("shellies/command", "announce")
}

func (s *Shelly) StopDiscovery() {
	log.Info(("Stop Shelly Device discovery"))

	s.unsubscribeMqttTopic("shellies/announce")
	s.unsubscribeMqttTopic("shellies/+/info")

}

func (s *Shelly) mqttDiscoverCallback() mqtt.MessageHandler {

	var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {

		log.WithFields(log.Fields{
			"Topic": string(msg.Topic()),
			"Msg":   string(msg.Payload()),
		}).Trace("MQTT Mesage for Shelly Device discovery")

		if strings.Contains(msg.Topic(), "announce") {
			log.WithFields(log.Fields{
				"Topic": string(msg.Topic()),
				"Msg":   string(msg.Payload()),
			}).Trace("Announce MQTT Mesage for Shelly Device discovery")

			shellyDevice := ShellyDevice{}

			err := json.Unmarshal(msg.Payload(), &shellyDevice)
			if err != nil {
				log.WithError(err).Fatal("Unmarshal to Shelly Device failed")
				return
			}

			shellyDevice.newShellyDevice(s)

			// Handle device discovered
			if s.handleDeviceDiscoveredFunc != nil {
				s.handleDeviceDiscoveredFunc(&shellyDevice)
			}

		}
		// if strings.Contains(msg.Topic(), "shellies") && strings.Contains(msg.Topic(), "info") {
		// 	log.Println("Shelly info found", string(msg.Payload()))
		// }
	}

	return f
}
