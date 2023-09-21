package denonavr

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"

	log "github.com/sirupsen/logrus"

	"net/http"
)

type DenonCommand string

const (
	DenonCommandPower       DenonCommand = "PV"
	DenonCommandVolume                   = "MV"
	DenonCommandMute                     = "MU"
	DenonCommandSelectInput              = "SI"
	DenonVolumeStep         float64      = 1
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
	VideoSelectLists []ValueLists `xml:"VideoSelectLists>value"`
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
	Host       string
	baseURL    string
	commandURL string

	data DenonXML

	updateTrigger chan string

	entityChangedFunction map[string][]func(interface{})
}

func NewDenonAVR(host string) *DenonAVR {

	denonavr := DenonAVR{}

	denonavr.Host = host

	denonavr.baseURL = "http://" + host + "/goform/formMainZone_MainZoneXml.xml"
	denonavr.commandURL = "http://" + host + "/goform/formiPhoneAppDirect.xml"

	denonavr.data = DenonXML{}
	denonavr.entityChangedFunction = make(map[string][]func(interface{}))

	denonavr.updateTrigger = make(chan string)

	return &denonavr
}

func (d *DenonAVR) AddHandleEntityChangeFunc(key string, f func(interface{})) {

	d.entityChangedFunction[key] = append(d.entityChangedFunction[key], f)

}

func (d *DenonAVR) getDataFromDevice() {

	d.data = DenonXML{} // Somehow the values in the array are added instead of replaced. Not sure if this is the solution, but it works...
	resp, err := http.Get(d.baseURL)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)

	if err := xml.Unmarshal(body, &d.data); err != nil {
		log.WithError(err).Info("Could not unmarshall")
	}
}

func (d *DenonAVR) sendCommandToDevice(denonCommandType DenonCommand, command string) error {

	url := d.commandURL + "?" + string(denonCommandType) + command
	log.WithFields(log.Fields{
		"type":    string(denonCommandType),
		"command": command,
		"url":     url}).Info("Send Command to Denon Device")

	_, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error sending command: %w", err)
	}

	// Trigger a updata data, handeld in the Listen Looü
	d.updateTrigger <- "update"

	return nil
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
			d.UpdateAndNotify()
		case <-ticker.C:
			// Update every 5 Seconds
			d.UpdateAndNotify()

		}
	}
}

func (d *DenonAVR) UpdateAndNotify() {
	oldData := d.data
	d.getDataFromDevice()

	// TODO: make the following part nicer?

	if len(d.entityChangedFunction["MasterVolume"]) > 0 {
		if oldData.MasterVolume != d.data.MasterVolume {
			for _, f := range d.entityChangedFunction["MasterVolume"] {
				f(d.data.MasterVolume)
			}
		}

	}

	if len(d.entityChangedFunction["Power"]) > 0 {
		if oldData.Power != d.data.Power {
			for _, f := range d.entityChangedFunction["Power"] {
				f(d.data.Power)
			}
		}
	}

	if len(d.entityChangedFunction["ZonePower"]) > 0 {
		if oldData.ZonePower != d.data.ZonePower {
			for _, f := range d.entityChangedFunction["ZonePower"] {
				f(d.data.ZonePower)
			}
		}
	}

	if len(d.entityChangedFunction["Mute"]) > 0 {
		if oldData.Mute != d.data.Mute {
			for _, f := range d.entityChangedFunction["Mute"] {
				f(d.data.Mute)
			}
		}
	}

	if len(d.entityChangedFunction["VideoSelectLists"]) > 0 {
		if !EqualValueList(oldData.VideoSelectLists, d.data.VideoSelectLists) {
			for _, f := range d.entityChangedFunction["VideoSelectLists"] {
				f(d.data.VideoSelectLists)
			}
		}
	}

	if len(d.entityChangedFunction["VideoSelect"]) > 0 {
		if oldData.VideoSelect != d.data.VideoSelect {
			for _, f := range d.entityChangedFunction["VideoSelect"] {
				f(d.data.VideoSelect)
			}
		}
	}
}

func (d *DenonAVR) GetFriendlyName() string {

	return d.data.FriendlyName
}

func (d *DenonAVR) GetVideoSelectList() []ValueLists {
	return d.data.VideoSelectLists
}
