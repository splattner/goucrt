package deconz

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

var (
	getAllSensorsURL = "http://%s/api/%s/sensors"
	getSensorURL     = "http://%s/api/%s/sensors/%d"
)

type DeconzSensor struct {
	ID               int
	Config           DeconzSensorConfig `json:"config,omitempty"`
	Ep               int                `json:"ep,omitempty"`
	ETag             string             `json:"etag"`
	ManufacturerName string             `json:"manufacturername,omitempty"`
	ModelID          string             `json:"modelid,omitempty"`
	Name             string             `json:"name"`
	State            DeconzState        `json:"state,omitempty"`
	SWVersion        string             `json:"swversion,omitempty"`
	Type             string             `json:"type,omitempty"`
	UniqueID         string             `json:"uniqueid,omitempty"`
}

type DeconzSensorConfig struct {
	On            bool   `json:"on"`
	Reachable     bool   `json:"reachable"`
	Battery       int16  `json:"battery,omitempty"`
	Long          string `json:"long,omitempty"`
	Lat           string `json:"lat,omitempty"`
	SunriseOffset int16  `json:"sunriseoffset,omitempty"`
	SunsetOffset  int16  `json:"sunsetoffset,omitempty"`
}

func (d *DeconzDevice) GetSensor(sensorID int) (DeconzSensor, error) {
	var ll DeconzSensor
	url := fmt.Sprintf(getSensorURL, fmt.Sprintf("%s:%d", d.deconz.host, d.deconz.port), d.deconz.apikey, sensorID)
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
	ll.ID = sensorID
	return ll, err
}

func (d *DeconzDevice) UpdateSensor(sensorID int, sensorName string) ([]ApiResponse, error) {
	url := fmt.Sprintf(getSensorURL, fmt.Sprintf("%s:%d", d.deconz.host, d.deconz.port), d.deconz.apikey, sensorID)
	data := fmt.Sprintf("{\"name\": \"%s\"}", sensorName)
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

func (d *Deconz) GetAllSensors() ([]DeconzSensor, error) {
	url := fmt.Sprintf(getAllSensorsURL, fmt.Sprintf("%s:%d", d.host, d.port), d.apikey)
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
	//fmt.Println(string(contents))

	sensorsMap := map[string]DeconzSensor{}
	err = json.Unmarshal(contents, &sensorsMap)
	if err != nil {
		return nil, err
	}
	sensors := make([]DeconzSensor, 0, len(sensorsMap))
	for sensorID, sensor := range sensorsMap {
		sensor.ID, _ = strconv.Atoi(sensorID)
		sensors = append(sensors, sensor)
	}

	sort.Slice(sensors, func(i, j int) bool { return sensors[i].ID < sensors[j].ID })

	return sensors, err
}

func (d *DeconzDevice) newDeconzSensorDevice() {

}
