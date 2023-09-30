package denonavr

import (
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"

	"net/http"
)

type DenonCommand string
type DenonZone string

const (
	DenonCommandPower         DenonCommand = "PW"
	DennonCommandZoneMain     DenonCommand = "ZM"
	DenonCommandVolume        DenonCommand = "MV"
	DenonCommandMute          DenonCommand = "MU"
	DenonCommandSelectInput   DenonCommand = "SI"
	DenonCommandCursorControl DenonCommand = "MN"
	DenonCommandNS            DenonCommand = "NS"
	DenonCommandMS            DenonCommand = "MS"
	DenonCommandVS            DenonCommand = "VS"
	DenonVolumeStep           float64      = 1
)

const (
	MainZone DenonZone = "MAIN"
	Zone2    DenonZone = "Z2"
	Zone3    DenonZone = "Z3"
)

const (
	STATUS_URL    string = "/goform/formMainZone_MainZoneXmlStatus.xml"
	STATUS_Z2_URL string = "/goform/formZone2_Zone2XmlStatus.xml"
	STATUS_Z3_URL string = "/goform/formZone3_Zone3XmlStatus.xml"
	MAINZONE_URL  string = "/goform/formMainZone_MainZoneXml.xml"
	COMMAND_URL   string = "/goform/formiPhoneAppDirect.xml"
)

type DenonXML struct {
	XMLName          xml.Name     `xml:"item"`
	FriendlyName     string       `xml:"FriendlyName>value"`
	Power            string       `xml:"Power>value"`
	ZonePower        string       `xml:"ZonePower>value"`
	RenameZone       string       `xml:"RenameZone>value"`
	TopMenuLink      string       `xml:"TopMenuLink>value"`
	VideoSelectDisp  string       `xml:"VideoSelectDisp>value"`
	VideoSelect      string       `xml:"VideoSelect>value"`
	VideoSelectOnOff string       `xml:"VideoSelectOnOff>value"`
	VideoSelectList  []ValueLists `xml:"VideoSelectLists>value"`
	ECOModeDisp      string       `xml:"ECOModeDisp>value"`
	ECOMode          string       `xml:"ECOMode>value"`
	ECOModeList      []ValueLists `xml:"ECOModeLists>value"`
	AddSourceDisplay string       `xml:"AddSourceDisplay>value"`
	ModelId          string       `xml:"ModelId>value"`
	BrandId          string       `xml:"BrandId>value"`
	SalesArea        string       `xml:"SalesArea>value"`
	InputFuncSelect  string       `xml:"InputFuncSelect>value"`
	NetFuncSelect    string       `xml:"NetFuncSelect>value"`
	SelectSurround   string       `xml:"selectSurround>value"`
	VolumeDisplay    string       `xml:"VolumeDisplay>value"`
	MasterVolume     string       `xml:"MasterVolume>value"`
	Mute             string       `xml:"Mute>value"`
}

type ValueLists struct {
	Index string `xml:"index,attr"`
	Table string `xml:"table,attr"`
}

type DenonAVR struct {
	Host string

	mainZoneData DenonXML

	// Zone Status
	mainZoneStatus DenonStatus
	zone2Status    DenonStatus
	zone3Status    DenonStatus

	updateTrigger chan string

	entityChangedFunction map[string][]func(interface{})
}

func NewDenonAVR(host string) *DenonAVR {

	denonavr := DenonAVR{}

	denonavr.Host = host

	denonavr.mainZoneData = DenonXML{}

	denonavr.mainZoneStatus = DenonStatus{}
	denonavr.zone2Status = DenonStatus{}
	denonavr.zone3Status = DenonStatus{}

	denonavr.entityChangedFunction = make(map[string][]func(interface{}))

	denonavr.updateTrigger = make(chan string)

	return &denonavr
}

func (d *DenonAVR) AddHandleEntityChangeFunc(key string, f func(interface{})) {

	d.entityChangedFunction[key] = append(d.entityChangedFunction[key], f)

}

func (d *DenonAVR) getMainZoneDataFromDevice() {

	d.mainZoneData = DenonXML{} // Somehow the values in the array are added instead of replaced. Not sure if this is the solution, but it works...
	resp, err := http.Get("http://" + d.Host + MAINZONE_URL)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("Cannot read response body")
	}

	if err := xml.Unmarshal(body, &d.mainZoneData); err != nil {
		log.WithError(err).Info("Could not unmarshall")
	}
}

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

func (d *DenonAVR) StartListenLoop() {

	log.Info("Start Denon Listen Loop")

	updateInterval := 5 * time.Second
	ticker := time.NewTicker(updateInterval)

	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-d.updateTrigger:
			// force manual update
			d.updateAndNotify()
		case <-ticker.C:
			// Update every 5 Seconds
			d.updateAndNotify()

		}
	}
}

