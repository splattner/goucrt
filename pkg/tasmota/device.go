package tasmota

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type TasmotaDevice struct {
	tasmota *Tasmota

	IPAddress       string         `json:"ip,omitempty"`
	DeviceName      string         `json:"dn,omitempty"`
	FriendlyName    []string       `json:"fn,omitempty"`
	Hostname        string         `json:"hn,omitempty"`
	MACAddress      string         `json:"mac,omitempty"`
	Module          string         `json:"md,omitempty"`
	TuyaMCUFlag     int            `json:"ty,omitempty"`
	IFAN            int            `json:"if,omitempty"`
	DOffline        string         `json:"ofln,omitempty"`
	DOnline         string         `json:"onln,omitempty"`
	State           []string       `json:"st,omitempty"`
	SoftwareVersion string         `json:"sw,omitempty"`
	Topic           string         `json:"t,omitempty"`
	Fulltopic       string         `json:"ft,omitempty"`
	TopicPrefix     []string       `json:"tp,omitempty"`
	Relays          []int          `json:"rl,omitempty"`
	Switches        []int          `json:"swc,omitempty"`
	SWN             []int          `json:"swn,omitempty"`
	Buttons         []int          `json:"btn,omitempty"`
	SetOptions      map[string]int `json:"so,omitempty"`
	LK              int            `json:"lk,omitempty"`    // LightColor (LC) and RGB LinKed https://github.com/arendst/Tasmota/blob/development/tasmota/xdrv_04_light.ino#L689
	LightSubtype    int            `json:"lt_st,omitempty"` // https://github.com/arendst/Tasmota/blob/development/tasmota/xdrv_04_light.ino
	ShutterOptions  []int          `json:"sho,omitempty"`
	Version         int            `json:"ver,omitempty"`

	LocalState      TasmotaResultMsg
	LastTeleMessame TasmotaTeleMsg

	handleMsgReceivedFunc map[string][]func(interface{})
}

func (d *TasmotaDevice) newTasmotaDevice(tasmota *Tasmota) {

	d.tasmota = tasmota
	//d.subscribe()

	d.handleMsgReceivedFunc = make(map[string][]func(interface{}))
}

// Add a function that is called when a message is eceiverd from a Shelly device on a selected topic
func (d *TasmotaDevice) AddMsgReceivedFunc(topic string, f func(message interface{})) {
	d.handleMsgReceivedFunc[topic] = append(d.handleMsgReceivedFunc[topic], f)
}

// Call all MsgReceivedFunc for this device and topic
func (d *TasmotaDevice) stateChangeHandler(topic string, message interface{}) {

	if d.handleMsgReceivedFunc[topic] != nil {
		for _, f := range d.handleMsgReceivedFunc[topic] {
			f(message)
		}
	}

}

func (e *TasmotaDevice) Subscribe() {
	log.WithField("Topic", e.Topic).Debug("Subscribe to Tasmota Topic for this device")

	// Add callback for stat
	topicStat := fmt.Sprintf("stat/%s/#", e.Topic)
	e.tasmota.subscribeMqttTopic(topicStat, e.mqttCallback())

	// Add callback for tele
	topicTele := fmt.Sprintf("tele/%s/#", e.Topic)
	e.tasmota.subscribeMqttTopic(topicTele, e.mqttCallback())

	// Make sure to have an initial State
	if err := e.GetResult(); err != nil {
		log.WithError(err).Debug("Cannot get a Result")
	}
}

func (e *TasmotaDevice) Unsubscribe() {
	log.WithField("Topic", e.Topic).Debug("Unsubscribe to Tasmota Topic for this device")

	// Add callback for stat
	topicStat := fmt.Sprintf("stat/%s/#", e.Topic)
	e.tasmota.unsubscribeMqttTopic(topicStat)

	// Add callback for tele
	topicTele := fmt.Sprintf("tele/%s/#", e.Topic)
	e.tasmota.unsubscribeMqttTopic(topicTele)
}

