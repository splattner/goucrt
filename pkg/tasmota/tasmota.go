package tasmota

import (
	"encoding/json"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type Tasmota struct {
	mqttClient mqtt.Client

	handleDeviceDiscoveredFunc func(*TasmotaDevice)
}

func NewTasmota(mqttClient mqtt.Client) *Tasmota {

	tasmota := Tasmota{}
	tasmota.mqttClient = mqttClient

	return &tasmota

}

// Set the function that get called when a new Deconz Device is discovered
func (s *Tasmota) SetDeviceDiscoveredHandler(f func(*TasmotaDevice)) {
	s.handleDeviceDiscoveredFunc = f
}

func (s *Tasmota) Start() error {
	// Connect to MQTT Broker
	if token := s.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.WithError(token.Error()).Error("MQTT connect failed")

		return token.Error()
	}

	return nil
}

func (s *Tasmota) Stop() {
	if s.mqttClient.IsConnected() {
		s.mqttClient.Disconnect(0)
	}
}

func (s *Tasmota) StartDiscovery() {

	log.Info(("Starting Shelly Device discovery"))

	s.subscribeMqttTopic("tasmota/discovery/#", s.mqttDiscoverCallback())
}

func (s *Tasmota) StopDiscovery() {
	log.Info(("Stop Shelly Device discovery"))

	s.unsubscribeMqttTopic("tasmota/discovery/#")

}

func (s *Tasmota) mqttDiscoverCallback() mqtt.MessageHandler {

	var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {

		log.WithFields(log.Fields{
			"Topic": string(msg.Topic()),
			"Msg":   string(msg.Payload()),
		}).Trace("MQTT Mesage for Tasmota Device discovery")

		if strings.Contains(msg.Topic(), "config") {
			log.WithFields(log.Fields{
				"Topic": string(msg.Topic()),
				"Msg":   string(msg.Payload()),
			}).Trace("MQTT Mesage for Tasmota Device discovery")

			tasmotaDevice := TasmotaDevice{}

			err := json.Unmarshal(msg.Payload(), &tasmotaDevice)
			if err != nil {
				log.WithError(err).Fatal("Unmarshal to Tasmota Device failed")
				return
			}

			tasmotaDevice.newTasmotaDevice(s)

			// Handle device discovered
			if s.handleDeviceDiscoveredFunc != nil {
				s.handleDeviceDiscoveredFunc(&tasmotaDevice)
			}

		}

	}

	return f
}
