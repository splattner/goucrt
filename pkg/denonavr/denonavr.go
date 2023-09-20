package denonavr

import (
	"encoding/xml"
	"io"
	"time"

	log "github.com/sirupsen/logrus"

	"net/http"
)

type DenonXML struct {
	XMLName          xml.Name `xml:"item"`
	FriendlyName     string   `xml:"FriendlyName>value"`
	ZonePower        string   `xml:"ZonePower>value"`
	RenameZone       string   `xml:"RenameZone>value"`
	TopMenuLink      string   `xml:"TopMenuLink>value"`
	VideoSelectDisp  string   `xml:"VideoSelectDisp>value"`
	VideoSelect      string   `xml:"VideoSelect>value"`
	VideoSelectOnOff string   `xml:"VideoSelectOnOff>value"`
	// TODO: VideoSelectLists
	ECOModeDisp string `xml:"ECOModeDisp>value"`
	ECOMode     string `xml:"ECOMode>value"`
	// TODO: ECOModeLists
	AddSourceDisplay string `xml:"AddSourceDisplay>value"`
	ModelId          string `xml:"ModelId>value"`
	BrandId          string `xml:"BrandId>value"`
	SalesArea        string `xml:"SalesArea>value"`
	InputFuncSelect  string `xml:"InputFuncSelect>value"`
	NetFuncSelect    string `xml:"NetFuncSelect>value"`
	SelectSurround   string `xml:"selectSurround>value"`
	VolumeDisplay    string `xml:"VolumeDisplay>value"`
	MasterVolume     string `xml:"MasterVolume>value"`
	Mute             string `xml:"Mute>value"`
}

type DenonXMLValue struct {
	Value string `xml:"value"`
}

type DenonAVR struct {
	host    string
	baseURL string

	data DenonXML

	entityChangedFunction map[string][]func(string)
}

func NewDenonAVR(host string) *DenonAVR {

	denonavr := DenonAVR{}

	denonavr.host = host

	denonavr.baseURL = "http://" + host + "/goform/formMainZone_MainZoneXml.xml"

	denonavr.data = DenonXML{}
	denonavr.entityChangedFunction = make(map[string][]func(string))

	return &denonavr
}

func (d *DenonAVR) AddHandleEntityChangeFunc(key string, f func(string)) {

	d.entityChangedFunction[key] = append(d.entityChangedFunction[key], f)

}

func (d *DenonAVR) getDataFromDevice() {

	resp, err := http.Get(d.baseURL)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)

	if err := xml.Unmarshal(body, &d.data); err != nil {
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

	for {
		select {
		case <-ticker.C:
			oldData := d.data
			d.getDataFromDevice()

			if len(d.entityChangedFunction["MasterVolume"]) > 0 {
				if oldData.MasterVolume != d.data.MasterVolume {
					for _, f := range d.entityChangedFunction["MasterVolume"] {
						f(d.data.MasterVolume)
					}
				}

			}
		}
	}
}

func (d *DenonAVR) GetFriendlyName() string {

	return d.data.FriendlyName
}

func (d *DenonAVR) GetMasterVolume() string {

	return d.data.MasterVolume
}
