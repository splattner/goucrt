package denonavr

func (d *DenonAVR) SetMoni1Out() error {
	_, err := d.sendCommandToDevice(DenonCommandCursorControl, "MONI1")

	return err
}

func (d *DenonAVR) SetMoni2Out() error {
	_, err := d.sendCommandToDevice(DenonCommandCursorControl, "MONI1")

	return err
}

func (d *DenonAVR) SetMoniAutoOut() error {
	_, err := d.sendCommandToDevice(DenonCommandCursorControl, "MONIAUTO")

	return err
}
