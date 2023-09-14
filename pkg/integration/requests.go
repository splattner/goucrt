package integration

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/splattner/goucrt/pkg/entities"
	"k8s.io/utils/strings/slices"
)

// Handle the request Message from Remote Two
func (i *Integration) handleRequest(req *RequestMessage, p []byte) interface{} {
	log.Println("Handle Request Message: " + req.Msg)

	var res interface{}

	switch req.Msg {
	case "auth_required":
		authRequiredReq := AuthRequestMessage{}
		json.Unmarshal(p, &authRequiredReq)

		// TODO
		//res = i.handleAuthRequired(&authRequiredReq)

	case "get_driver_version":
		driverVersionReq := DriverVersionReq{}
		json.Unmarshal(p, &driverVersionReq)

		res = i.handleGetDriverVersionRequest(&driverVersionReq)

	case "get_driver_metadata":
		driverMetadataReq := DriverMetadataReq{}
		json.Unmarshal(p, &driverMetadataReq)

		res = i.getDriverMetadata(&driverMetadataReq)

	case "get_device_state":
		deviceStateMessageReq := DeviceStateMessageReq{}
		json.Unmarshal(p, &deviceStateMessageReq)

		i.handleGetDeviceStateRequest(&deviceStateMessageReq)

	case "get_available_entities":
		availableEntityMessageReq := AvailableEntityMessageReq{}
		json.Unmarshal(p, &availableEntityMessageReq)

		res = i.handleGetAvailableEntitiesRequest(&availableEntityMessageReq)
	case "subscribe_events":
		subscribeEventMessageReq := SubscribeEventMessageReq{}
		json.Unmarshal(p, &subscribeEventMessageReq)

		res = i.handleSubscribeEventRequest(&subscribeEventMessageReq)
	case "unsubscribe_events":
		unsubscribeEventMessageReq := UnubscribeEventMessageReq{}
		json.Unmarshal(p, &unsubscribeEventMessageReq)

		res = i.handleUnsubscribeEventsRequest(&unsubscribeEventMessageReq)

	case "get_entity_states":
		entityStatesReq := GetEntityStatesMessageReq{}
		json.Unmarshal(p, &entityStatesReq)

		res = i.getEntityStates(&entityStatesReq)

	case "entity_command":
		entityCommandReq := EntityCommandReq{}
		json.Unmarshal(p, &entityCommandReq)

		res = i.handleEntityCommandRequest(&entityCommandReq)

	case "setup_driver":
		setupDriverReq := SetupDriverMessageReq{}
		json.Unmarshal(p, &setupDriverReq)

		res = i.handleSetupDriverRequest(&setupDriverReq)

	case "set_driver_user_data":
		setUserData := SetDriverUserDataRequest{}
		json.Unmarshal(p, &setUserData)

		res = i.handleSetUserDataRequest(&setUserData)

	default:
		log.Println("mesage not know")
	}

	if res != nil {
		if err := i.sendResponseMessage(&res, websocket.TextMessage); err != nil {
			log.Println(err)
		}
	}

	return res
}

// Called by the Remote Two when it needs to synchronize the device state,
// e.g. after waking up from standby, or if it doesn't receive regular device_state events.
func (i *Integration) handleGetDeviceStateRequest(req *DeviceStateMessageReq) {
	i.sendDeviceStateEvent()
}

// Get version information about the integration driver.
func (i *Integration) handleGetDriverVersionRequest(req *DriverVersionReq) *ResponseMessage {

	msg_data := DriverVersionData{
		Name: i.Metadata.Name.En,
		Version: Version{
			Api:    "test",
			Driver: i.Metadata.Version,
		},
	}

	res := ResponseMessage{
		CommonResp{
			Kind: "resp",
			Id:   req.Id,
			Msg:  "driver_version",
			Code: 200,
		},
		msg_data,
	}

	return &res

}

// The metadata is used to setup the driver in the remote / web-configurator and start the setup flow.
func (i *Integration) getDriverMetadata(req *DriverMetadataReq) *DriverMetadataReponse {

	res := DriverMetadataReponse{
		CommonResp{
			Kind: "resp",
			Id:   req.Id,
			Msg:  "driver_metadata",
			Code: 200,
		},
		i.Metadata,
	}

	return &res

}

