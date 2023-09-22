package denonavr

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

var SOURCE_MAPPING = map[string]string{
	"TV AUDIO":       "TV",
	"iPod/USB":       "USB/IPOD",
	"Bluetooth":      "BT",
	"Blu-ray":        "BD",
	"CBL/SAT":        "SAT/CBL",
	"NETWORK":        "NET",
	"Media Player":   "MPLAY",
	"AUX1":           "AUX1",
	"Tuner":          "TUNER",
	"FM":             "TUNER",
	"SpotifyConnect": "Spotify Connect",
}

func (d *DenonAVR) GetZoneInputFuncList(zone DenonZone) map[string]string {

	var inputFuncList map[string]string

	switch zone {
	case MainZone:
		inputFuncList = d.getInputFuncList(d.mainZoneStatus)
	case Zone2:
		inputFuncList = d.getInputFuncList(d.zone2Status)
	case Zone3:
		inputFuncList = d.getInputFuncList(d.zone3Status)
	}

	return inputFuncList
}

func (d *DenonAVR) getInputFuncList(zoneStatus DenonStatus) map[string]string {

	inputFuncList := make(map[string]string)

	// Only add those not deleted
	// Use renamed value
	for i, input := range zoneStatus.InputFuncList {
		if zoneStatus.SourceDelete[i] == "USE" {
			inputFuncList[input] = strings.TrimRight(zoneStatus.RenameSource[i], " ")
		}
	}

	return inputFuncList

}

func (d *DenonAVR) SetSelectSourceMainZone(source string) int {

	inputFuncList := d.GetZoneInputFuncList(MainZone)
	log.WithFields(log.Fields{
		"source":        source,
		"inputFuncList": inputFuncList,
		"sourceMapping": SOURCE_MAPPING}).Debug("Select Source Main Zone")

	var selectedSource string
	for sourceOrigin, renamedSource := range inputFuncList {
		if renamedSource == source {
			selectedSource = sourceOrigin
			break
		}

	}
	if SOURCE_MAPPING[selectedSource] != "" {
		status, _ := d.sendCommandToDevice(DenonCommandSelectInput, SOURCE_MAPPING[selectedSource])
		return status
	}

	return 404
}
