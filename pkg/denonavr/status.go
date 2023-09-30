package denonavr

import (
	"encoding/xml"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type DenonStatus struct {
	XMLName         xml.Name `xml:"item"`
	Zone            string   `xml:"Zone>value"`
	Power           string   `xml:"Power>value"`
	InputFuncList   []string `xml:"InputFuncList>value"`
	RenameSource    []string `xml:"RenameSource>value>value"`
	SourceDelete    []string `xml:"SourceDelete>value"`
	InputFuncSelect string   `xml:"InputFuncSelect>value"`
	VolumeDisplay   string   `xml:"VolumeDisplay>value"`
	RestorerMode    string   `xml:"RestorerMode>value"`
	SurrMode        string   `xml:"SurrMode>value"`
	MasterVolume    string   `xml:"MasterVolume>value"`
	Mute            string   `xml:"Mute>value"`
	Model           string   `xml:"Model>value"`
}

func (d *DenonAVR) getZoneStatus(zone DenonZone) {
	switch zone {
	case MainZone:
		url := "http://" + d.Host + STATUS_URL
		d.mainZoneStatus = d.getZoneStatusFromDevice(url)

	case Zone2:

		url := "http://" + d.Host + STATUS_Z2_URL
		d.zone2Status = d.getZoneStatusFromDevice(url)

	case Zone3:
		url := "http://" + d.Host + STATUS_Z3_URL
		d.zone2Status = d.getZoneStatusFromDevice(url)
	}

}

// Return the Status from a Zone
func (d *DenonAVR) getZoneStatusFromDevice(url string) DenonStatus {
	status := DenonStatus{} // Somehow the values in the array are added instead of replaced. Not sure if this is the solution, but it works...
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("Cannot read response body")
	}

	if err := xml.Unmarshal(body, &status); err != nil {
		log.WithError(err).Info("Could not unmarshall")
	}

	return status
}
