package denonavr

import (
	"encoding/xml"
	"fmt"
	"hash/fnv"
	"io"
	"reflect"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/utils/strings/slices"

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

	media_title     string
	media_image_url string

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

// Call the registred entity change function with the new value for a attribute
func (d *DenonAVR) callEntityChangeFunction(attribute string, newValue interface{}) {
	if len(d.entityChangedFunction[attribute]) > 0 {
		for _, f := range d.entityChangedFunction[attribute] {
			f(newValue)
		}
	}
}

// Get the current Media Title
// Title of the Playing media or the current Input Function
func (d *DenonAVR) getMediaTitle() string {
	var media_title string

	if slices.Contains(PLAYING_SOURCES, d.mainZoneData.InputFuncSelect) {
		// This is a source that is playing audio
		// fot the moment, also set this to the input func
		media_title = d.mainZoneData.InputFuncSelect
	} else {
		// Not a playing source
		media_title = d.mainZoneData.InputFuncSelect
	}

	return media_title
}

// Get the current Media Title
// Title of the Playing media or the current Input Function
func (d *DenonAVR) getMediaImageURL() string {
	var media_image_url string

	if slices.Contains(PLAYING_SOURCES, d.mainZoneData.InputFuncSelect) {
		// This is a source that is playing audio
		// fot the moment, also set this to the input func

		hash := fnv.New32a()
		hash.Write([]byte(d.media_title))
		media_image_url = fmt.Sprintf("http://%s:%d/NetAudio/art.asp-jpg?%d", d.Host, 80, hash.Sum32())
	}

	return media_image_url
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

	// Media Title
	var media_title = d.getMediaTitle()
	if d.media_title != media_title {
		d.media_title = media_title
		d.callEntityChangeFunction("media_title", media_title)
	}

	// Media Image URL
	var media_image_url = d.getMediaImageURL()
	if d.media_image_url != media_image_url {
		d.media_image_url = media_image_url
		d.callEntityChangeFunction("media_image_url", media_image_url)
	}

	// Power
	if oldMainZonData.Power != d.mainZoneData.Power {
		d.callEntityChangeFunction("POWER", d.mainZoneData.Power)
	}

	// Zone Power
	if oldMainZoneStatus.Power != d.mainZoneData.Power {
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
		d.callEntityChangeFunction("MainZoneMute", d.mainZoneStatus.MasterVolume)
	}
	if oldZone2Status.Mute != d.zone2Status.Mute {
		d.callEntityChangeFunction("Zone2Mute", d.zone2Status.MasterVolume)
	}
	if oldZone3Status.Mute != d.zone3Status.Mute {
		d.callEntityChangeFunction("Zone3Mute", d.zone3Status.MasterVolume)
	}

	// Video Select
	if !reflect.DeepEqual(oldMainZoneStatus.InputFuncList, d.mainZoneStatus.InputFuncList) {
		d.callEntityChangeFunction("MainZoneInputFuncList", d.mainZoneData.VideoSelectList)
	}
	if oldMainZonData.VideoSelect != d.mainZoneData.VideoSelect {
		d.callEntityChangeFunction("MainZoneInputFuncSelect", d.mainZoneData.VideoSelect)
	}

	// Surround Mode
	if oldMainZoneStatus.SurrMode != d.mainZoneStatus.SurrMode {
		d.callEntityChangeFunction("MainZoneSurroundMode", strings.TrimLeft(d.mainZoneStatus.SurrMode, ""))
	}
	if oldZone2Status.SurrMode != d.zone2Status.SurrMode {
		d.callEntityChangeFunction("Zone2SurroundMode", strings.TrimLeft(d.zone2Status.SurrMode, ""))
	}
	if oldZone3Status.SurrMode != d.zone3Status.SurrMode {
		d.callEntityChangeFunction("Zone3SurroundMode", strings.TrimLeft(d.zone3Status.SurrMode, ""))
	}

}

func (d *DenonAVR) GetFriendlyName() string {

	return d.mainZoneData.FriendlyName
}
