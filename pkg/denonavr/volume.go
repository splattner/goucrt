package denonavr

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func (d *DenonAVR) GetMainZoneVolume() string {

	return d.mainZoneStatus.MasterVolume
}

func (d *DenonAVR) GetVolume() float64 {

	var volume float64
	if s, err := strconv.ParseFloat(d.GetMainZoneVolume(), 64); err == nil {
		volume = s
	}

	return volume + 80

}

func (d *DenonAVR) SetVolume(volume float64) error {

	// The Volume command need the following
	// 10.5 -> MV105
	// 11 -> MV11

	var convertedVolume string

	if volume != math.Trunc(volume) {
		convertedVolume = strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", volume*10), "0"), ".")
	} else {
		convertedVolume = strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", volume), "0"), ".")
	}

	_, err := d.sendCommandToDevice(DenonCommandVolume, convertedVolume)

	return err

}

func (d *DenonAVR) MainZoneMute() error {

	_, err := d.sendCommandToDevice(DenonCommandMute, "ON")
	return err
}

func (d *DenonAVR) MainZoneUnMute() error {

	_, err := d.sendCommandToDevice(DenonCommandMute, "OFF")
	return err
}

func (d *DenonAVR) MainZoneMuteToggle() error {

	if d.MainZoneMuted() {
		return d.MainZoneUnMute()

	}
	return d.MainZoneMute()
}

func (d *DenonAVR) MainZoneMuted() bool {

	switch d.mainZoneStatus.Mute {
	case "on":
		return true
	default:
		return false
	}
}

func (d *DenonAVR) SetVolumeUp() error {

	newVolume := d.GetVolume() + DenonVolumeStep
	return d.SetVolume(newVolume)

}

func (d *DenonAVR) SetVolumeDown() error {

	newVolume := d.GetVolume() - DenonVolumeStep
	return d.SetVolume(newVolume)

}
