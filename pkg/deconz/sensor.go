package deconz

import (
	"fmt"
)

type ButtonEvent int

// http://developer.digitalstrom.org/Architecture/ds-basics.pdf
const (
	Hold ButtonEvent = iota + 1
	ShortRelease
	LongRelease
	DoublePress
	TreeplePress
)

const (
	// milisecs
	SingleTip   int = 150
	SingleClick int = 50
)

func (d *DeconzDevice) newDeconzSensorDevice() {

}

func (e *DeconzDevice) getUniqueId() string {
	uniqueID := fmt.Sprintf("%s-%d", e.Sensor.UniqueID, e.sensorButtonId)
	return uniqueID
}

func (e *DeconzDevice) getName() string {
	name := fmt.Sprintf("%s Button %d", e.Sensor.Name, e.sensorButtonId+1)
	return name
}
