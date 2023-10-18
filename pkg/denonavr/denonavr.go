package denonavr

import (
	"encoding/xml"
	"io"
	"reflect"
	"strings"
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
	DenonVolumeStep           float64      = 0.5
)

const (
	MainZone DenonZone = "MAIN"
	Zone2    DenonZone = "Z2"
	Zone3    DenonZone = "Z3"
)

const (
	STATUS_URL           string = "/goform/formMainZone_MainZoneXmlStatus.xml"
	STATUS_Z2_URL        string = "/goform/formZone2_Zone2XmlStatus.xml"
	STATUS_Z3_URL        string = "/goform/formZone3_Zone3XmlStatus.xml"
	MAINZONE_URL         string = "/goform/formMainZone_MainZoneXml.xml"
	COMMAND_URL          string = "/goform/formiPhoneAppDirect.xml"
	NET_AUDIO_STATUR_URL string = "/goform/formNetAudio_StatusXml.xml"
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
	netAudioStatus DenonNetAudioStatus

	attributes map[string]interface{}

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
	denonavr.netAudioStatus = DenonNetAudioStatus{}

	denonavr.entityChangedFunction = make(map[string][]func(interface{}))

	denonavr.attributes = make(map[string]interface{})

	denonavr.updateTrigger = make(chan string)

	return &denonavr
}

// Add a new function that is called when a attribute of this entity has changed
func (d *DenonAVR) AddHandleEntityChangeFunc(attribute string, f func(interface{})) {
	d.entityChangedFunction[attribute] = append(d.entityChangedFunction[attribute], f)
}

// Call the registred entity change function with the new value for a attribute
func (d *DenonAVR) callEntityChangeFunction(attribute string, newValue interface{}) {
	if len(d.entityChangedFunction[attribute]) > 0 {
		for _, f := range d.entityChangedFunction[attribute] {
			f(newValue)
		}
	}
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

func (d *DenonAVR) StartListenLoop() {

	log.Info("Start Denon Listen Loop")

	updateInterval := 5 * time.Second
	ticker := time.NewTicker(updateInterval)

	defer func() {
		ticker.Stop()
	}()

	// do an intial update to make sure we have up to date values
	d.updateAndNotify()

	for {
		select {
		case <-d.updateTrigger:
			log.Debug("Update Trigger")
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
	oldMainZoneData := d.mainZoneData
	oldMainZoneStatus := d.mainZoneStatus
	oldZone2Status := d.zone2Status
	oldZone3Status := d.zone3Status

	// Get Data from Denon AVR
	d.getMainZoneDataFromDevice()
	d.getZoneStatus(MainZone)
	d.getZoneStatus(Zone2)
	d.getZoneStatus(Zone3)
	d.getNetAudioStatus()

	// TODO: make the following part nicer?

	// Media Title
	d.getMediaTitle()

	// Media Image URL
	d.getMediaImageURL()

	// Power
	if oldMainZoneData.Power != d.mainZoneData.Power {
		d.callEntityChangeFunction("POWER", d.mainZoneData.Power)
	}

	// Zone Power
	if oldMainZoneStatus.Power != d.mainZoneStatus.Power {
		d.callEntityChangeFunction("POWER", d.mainZoneStatus.Power)
	}
	if oldZone2Status.Power != d.zone2Status.Power {
		d.callEntityChangeFunction("Zone2Power", d.zone2Status.Power)
	}
	if oldZone3Status.Power != d.zone3Status.Power {
		d.callEntityChangeFunction("Zone3Power", d.zone3Status.Power)
	}

	// Volume
	if oldMainZoneStatus.MasterVolume != d.mainZoneStatus.MasterVolume {
		d.callEntityChangeFunction("MainZoneVolume", d.mainZoneStatus.MasterVolume)
	}
	if oldZone2Status.MasterVolume != d.zone2Status.MasterVolume {
		d.callEntityChangeFunction("Zone2Volume", d.zone2Status.MasterVolume)
	}
	if oldZone3Status.MasterVolume != d.zone3Status.MasterVolume {
		d.callEntityChangeFunction("Zone3Volume", d.zone3Status.MasterVolume)
	}
	if oldMainZoneStatus.Mute != d.mainZoneStatus.Mute {
		d.callEntityChangeFunction("MainZoneMute", d.mainZoneStatus.Mute)
	}
	if oldZone2Status.Mute != d.zone2Status.Mute {
		d.callEntityChangeFunction("Zone2Mute", d.zone2Status.Mute)
	}
	if oldZone3Status.Mute != d.zone3Status.Mute {
		d.callEntityChangeFunction("Zone3Mute", d.zone3Status.Mute)
	}

	// Video Select
	if !reflect.DeepEqual(oldMainZoneStatus.InputFuncList, d.mainZoneStatus.InputFuncList) {
		d.callEntityChangeFunction("MainZoneInputFuncList", d.mainZoneData.VideoSelectList)
	}
	if oldMainZoneData.VideoSelect != d.mainZoneData.VideoSelect {
		d.callEntityChangeFunction("MainZoneInputFuncSelect", d.mainZoneData.VideoSelect)
	}

	// Surround Mode
	if oldMainZoneStatus.SurrMode != d.mainZoneStatus.SurrMode {
		d.callEntityChangeFunction("MainZoneSurroundMode", strings.TrimRight(d.mainZoneStatus.SurrMode, " "))
	}
	if oldZone2Status.SurrMode != d.zone2Status.SurrMode {
		d.callEntityChangeFunction("Zone2SurroundMode", strings.TrimRight(d.zone2Status.SurrMode, " "))
	}
	if oldZone3Status.SurrMode != d.zone3Status.SurrMode {
		d.callEntityChangeFunction("Zone3SurroundMode", strings.TrimRight(d.zone3Status.SurrMode, " "))
	}

}