// Called while configuring profiles and assigning entities to pages or groups in the web-configurator or the embedded editor of the remote UI.
// With the optional filter, only entities of a given type can be requested.
func (i *Integration) handleGetAvailableEntitiesRequest(req *AvailableEntityMessageReq) *AvailableEntityMessage {

	// TODO: get filtered entities

	res := AvailableEntityMessage{
		CommonResp{Kind: "resp", Id: req.Id, Msg: "available_entities", Code: 200},
		AvailableEntityData{
			Filter:            req.MsgData.Filter,
			AvailableEntities: i.Entities,
		},
	}

	return &res

}

// start driver setup
// https://studio.asyncapi.com/?url=https://raw.githubusercontent.com/unfoldedcircle/core-api/main/integration-api/asyncapi.yaml#message-setup_driver
func (i *Integration) handleSetupDriverRequest(req *SetupDriverMessageReq) *ResponseMessage {
	if i.handleSetupFunction != nil {
		// The handleSetupFunction is where the driver specific implmenentation for driver setup is
		go i.handleSetupFunction()
	}

	res := ResponseMessage{
		CommonResp{
			Kind: "resp",
			Id:   req.Id,
			Msg:  req.Msg,
			Code: 200,
		},
		nil,
	}

	return &res

}

// Subscribe to entity state change events to receive entity_change events from the integration driver.
// If no entity IDs are specified then events for all available entities are sent to the Remote Two.
func (i *Integration) handleSubscribeEventRequest(req *SubscribeEventMessageReq) *SubscribeEventMessage {

	// Add entities to SubscribedEntities if not already in there
	for _, e := range i.Entities {
		if req.MsgData.EntityIds == nil || slices.Contains(req.MsgData.EntityIds, e.(entities.Entity).Id) {
			if !slices.Contains(i.SubscribedEntities, e.(entities.Entity).Id) {
				i.SubscribedEntities = append(i.SubscribedEntities, e.(entities.Entity).Id)
			}
		}
	}

	res := SubscribeEventMessage{
		CommonResp{Kind: "resp", Id: req.Id, Msg: req.Msg, Code: 200},
	}

	return &res

}

// If no entity IDs are specified then all events for all available entities are stopped.
// This message is sent by the Remote Two if a previously configured entity is no longer used and therefore no longer interested in entity updates. If the integration driver keeps sending events for the unsubscribed entities then they are simply discarded.
func (i *Integration) handleUnsubscribeEventsRequest(req *UnubscribeEventMessageReq) *UnubscribeEventMessage {

	for ix, e := range i.SubscribedEntities {
		if req.MsgData.EntityIds == nil || slices.Contains(i.SubscribedEntities, e) {

			i.SubscribedEntities[ix] = i.SubscribedEntities[len(i.SubscribedEntities)-1] // Copy last element to index i.
			i.SubscribedEntities[len(i.SubscribedEntities)-1] = ""                       // Erase last element (write zero value).
			i.SubscribedEntities = i.SubscribedEntities[:len(i.SubscribedEntities)-1]    // Truncate slice.
		}
	}

	res := UnubscribeEventMessage{
		CommonResp{Kind: "resp", Id: req.Id, Msg: req.Msg, Code: 200},
	}

	return &res
}

// Called by the Remote Two when it needs to synchronize the dynamic entity attributes, e.g. after connection setup or waking up from standby.
func (i *Integration) getEntityStates(req *GetEntityStatesMessageReq) *GetEntityStatesMessage {

	var entityStates []entities.EntityStateData

	for _, e := range i.Entities {
		entity := e.(entities.Entity)
		entityStates = append(entityStates, *entity.GetEntityState())
	}

	res := GetEntityStatesMessage{
		CommonResp{Kind: "resp", Id: req.Id, Msg: "entity_states", Code: 200},
		entityStates,
	}

	return &res
}

// Handle the entity command request sent by the remote
func (i *Integration) handleEntityCommandRequest(req *EntityCommandReq) *EntityCommandResponse {

	entity, _, err := i.GetEntityById(req.MsgData.EntityId)

	var returnCode int

	if err == nil {
		returnCode = 200
	} else {
		returnCode = 404
	}

	HandleCommand(entity, req)

	res := EntityCommandResponse{
		CommonResp{Kind: "resp", Id: req.Id, Msg: req.Msg, Code: returnCode},
	}

	return &res

}

func (i *Integration) handleSetUserDataRequest(req *SetDriverUserDataRequest) *EntityCommandResponse {

	if req.MsgData.InputValues != nil {
		i.UserInputValues = req.MsgData.InputValues
	}

	i.UserInputConfirmation = req.MsgData.Confirm

	res := EntityCommandResponse{
		CommonResp{Kind: "resp", Id: req.Id, Msg: req.Msg, Code: 200},
	}

	return &res
}