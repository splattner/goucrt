package denonavr

func (d *DenonAVR) SetMoni1Out() int {
	status, _ := d.sendCommandToDevice(DenonCommandCursorControl, "MONI1")

	return status
}

func (d *DenonAVR) SetMoni2Out() int {
	status, _ := d.sendCommandToDevice(DenonCommandCursorControl, "MONI1")

	return status
}

func (d *DenonAVR) SetMoniAutoOut() int {
	status, _ := d.sendCommandToDevice(DenonCommandCursorControl, "MONIAUTO")

	return status
}
