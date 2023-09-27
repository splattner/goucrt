package shelly

import (
	"fmt"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type ShellyDevice struct {
	shelly               *Shelly
	Id                   string `json:"id,omitempty"`
	Model                string `json:"model,omitempty"`
	MACAddress           string `json:"mac,omitempty"`
	IPAddress            string `json:"ip,omitempty"`
	NewFirewareAvailable bool   `json:"new_fw,omitempty"`
	FirmewareVersion     string `json:"fw_ver,omitempty"`

	State string

	handleMsgReceivedFunc map[string][]func([]byte)
}

func (d *ShellyDevice) newShellyDevice(shelly *Shelly) {

	d.shelly = shelly
	d.configureCallbacks()

	d.handleMsgReceivedFunc = make(map[string][]func([]byte))
}

// Add a function that is called when a message is eceiverd from a Shelly device on a selected topic
func (d *ShellyDevice) AddMsgReceivedFunc(topic string, f func(payload []byte)) {
	d.handleMsgReceivedFunc[topic] = append(d.handleMsgReceivedFunc[topic], f)
}

// Call all MsgReceivedFunc for this device and topic
func (d *ShellyDevice) stateChangeHandler(topic string, payload []byte) {

	if d.handleMsgReceivedFunc[topic] != nil {
		for _, f := range d.handleMsgReceivedFunc[topic] {
			f(payload)
		}
	}

}

func (e *ShellyDevice) configureCallbacks() {
	log.WithField("ID", e.Id).Debug("Subscribe to Shelly Topic for this device")
	// Add callback for this device
	topic := fmt.Sprintf("shellies/%s/#", e.Id)
	e.shelly.subscribeMqttTopic(topic, e.mqttCallback())
}

func (e *ShellyDevice) mqttCallback() mqtt.MessageHandler {
	var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {

		log.WithFields(log.Fields{
			"ID":    e.Id,
			"Topic": msg.Topic(),
			"Msg:":  string(msg.Payload()),
		}).Trace("Received Message from Shelly")

		topic := strings.TrimLeft(msg.Topic(), "shellies/"+e.Id+"/")

		switch topic {
		case "relay/0":
			// Only call state chage handler for relay/0 when something has schanged
			if e.State != string(msg.Payload()) {
				// Call the state change handler function
				e.stateChangeHandler(topic, msg.Payload())
				// Set internal state
				e.State = string(msg.Payload())
			}
		default:
			// Call the state change handler function
			e.stateChangeHandler(topic, msg.Payload())
		}

	}

	return f
}

func (e *ShellyDevice) TurnOn() error {
	return e.shelly.publishMqttCommand("shellies/"+e.Id+"/relay/0/command", "on")
}

func (e *ShellyDevice) TurnOff() error {
	return e.shelly.publishMqttCommand("shellies/"+e.Id+"/relay/0/command", "off")
}

func (e *ShellyDevice) IsOn() bool {
	return e.State == "on"
}

func (e *ShellyDevice) Toggle() error {

	if e.IsOn() {
		return e.TurnOff()
	}

	return e.TurnOn()
}
