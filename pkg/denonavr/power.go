package denonavr

func (d *DenonAVR) TurnOn() {
	d.sendCommandToDevice(DenonCommandPower, "ON")
}

func (d *DenonAVR) TurnOff() {
	d.sendCommandToDevice(DenonCommandPower, "STANDBY")
}

func (d *DenonAVR) TogglePower() {

	if d.IsOn() {
		d.TurnOn()
	} else {
		d.TurnOff()
	}
}

func (d *DenonAVR) IsOn() bool {

	switch d.data.Power {
	case "ON":
		return true
	default:
		return false
	}
}
