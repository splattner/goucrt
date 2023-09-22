package denonavr

func (d *DenonAVR) Play() int {
	status, _ := d.sendCommandToDevice(DenonCommandNS, "9A")

	return status
}

func (d *DenonAVR) Pause() int {
	status, _ := d.sendCommandToDevice(DenonCommandNS, "9B")

	return status
}
