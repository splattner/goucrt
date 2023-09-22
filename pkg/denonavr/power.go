package denonavr

func (d *DenonAVR) TurnOn() int {
	d.sendCommandToDevice(DenonCommandPower, "ON")
	status, _ := d.sendCommandToDevice(DennonCommandZoneMain, "ON")

	return status
}

func (d *DenonAVR) TurnOff() int {
	d.sendCommandToDevice(DenonCommandPower, "STANDBY")
	status, _ := d.sendCommandToDevice(DennonCommandZoneMain, "OFF")

	return status
}

func (d *DenonAVR) TogglePower() int {

	if d.IsOn() {
		return d.TurnOn()
	}

	return d.TurnOff()
}

func (d *DenonAVR) IsOn() bool {

	switch d.mainZoneData.Power {
	case "ON":
		return true
	default:
		return false
	}
}
