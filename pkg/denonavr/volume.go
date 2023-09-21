package denonavr

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func (d *DenonAVR) GetMasterVolume() string {

	return d.data.MasterVolume
}

func (d *DenonAVR) GetVolume() float64 {

	var volume float64
	if s, err := strconv.ParseFloat(d.GetMasterVolume(), 64); err == nil {
		volume = s
	}

	return volume + 80

}

func (d *DenonAVR) SetVolume(volume float64) {

	// The Volume command need the following
	// 10.5 -> MV105
	// 11 -> MV11

	var convertedVolume string

	if volume != math.Trunc(volume) {
		convertedVolume = strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", volume*10), "0"), ".")
	} else {
		convertedVolume = strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", volume), "0"), ".")
	}

	d.sendCommandToDevice(DenonCommandVolume, convertedVolume)

}

func (d *DenonAVR) Mute() {

	d.sendCommandToDevice(DenonCommandMute, "ON")
}

func (d *DenonAVR) UnMute() {

	d.sendCommandToDevice(DenonCommandMute, "OFF")
}

func (d *DenonAVR) MuteToggle() {

	if d.Muted() {
		d.UnMute()

	} else {
		d.Mute()
	}

}

func (d *DenonAVR) Muted() bool {

	switch d.data.Mute {
	case "on":
		return true
	default:
		return false
	}
}

func (d *DenonAVR) SetVolumeUp() {

	newVolume := d.GetVolume() + DenonVolumeStep
	d.SetVolume(newVolume)

}

func (d *DenonAVR) SetVolumeDown() {

	newVolume := d.GetVolume() - DenonVolumeStep
	d.SetVolume(newVolume)

}
