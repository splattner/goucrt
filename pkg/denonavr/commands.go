package denonavr

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (d *DenonAVR) sendCommandToDevice(denonCommandType DenonCommand, command string) (int, error) {

	url := "http://" + d.Host + COMMAND_URL + "?" + string(denonCommandType) + command
	log.WithFields(log.Fields{
		"type":    string(denonCommandType),
		"command": command,
		"url":     url}).Info("Send Command to Denon Device")

	req, err := http.Get(url)
	if err != nil {
		return req.StatusCode, fmt.Errorf("Error sending command: %w", err)
	}

	// Trigger a updata data, handeld in the Listen Looü
	d.updateTrigger <- "update"

	return req.StatusCode, nil
}

func (d *DenonAVR) SetMoni1Out() error {
	_, err := d.sendCommandToDevice(DenonCommandVS, "MONI1")

	return err
}

func (d *DenonAVR) SetMoni2Out() error {
	_, err := d.sendCommandToDevice(DenonCommandVS, "MONI2")

	return err
}

func (d *DenonAVR) SetMoniAutoOut() error {
	_, err := d.sendCommandToDevice(DenonCommandVS, "MONIAUTO")

	return err
}
