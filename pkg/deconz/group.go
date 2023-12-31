package deconz

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/jurgen-kluft/go-conbee/scenes"
	log "github.com/sirupsen/logrus"
)

var (
	getAllGroupsURL  = "http://%s/api/%s/groups"
	getGroupAttrsURL = "http://%s/api/%s/groups/%d"
	setGroupAttrsURL = "http://%s/api/%s/groups/%d"
	setGroupStateURL = "http://%s/api/%s/groups/%d/action"
)

type DeconzGroup struct {
	ID               int
	TID              string         `json:"id,omitempty"`
	ETag             string         `json:"etag,omitempty"`
	Name             string         `json:"name,omitempty"`
	Hidden           bool           `json:"hidden,omitempty"`
	Action           DeconzState    `json:"action,omitempty"`
	LightIDs         []string       `json:"lights,omitempty"`
	LightSequence    []string       `json:"lightsequence,omitempty"`
	MultiDeviceIDs   []string       `json:"multideviceids,omitempty"`
	DeviceMembership []string       `json:"devicemembership,omitempty"`
	Scenes           []scenes.Scene `json:"scenes,omitempty"`
	State            DeconzState    `json:"state,omitempty"`
	Lights           []*DeconzDevice
}

// Return all Groups from DeCONZ Rest API
func (d *Deconz) GetAllGroups() ([]DeconzGroup, error) {
	url := fmt.Sprintf(getAllGroupsURL, fmt.Sprintf("%s:%d", d.host, d.port), d.apikey)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.WithError(err).Error("Cannot create new http Request")
		return nil, err
	}
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.WithError(err).Error("Cannot Do the request")
		return nil, err
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		log.WithError(err).Error("Cannot Read the response body")
		return nil, err
	}
	groupsMap := map[string]DeconzGroup{}
	if err := json.Unmarshal(contents, &groupsMap); err != nil {
		log.WithError(err).Error("Cannon unmarshall data into map[string]DeconzGroup")
	}
	groups := make([]DeconzGroup, 0, len(groupsMap))
	for groupID, group := range groupsMap {
		group.TID = groupID
		group.ID, err = strconv.Atoi(groupID)
		if err != nil {
			return nil, err
		}

		// Get State of all Lights in this group
		group.Lights = make([]*DeconzDevice, len(group.LightIDs))
		for i, id := range group.LightIDs {
			lightid, _ := strconv.Atoi(id)
			light, err := d.GetLight(lightid)

			if err != nil {
				return nil, err
			}
			lightDevice, err := d.GetDeviceByID(light.ID)
			if err != nil {
				return nil, err
			}
			group.Lights[i] = lightDevice
		}

		groups = append(groups, group)
	}

	sort.Slice(groups, func(i, j int) bool { return groups[i].ID < groups[j].ID })

	return groups, err
}

func (d *Deconz) GetGroup(groupID int) (DeconzGroup, error) {
	var gg DeconzGroup
	url := fmt.Sprintf(getGroupAttrsURL, fmt.Sprintf("%s:%d", d.host, d.port), d.apikey, groupID)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return gg, err
	}
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return gg, err
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return gg, err
	}
	gg.ID = groupID
	gg.TID = fmt.Sprintf("%d", groupID)
	err = json.Unmarshal(contents, &gg)
	if err != nil {
		return gg, err
	}
	return gg, err
}

func (d *DeconzDevice) SetGroupAttrs() ([]ApiResponse, error) {
	var apiResponse []ApiResponse
	url := fmt.Sprintf(setGroupAttrsURL, fmt.Sprintf("%s:%d", d.deconz.host, d.deconz.port), d.deconz.apikey, d.Group.ID)

	jsonData, err := json.Marshal(&d.Group)
	if err != nil {
		return apiResponse, err
	}
	body := strings.NewReader(string(jsonData))
	request, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return apiResponse, err
	}
	request.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return apiResponse, err
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return apiResponse, err
	}
	if err := json.Unmarshal(contents, &apiResponse); err != nil {
		log.WithError(err).Error("Cannon unmarshall data into []ApiResponse")
	}
	return apiResponse, err
}

func (d *DeconzDevice) SetGroupState() ([]ApiResponse, error) {

	log.WithFields(log.Fields{"ID": d.Group.ID, "State": d.Group.Action}).Debug("Set Group State")

	var apiResponse []ApiResponse
	url := fmt.Sprintf(setGroupStateURL, fmt.Sprintf("%s:%d", d.deconz.host, d.deconz.port), d.deconz.apikey, d.Group.ID)
	jsonData, err := json.Marshal(&d.Group.Action)
	if err != nil {
		return apiResponse, err
	}
	body := strings.NewReader(string(jsonData))
	request, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return apiResponse, err
	}
	request.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return apiResponse, err
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return apiResponse, err
	}
	err = json.Unmarshal(contents, &apiResponse)
	if err != nil {
		return apiResponse, err
	}
	return apiResponse, err
}

func (d *DeconzDevice) newDeconzGroupDevice() {

}

func (d *DeconzDevice) setGroupState() error {

	_, err := d.SetGroupState()
	if err != nil {
		log.WithError(err).Debug("Deconz, SetGroupState Error")
		return err
	}

	return nil
}
