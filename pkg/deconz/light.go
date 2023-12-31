package deconz

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	getAllLightsURL  = "http://%s/api/%s/lights"
	getLightStateURL = "http://%s/api/%s/lights/%d"
	setLightStateURL = "http://%s/api/%s/lights/%d/state"
	setLightAttrsURL = "http://%s/api/%s/lights/%d"
)

type DeconzLight struct {
	Name              string      `json:"name"`
	ID                int         `json:"id,omitempty"`
	ETag              string      `json:"etag,omitempty"`
	State             DeconzState `json:"state,omitempty"`
	HasColor          bool        `json:"hascolor,omitempty"`
	Type              string      `json:"type,omitempty"`
	ModelID           string      `json:"modelid,omitempty"`
	UniqueID          string      `json:"uniqueid,omitempty"`
	SWVersion         string      `json:"swversion,omitempty"`
	LastAnnounced     string      `json:"lastannounced,omitempty"`
	LastSeen          string      `json:"lastseen,omitempty"`
	ManufacturerName  string      `json:"manufacturername,omitempty"`
	ColorCapabilities int         `json:"colorcapabilities,omitempty"`
	Ctmax             int         `json:"ctmax,omitempty"`
	Ctmin             int         `json:"ctmin,omitempty"`
}

func (d *Deconz) GetLight(lightID int) (DeconzLight, error) {
	var ll DeconzLight
	url := fmt.Sprintf(getLightStateURL, fmt.Sprintf("%s:%d", d.host, d.port), d.apikey, lightID)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ll, err
	}
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return ll, err
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return ll, err
	}
	err = json.Unmarshal(contents, &ll)
	if err != nil {
		return ll, err
	}
	ll.ID = lightID
	return ll, err
}

func (d *DeconzDevice) SetLightAttrs() ([]ApiResponse, error) {
	url := fmt.Sprintf(setLightAttrsURL, fmt.Sprintf("%s:%d", d.deconz.host, d.deconz.port), d.deconz.apikey, d.Light.ID)
	data := fmt.Sprintf("{\"name\": \"%s\"}", d.Light.Name)
	postbody := strings.NewReader(data)
	request, err := http.NewRequest("PUT", url, postbody)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var apiResponse []ApiResponse
	err = json.Unmarshal(contents, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

func (d *DeconzDevice) SetLightState() ([]ApiResponse, error) {
	url := fmt.Sprintf(setLightStateURL, fmt.Sprintf("%s:%d", d.deconz.host, d.deconz.port), d.deconz.apikey, d.Light.ID)
	stateJSON, err := json.Marshal(&d.Light.State)
	if err != nil {
		return nil, err
	}
	postbody := strings.NewReader(string(stateJSON))
	request, err := http.NewRequest("PUT", url, postbody)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var apiResponse []ApiResponse
	err = json.Unmarshal(contents, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

func (d *Deconz) GetAllLights() ([]DeconzLight, error) {
	url := fmt.Sprintf(getAllLightsURL, fmt.Sprintf("%s:%d", d.host, d.port), d.apikey)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	lightsMap := map[string]DeconzLight{}
	err = json.Unmarshal(contents, &lightsMap)
	if err != nil {
		return nil, err
	}
	lights := make([]DeconzLight, 0, len(lightsMap))
	for lightID, light := range lightsMap {
		light.ID, _ = strconv.Atoi(lightID)
		lights = append(lights, light)
	}

	sort.Slice(lights, func(i, j int) bool { return lights[i].ID < lights[j].ID })
	return lights, err
}

func (d *DeconzDevice) newDeconzLightDevice() {

}

func (d *DeconzDevice) setLightState() error {

	log.WithFields(log.Fields{
		"ID":    d.Group.ID,
		"State": d.Light.State,
	}).Info("Deconz, call SetGroupState")

	_, err := d.SetLightState()
	if err != nil {
		log.Debugln("Deconz, SetLightState Error", err)
		return err
	}

	return nil
}
