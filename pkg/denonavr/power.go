package denonavr

func (d *DenonAVR) TurnOn() error {
	if _, err := d.sendCommandToDevice(DenonCommandPower, "ON"); err != nil {
		return err
	}
	_, err := d.sendCommandToDevice(DennonCommandZoneMain, "ON")

	return err
}

func (d *DenonAVR) TurnOff() error {

	if _, err := d.sendCommandToDevice(DenonCommandPower, "STANDBY"); err != nil {
		return err
	}
	_, err := d.sendCommandToDevice(DennonCommandZoneMain, "OFF")

	return err
}

func (d *DenonAVR) TogglePower() error {

	if d.IsOn() {
		return d.TurnOff()
	}

	return d.TurnOn()
}

func (d *DenonAVR) IsOn() bool {

	switch d.mainZoneData.ZonePower {
	case "ON":
		return true
	default:
		return false
	}
}
