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
	zoneStatus     map[DenonZone]DenonZoneStatus
	netAudioStatus DenonNetAudioStatus

	attributes map[string]interface{}

	updateTrigger chan string

	entityChangedFunction map[string][]func(interface{})
}

func NewDenonAVR(host string) *DenonAVR {

	denonavr := DenonAVR{}

	denonavr.Host = host

	denonavr.mainZoneData = DenonXML{}
	denonavr.zoneStatus = make(map[DenonZone]DenonZoneStatus)
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
			go f(newValue)
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
			// force manual update
			d.updateAndNotify()
		case <-ticker.C:
			// Update every 5 Seconds
			d.updateAndNotify()

		}
	}
}

func (d *DenonAVR) updateAndNotify() {

	// Don't wait on each Call, handle them individually
	go d.updateMainZoneDataAndNotify()
	go d.updateZoneStatusAndNotify(MainZone)
	go d.updateZoneStatusAndNotify(Zone2)
	go d.updateZoneStatusAndNotify(Zone3)
}

func (d *DenonAVR) updateMainZoneDataAndNotify() {
	// Make copy of data to compare and update on changes
	oldMainZoneData := d.mainZoneData
	d.getMainZoneDataFromDevice()

	// Power
	if oldMainZoneData.Power != d.mainZoneData.Power {
		d.callEntityChangeFunction("POWER", d.mainZoneData.Power)
	}

	d.getNetAudioStatus()

	// Media Title
	d.getMediaTitle()

	// Media Image URL
	d.getMediaImageURL()
}

func (d *DenonAVR) updateZoneStatusAndNotify(zone DenonZone) {
	// Make copy of data to compare and update on changes
	oldZoneStatus := d.getZoneStatus(zone)
	// Get Data from Denon AVR
	newZoneStatus := d.getZoneStatus(zone)
	zoneName := d.getZoneName(zone)

	if oldZoneStatus.Power != newZoneStatus.Power {
		d.callEntityChangeFunction(zoneName+"Power", newZoneStatus.Power)
	}
	if oldZoneStatus.MasterVolume != newZoneStatus.MasterVolume {
		d.callEntityChangeFunction(zoneName+"Volume", newZoneStatus.MasterVolume)
	}
	if oldZoneStatus.Mute != newZoneStatus.Mute {
		d.callEntityChangeFunction(zoneName+"MainZoneMute", newZoneStatus.Mute)
	}

	// Input Func
	if !reflect.DeepEqual(oldZoneStatus.InputFuncList, newZoneStatus.InputFuncList) {

		var sourceList []string
		inputFuncSelectList := d.GetZoneInputFuncList(zone)
		for _, renamedSource := range inputFuncSelectList {
			sourceList = append(sourceList, renamedSource)
		}

		d.callEntityChangeFunction(zoneName+"InputFuncList", sourceList)
	}
	if oldZoneStatus.InputFuncSelect != newZoneStatus.InputFuncSelect {

		inputFuncSelect := newZoneStatus.InputFuncSelect

		// Rename Source with the SOURCE_MAPPING if necessary
		for source, origin := range SOURCE_MAPPING {
			if origin == newZoneStatus.InputFuncSelect {
				inputFuncSelect = source
				break
			}
		}
		// And then custom renames
		inputFuncSelectList := d.GetZoneInputFuncList(zone)
		if inputFuncSelectList[inputFuncSelect] != "" {
			inputFuncSelect = inputFuncSelectList[inputFuncSelect]
		}

		d.callEntityChangeFunction(zoneName+"InputFuncSelect", inputFuncSelect)
	}

	if oldZoneStatus.SurrMode != newZoneStatus.SurrMode {
		d.callEntityChangeFunction(zoneName+"SurroundMode", strings.TrimRight(newZoneStatus.SurrMode, " "))
	}
}

func (d *DenonAVR) getZoneName(zone DenonZone) string {

	switch zone {
	case MainZone:
		return "MainZone"

	case Zone2:
		return "Zone2"

	case Zone3:
		return "Zone3"
	}

	return ""
}