func (d *DenonAVR) updateAndNotify() {

	// Make copy of data to compare and update on changes
	oldMainZonData := d.mainZoneData
	oldMainZoneStatus := d.mainZoneStatus
	oldZone2Status := d.zone2Status
	oldZone3Status := d.zone3Status

	d.getMainZoneDataFromDevice()
	d.getZoneStatus(MainZone)
	d.getZoneStatus(Zone2)
	d.getZoneStatus(Zone3)

	// TODO: make the following part nicer?
	// Power
	if len(d.entityChangedFunction["Power"]) > 0 {
		if oldMainZonData.Power != d.mainZoneData.Power {
			for _, f := range d.entityChangedFunction["Power"] {
				f(d.mainZoneData.Power)
			}
		}
	}

	// Zone Power
	if len(d.entityChangedFunction["MainZonePower"]) > 0 {
		if oldMainZoneStatus.Power != d.mainZoneStatus.Power {
			for _, f := range d.entityChangedFunction["MainZonePower"] {
				f(d.mainZoneStatus.Power)
			}
		}
	}

	if len(d.entityChangedFunction["Zone2Power"]) > 0 {
		if oldZone2Status.Power != d.zone2Status.Power {
			for _, f := range d.entityChangedFunction["Zone2Power"] {
				f(d.zone2Status.Power)
			}
		}
	}

	if len(d.entityChangedFunction["Zone3Power"]) > 0 {
		if oldZone3Status.Power != d.zone3Status.Power {
			for _, f := range d.entityChangedFunction["Zone3Power"] {
				f(d.zone3Status.Power)
			}
		}
	}

	// Volume
	if len(d.entityChangedFunction["MainZoneVolume"]) > 0 {
		if oldMainZoneStatus.MasterVolume != d.mainZoneStatus.MasterVolume {
			for _, f := range d.entityChangedFunction["MainZoneVolume"] {
				f(d.mainZoneData.MasterVolume)
			}
		}
	}

	if len(d.entityChangedFunction["Zone2Volume"]) > 0 {
		if oldZone2Status.MasterVolume != d.zone2Status.MasterVolume {
			for _, f := range d.entityChangedFunction["Zone2Volume"] {
				f(d.zone2Status.MasterVolume)
			}
		}
	}

	if len(d.entityChangedFunction["Zone3Volume"]) > 0 {
		if oldZone3Status.MasterVolume != d.zone3Status.MasterVolume {
			for _, f := range d.entityChangedFunction["Zone23olume"] {
				f(d.zone3Status.MasterVolume)
			}
		}
	}

	if len(d.entityChangedFunction["MainZoneMute"]) > 0 {
		if oldMainZoneStatus.Mute != d.mainZoneStatus.Mute {
			for _, f := range d.entityChangedFunction["MainZoneMute"] {
				f(d.mainZoneStatus.Mute)
			}
		}
	}

	if len(d.entityChangedFunction["Zone2Mute"]) > 0 {
		if oldZone2Status.Mute != d.zone2Status.Mute {
			for _, f := range d.entityChangedFunction["Zone2Mute"] {
				f(d.zone2Status.Mute)
			}
		}
	}

	if len(d.entityChangedFunction["Zone3Mute"]) > 0 {
		if oldZone3Status.Mute != d.zone3Status.Mute {
			for _, f := range d.entityChangedFunction["Zone3Mute"] {
				f(d.zone3Status.Mute)
			}
		}
	}

	// Video Select

	if len(d.entityChangedFunction["MainZoneInputFuncList"]) > 0 {
		if !reflect.DeepEqual(oldMainZoneStatus.InputFuncList, d.mainZoneStatus.InputFuncList) {
			for _, f := range d.entityChangedFunction["MainZoneInputFuncList"] {
				f(d.mainZoneData.VideoSelectList)
			}
		}
	}

	if len(d.entityChangedFunction["MainZoneInputFuncSelect"]) > 0 {
		if oldMainZonData.VideoSelect != d.mainZoneData.VideoSelect {
			for _, f := range d.entityChangedFunction["MainZoneInputFuncSelect"] {
				f(d.mainZoneData.VideoSelect)
			}
		}
	}

	// Surround Mode
	if len(d.entityChangedFunction["MainZoneSurroundMode"]) > 0 {
		if oldMainZoneStatus.SurrMode != d.mainZoneStatus.SurrMode {
			for _, f := range d.entityChangedFunction["MainZoneSurroundMode"] {
				f(d.GetSurroundMode(MainZone))
			}
		}
	}

	if len(d.entityChangedFunction["Zone2SurroundMode"]) > 0 {
		if oldZone2Status.SurrMode != d.zone2Status.SurrMode {
			for _, f := range d.entityChangedFunction["Zone2SurroundMode"] {
				f(d.GetSurroundMode(Zone2))
			}
		}
	}

	if len(d.entityChangedFunction["Zone3SurroundMode"]) > 0 {
		if oldZone3Status.SurrMode != d.zone3Status.SurrMode {
			for _, f := range d.entityChangedFunction["Zone3SurroundMode"] {
				f(d.GetSurroundMode(Zone3))
			}
		}
	}

}

func (d *DenonAVR) GetFriendlyName() string {

	return d.mainZoneData.FriendlyName
}
