package tasmota

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (e *Tasmota) publishMqttCommand(topic string, value interface{}) error {
	log.WithFields(log.Fields{
		"topic": topic,
		"value": fmt.Sprintf("%v", value)}).Debug("Publish MQTT Command")

	if token := e.mqttClient.Publish(topic, 0, false, fmt.Sprintf("%v", value)); token.Wait() && token.Error() != nil {
		log.WithError(token.Error()).Error("MQTT publish failed")
		return token.Error()
	}
	return nil
}

func (e *Tasmota) subscribeMqttTopic(topic string, callback mqtt.MessageHandler) {

	log.WithField("topic", topic).Debug("MQTT Subscribe to topic")
	if token := e.mqttClient.Subscribe(topic, 0, callback); token.Wait() && token.Error() != nil {
		log.WithError(token.Error()).Error("MQTT subscribe failed")
	}
}

func (e *Tasmota) unsubscribeMqttTopic(topic string) {

	log.WithField("topic", topic).Debug("MQTT Unsubscribe to topic")
	if token := e.mqttClient.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		log.WithError(token.Error()).Error("MQTT unsubscribe failed")
	}
}
