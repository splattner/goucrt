package shelly

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (e *Shelly) publishMqttCommand(topic string, value interface{}) error {
	if token := e.mqttClient.Publish(topic, 0, false, fmt.Sprintf("%v", value)); token.Wait() && token.Error() != nil {
		log.WithError(token.Error()).Error("MQTT publish failed")
		return token.Error()
	}
	return nil
}

func (e *Shelly) subscribeMqttTopic(topic string, callback mqtt.MessageHandler) {

	log.WithField("topic", topic).Debug("MQTT Subscribe to topic")
	if token := e.mqttClient.Subscribe(topic, 0, callback); token.Wait() && token.Error() != nil {
		log.WithError(token.Error()).Error("MQTT subscribe failed")
	}
}