func (d *TasmotaDevice) mqttCallback() mqtt.MessageHandler {
	var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {

		log.WithFields(log.Fields{
			"DeviceName": d.DeviceName,
			"Topic":      msg.Topic(),
			"Msg:":       string(msg.Payload()),
		}).Trace("Received Message from Shelly")

		topic := strings.TrimLeft(msg.Topic(), "tele/"+d.Topic+"/")
		topic = strings.TrimLeft(topic, "stat/"+d.Topic+"/")

		switch topic {
		case "RESULT":
			var resultMessage TasmotaResultMsg
			err := json.Unmarshal(msg.Payload(), &resultMessage)
			if resultMessage.CustomSend == "Done" {
				// We don't care for this message
				return
			}
			if err != nil {
				log.WithError(err).Debug("Unmarshal to TasmotaPowerMsg failed")
				return

			}

			// Set internal state
			d.LocalState = resultMessage

			d.stateChangeHandler(topic, resultMessage)
		case "SENSOR":
			err := json.Unmarshal(msg.Payload(), &d.LastTeleMessame)
			if err != nil {
				log.WithError(err).Debug("Unmarshal to TasmotaTeleMsg failed")
				return
			}

			d.stateChangeHandler(topic, msg.Payload())

		}

	}

	return f
}

func (e *TasmotaDevice) TurnOn() error {

	if err := e.tasmota.publishMqttCommand("cmnd/"+e.Topic+"/POWER", "ON"); err != nil {
		return err
	}

	// TUrn on does not send the current state of the device, so force it
	return e.GetResult()
}

func (e *TasmotaDevice) TurnOff() error {
	return e.tasmota.publishMqttCommand("cmnd/"+e.Topic+"/POWER", "OFF")
}

func (e *TasmotaDevice) IsOn() bool {
	return e.LocalState.Power == "ON" || e.LocalState.Power1 == "ON"
}

func (e *TasmotaDevice) Toggle() error {

	if e.IsOn() {
		return e.TurnOff()
	}

	return e.TurnOn()
}

func (d *TasmotaDevice) GetResult() error {
	return d.tasmota.publishMqttCommand("cmnd/"+d.Topic+"/HsbColor", "")
}

func (d *TasmotaDevice) SetBrightness(brightness int) error {
	return d.tasmota.publishMqttCommand("cmnd/"+d.Topic+"/HsbColor3", brightness)
}

func (d *TasmotaDevice) SetHue(hue float32) error {
	return d.tasmota.publishMqttCommand("cmnd/"+d.Topic+"/HsbColor1", hue)
}

func (d *TasmotaDevice) SetSaturation(saturation float32) error {
	return d.tasmota.publishMqttCommand("cmnd/"+d.Topic+"/HsbColor2", saturation)
}

func (d *TasmotaDevice) SetHSB(hue float32, saturation float32, brightness int) error {
	return d.tasmota.publishMqttCommand("cmnd/"+d.Topic+"/HsbColor", fmt.Sprintf("%.0f,%.0f,%d", hue, saturation, brightness))
}

func (d *TasmotaDevice) SetWhite(white int) error {
	//e.publishMqttCommand("cmnd/"+e.Topic+"/Color1", "0,0,0")
	return d.tasmota.publishMqttCommand("cmnd/"+d.Topic+"/White", white)
}

func (e *TasmotaDevice) SetColorTemp(ct float32) error {
	log.Warningln("Setting Color Temp not implemented")
	return nil
}

func (d *TasmotaDevice) GetHSB(color string) (float64, float64, int) {

	if color == "" {
		color = d.LocalState.HSBCOlor
	}

	hsb := strings.Split(color, ",")

	if len(hsb) == 3 {

		hue, err := strconv.Atoi(hsb[0])
		if err != nil {
			log.WithError(err).Error("Unable to parse HSB")
		}
		sat, err := strconv.Atoi(hsb[1])
		if err != nil {
			log.WithError(err).Error("Unable to parse HSB")
		}
		bri, err := strconv.Atoi(hsb[2])
		if err != nil {
			log.WithError(err).Error("Unable to parse HSB")
		}

		return float64(hue), float64(sat), bri
	}

	return 0, 0, 0
}
