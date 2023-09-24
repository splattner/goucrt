package deconz

import (
	"fmt"
	"math"

	deconzgroup "github.com/jurgen-kluft/go-conbee/groups"
	deconzlight "github.com/jurgen-kluft/go-conbee/lights"
	deconzsensor "github.com/jurgen-kluft/go-conbee/sensors"
)

type DeconzDeviceType string

const (
	LightDeconzDeviceType  DeconzDeviceType = "light"
	GroupDeconzDeviceType                   = "group"
	SensorDeconzDeviceType                  = "sensor"
)

type DeconzDevice struct {
	deconz *Deconz

	Type DeconzDeviceType

	Light  deconzlight.Light
	Group  deconzgroup.Group
	Sensor deconzsensor.Sensor

	// For when there are multiple buttons on one sensor
	// sensorButtonId is the identifier to get the correct device
	sensorButtonId int

	handleStateChangeFunc func(state *DeconzState)
}

// Call Tyoe specific functions
// Add Device to DeCONZ Client
func (d *DeconzDevice) NewDeconzDevice(deconz *Deconz) {

	d.deconz = deconz

	switch d.Type {
	case LightDeconzDeviceType:
		d.newDeconzLightDevice()
	case GroupDeconzDeviceType:
		d.newDeconzGroupDevice()
	case SensorDeconzDeviceType:
		d.newDeconzSensorDevice()
	}

	d.deconz.addDevice(d)

}

// Set the function that is called when a Stage change event is receiverd from DeCONZ Websocket
func (d *DeconzDevice) SetHandleChangeStateFunc(f func(state *DeconzState)) {
	d.handleStateChangeFunc = f
}

// Call the State Change Handler for this device
func (d *DeconzDevice) stateChangeHandler(state *DeconzState) {
	if d.handleStateChangeFunc != nil {
		d.handleStateChangeFunc(state)
	}

}

// Return the ID of the Device based on its type
func (d *DeconzDevice) GetID() int {

	switch d.Type {
	case LightDeconzDeviceType:
		return d.Light.ID
	case GroupDeconzDeviceType:
		return d.Group.ID
	case SensorDeconzDeviceType:
		return d.Sensor.ID
	}

	return 0
}

// Return the Name of the Device based on its type
func (d *DeconzDevice) GetName() string {

	switch d.Type {
	case LightDeconzDeviceType:
		return d.Light.Name
	case GroupDeconzDeviceType:
		return d.Group.Name
	case SensorDeconzDeviceType:
		return d.Sensor.Name
	}

	return ""
}

// Apply update from deconz device
func (d *DeconzDevice) SetValue(value float32, channelName string) error {

	switch channelName {

	case "basic_switch", "brightness":
		brightness := float32(math.Round(float64(value)))
		return d.SetBrightness(brightness)
	case "hue":
		return d.SetHue(value)
	case "saturation":
		return d.SetSaturation(value)
	case "colortemp":
		return d.SetColorTemp(value)
	}

	return fmt.Errorf("Channel Name not found")
}

func (d *DeconzDevice) TurnOn() error {

	switch d.Type {
	case LightDeconzDeviceType:
		d.Light.State.SetOn(true)
	case GroupDeconzDeviceType:
		d.Group.Action.SetOn(true)
	}

	return d.setState()
}

func (d *DeconzDevice) TurnOff() error {

	switch d.Type {
	case LightDeconzDeviceType:
		d.Light.State.SetOn(false)
	case GroupDeconzDeviceType:
		d.Group.Action.SetOn(false)
	}

	return d.setState()
}

func (d *DeconzDevice) IsOn() bool {
	switch d.Type {
	case LightDeconzDeviceType:
		return *d.Light.State.On
	case GroupDeconzDeviceType:
		return *d.Group.Action.On
	}

	return false
}

func (d *DeconzDevice) SetBrightness(brightness float32) error {

	switch d.Type {
	case LightDeconzDeviceType:
		if brightness == 0 {
			d.Light.State.SetOn(false)
		} else {
			d.Light.State.SetOn(true)
		}

		bri_converted := uint8(math.Round(float64(brightness) / 100 * 255))
		d.Light.State.Bri = &bri_converted
	case GroupDeconzDeviceType:
		if brightness == 0 {
			d.Group.Action.SetOn(false)
		} else {
			d.Group.Action.SetOn(true)
		}

		bri_converted := uint8(math.Round(float64(brightness) / 100 * 255))
		d.Group.Action.Bri = &bri_converted
	}

	return d.setState()
}

func (d *DeconzDevice) SetColorTemp(ct float32) error {

	converted := uint16(ct)

	switch d.Type {
	case LightDeconzDeviceType:
		d.Light.State.CT = &converted
	case GroupDeconzDeviceType:
		d.Group.Action.CT = &converted
	}

	return d.setState()
}

func (d *DeconzDevice) SetHue(hue float32) error {

	converted := uint16(hue)

	switch d.Type {
	case LightDeconzDeviceType:
		d.Light.State.Hue = &converted
	case GroupDeconzDeviceType:
		d.Group.Action.Hue = &converted
	}

	return d.setState()
}

func (d *DeconzDevice) SetSaturation(saturation float32) error {

	converted := uint8(saturation)

	switch d.Type {
	case LightDeconzDeviceType:
		d.Light.State.Sat = &converted
	case GroupDeconzDeviceType:
		d.Group.Action.Sat = &converted
	}

	return d.setState()
}

func (d *DeconzDevice) setState() error {

	switch d.Type {
	case LightDeconzDeviceType:
		return d.setLightState()
	case GroupDeconzDeviceType:
		return d.setGroupState()
	}

	return fmt.Errorf("Device Type not found")

}
