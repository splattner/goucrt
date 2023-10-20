package denonavr

import (
	"encoding/xml"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Zones
var ALL_ZONES = "All"
var MAIN_ZONE = "Main"
var ZONE2 = "Zone2"
var ZONE3 = "Zone3"

type DenonAVRDeviceInfo struct {
	api DenonAVRApi

	receiver         ReceiverType
	urls             ReceiverURLs
	Zone             string
	FriendlyName     string
	Manufacturer     string
	ModelName        string
	SerialNumber     string
	UseAvr2016Update bool
	power            string

	isSetup bool
}

func (deviceInfo DenonAVRDeviceInfo) getOwnZone() string {
	if deviceInfo.Zone == MAIN_ZONE {
		return "zone1"
	}
	return deviceInfo.Zone
}

func (deviceInfo DenonAVRDeviceInfo) setup() {

	switch deviceInfo.Zone {
	case MAIN_ZONE:
		deviceInfo.urls = DENONVAR_URLS
	case ZONE2:
		deviceInfo.urls = ZONE2_URLS
	case ZONE3:
		deviceInfo.urls = ZONE3_URLS
	}

	deviceInfo.identifyReceiver()
	//deviceInfo.getDeviceInfo()

	deviceInfo.isSetup = true
}

func (deviceInfo DenonAVRDeviceInfo) update() {

	if !deviceInfo.isSetup {
		deviceInfo.setup()
	}

	deviceInfo.updatePower()
}

func (deviceInfo DenonAVRDeviceInfo) identifyReceiver() {

	r_types := []ReceiverType{
		AVR_X,
		AVR_X_2016,
	}

	for _, r_type := range r_types {
		deviceInfo.api.Port = r_type.Port

		deviceInfoXML, err := deviceInfo.getDeviceInfo(deviceInfo.urls["deviceInfo"])

		if err == nil {
			if deviceInfo.isAVRX(*deviceInfoXML) {
				deviceInfo.receiver = r_type
			}
		}
	}

	//iff check of Deviceinfo.xml was not successful, receiver is type AVR
	deviceInfo.receiver = AVR
	deviceInfo.api.Port = AVR.Port

}

// TODO
func (deviceInfo DenonAVRDeviceInfo) getDeviceInfo(url string) (*DeviceInfoXML, error) {
	var deviceInfoXML DeviceInfoXML

	resp, err := http.Get("http://" + deviceInfo.api.Host + url)
	if err != nil {
		log.WithError(err).Error("Cannot make Get call")
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("Cannot read response body")
		return nil, err
	}

	if err := xml.Unmarshal(body, &deviceInfoXML); err != nil {
		log.WithError(err).Error("Could not unmarshall")
		return nil, err
	}

	return &deviceInfoXML, nil
}

// TODO
func (deviceInfo DenonAVRDeviceInfo) isAVRX(deviceInfoXML DeviceInfoXML) bool {

	if DEVICEINFO_COMMAPI_PATTERN.MatchString(deviceInfoXML.CommApiVers) {
		return true
	}

	if DEVICEINFO_AVR_X_PATTERN.MatchString(deviceInfo.ModelName) {
		return true
	}

	return false

}

// TODO
func (deviceInfo DenonAVRDeviceInfo) updatePower() {

}

func (d DenonAVR) setAPIHost(host string) {
	d.Device.api.Host = host
}

func (d DenonAVR) setAPITimeout(timeout int) {
	d.Device.api.timeout = timeout
}
