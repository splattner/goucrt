package denonavr

func (d *DenonAVR) TurnOn() error {
	d.sendCommandToDevice(DenonCommandPower, "ON")
	_, err := d.sendCommandToDevice(DennonCommandZoneMain, "ON")

	return err
}

func (d *DenonAVR) TurnOff() error {
	d.sendCommandToDevice(DenonCommandPower, "STANDBY")
	_, err := d.sendCommandToDevice(DennonCommandZoneMain, "OFF")

	return err
}

func (d *DenonAVR) TogglePower() error {

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
