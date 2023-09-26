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
	Lights           []string       `json:"lights,omitempty"`
	LightSequence    []string       `json:"lightsequence,omitempty"`
	MultiDeviceIDs   []string       `json:"multideviceids,omitempty"`
	DeviceMembership []string       `json:"devicemembership,omitempty"`
	Scenes           []scenes.Scene `json:"scenes,omitempty"`
	State            DeconzState    `json:"state,omitempty"`
}

// Return all Groups from DeCONZ Rest API
func (d *Deconz) GetAllGroups() ([]DeconzGroup, error) {
	url := fmt.Sprintf(getAllGroupsURL, fmt.Sprintf("%s:%d", d.host, d.port), d.apikey)
	request, err := http.NewRequest("GET", url, nil)
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
	groupsMap := map[string]DeconzGroup{}
	json.Unmarshal(contents, &groupsMap)
	groups := make([]DeconzGroup, 0, len(groupsMap))
	for groupID, group := range groupsMap {
		group.TID = groupID
		group.ID, err = strconv.Atoi(groupID)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	sort.Slice(groups, func(i, j int) bool { return groups[i].ID < groups[j].ID })

	return groups, err
}

func (d *DeconzDevice) GetGroupAttrs(groupID int) (DeconzGroup, error) {
	var gg DeconzGroup
	url := fmt.Sprintf(getGroupAttrsURL, fmt.Sprintf("%s:%d", d.deconz.host, d.deconz.port), d.deconz.apikey, groupID)
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

func (d *DeconzDevice) SetGroupAttrs(groupID int, group DeconzGroup) ([]ApiResponse, error) {
	var apiResponse []ApiResponse
	url := fmt.Sprintf(setGroupAttrsURL, fmt.Sprintf("%s:%d", d.deconz.host, d.deconz.port), d.deconz.apikey, groupID)
	gg := DeconzGroup{}
	gg.Name = group.Name
	gg.Lights = group.Lights
	gg.Hidden = group.Hidden
	gg.LightSequence = group.LightSequence
	gg.MultiDeviceIDs = group.MultiDeviceIDs
	jsonData, err := json.Marshal(&gg)
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
	json.Unmarshal(contents, &apiResponse)
	return apiResponse, err
}

func (d *DeconzDevice) SetGroupState(groupID int, state DeconzState) ([]ApiResponse, error) {

	log.WithFields(log.Fields{"ID": groupID, "State": state}).Debug("Set Group State")

	var apiResponse []ApiResponse
	url := fmt.Sprintf(setGroupStateURL, fmt.Sprintf("%s:%d", d.deconz.host, d.deconz.port), d.deconz.apikey, groupID)
	jsonData, err := json.Marshal(&state)
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

	log.WithFields(log.Fields{
		"ID":    d.Group.ID,
		"State": d.Group.Action,
	}).Info("Deconz, call SetGroupState")

	_, err := d.SetGroupState(d.Group.ID, d.Group.Action)
	if err != nil {
		log.Debugln("Deconz, SetGroupState Error", err)
		return err
	}

	return nil
}
