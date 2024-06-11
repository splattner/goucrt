package integration

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
	"github.com/splattner/goucrt/pkg/entities"
	"k8s.io/utils/strings/slices"
)

// Handle the request Message from Remote Two
func (i *Integration) handleRequest(req *RequestMessage, p []byte) {
	var res interface{}

	log.WithField("RawMessage", string(p)).Debug("Request received")

	switch req.Msg {
	case "auth":
		authRequiredReq := AuthRequestMessage{}
		if err := json.Unmarshal(p, &authRequiredReq); err != nil {
			log.WithError(err).Error("Cannot unmarshall authRequiredReq")
		}

		// TODO
		//res = i.handleAuthRequired(&authRequiredReq)

	case "get_driver_version":
		driverVersionReq := DriverVersionReq{}
		if err := json.Unmarshal(p, &driverVersionReq); err != nil {
			log.WithError(err).Error("Cannot unmarshall DriverVersionReq")
		}

		res = i.handleGetDriverVersionRequest(&driverVersionReq)

	case "get_driver_metadata":
		driverMetadataReq := DriverMetadataReq{}
		if err := json.Unmarshal(p, &driverMetadataReq); err != nil {
			log.WithError(err).Error("Cannot unmarshall DriverMetadataReq")
		}

		res = i.handleGetDriverMetadataRequest(&driverMetadataReq)

	case "get_device_state":
		deviceStateMessageReq := DeviceStateMessageReq{}
		if err := json.Unmarshal(p, &deviceStateMessageReq); err != nil {
			log.WithError(err).Error("Cannot unmarshall deviceStateMessageReq")
		}

		i.handleGetDeviceStateRequest(&deviceStateMessageReq)

	case "get_available_entities":
		availableEntityMessageReq := AvailableEntityMessageReq{}
		if err := json.Unmarshal(p, &availableEntityMessageReq); err != nil {
			log.WithError(err).Error("Cannot unmarshall availableEntityMessageReq")
		}

		res = i.handleGetAvailableEntitiesRequest(&availableEntityMessageReq)
	case "subscribe_events":
		subscribeEventMessageReq := SubscribeEventMessageReq{}
		if err := json.Unmarshal(p, &subscribeEventMessageReq); err != nil {
			log.WithError(err).Error("Cannot unmarshall subscribeEventMessageReq")
		}

		res = i.handleSubscribeEventRequest(&subscribeEventMessageReq)
	case "unsubscribe_events":
		unsubscribeEventMessageReq := UnubscribeEventMessageReq{}
		if err := json.Unmarshal(p, &unsubscribeEventMessageReq); err != nil {
			log.WithError(err).Error("Cannot unmarshall unsubscribeEventMessageReq")
		}

		res = i.handleUnsubscribeEventsRequest(&unsubscribeEventMessageReq)

	case "get_entity_states":
		entityStatesReq := GetEntityStatesMessageReq{}
		if err := json.Unmarshal(p, &entityStatesReq); err != nil {
			log.WithError(err).Error("Cannot unmarshall entityStatesReq")
		}

		res = i.handleGetEntityStatesRequest(&entityStatesReq)

	case "entity_command":
		entityCommandReq := EntityCommandReq{}
		if err := json.Unmarshal(p, &entityCommandReq); err != nil {
			log.WithError(err).Error("Cannot unmarshall entityCommandReq")
		}

		res = i.handleEntityCommandRequest(&entityCommandReq)

	case "setup_driver":
		setupDriverReq := SetupDriverMessageReq{}
		if err := json.Unmarshal(p, &setupDriverReq); err != nil {
			log.WithError(err).Error("Cannot unmarshall setupDriverReq")
		}

		res = i.handleSetupDriverRequest(&setupDriverReq)

	case "set_driver_user_data":
		setUserData := SetDriverUserDataRequest{}
		if err := json.Unmarshal(p, &setUserData); err != nil {
			log.WithError(err).Error("Cannot unmarshall setUserData")
		}

		res = i.handleSetDriverUserDataRequest(&setUserData)

	default:
		log.Debug("mesage not know")
	}

	if res != nil {
		if err := i.sendResponseMessage(&res, websocket.TextMessage); err != nil {
			log.Error(err)
		}
	}
}

// Called by the Remote Two when it needs to synchronize the device state,
// e.g. after waking up from standby, or if it doesn't receive regular device_state events.
func (i *Integration) handleGetDeviceStateRequest(req *DeviceStateMessageReq) {

	// The response is a event Message and not a response
	i.sendDeviceStateEvent()
}

// Get version information about the integration driver.
func (i *Integration) handleGetDriverVersionRequest(req *DriverVersionReq) *ResponseMessage {

	msg_data := DriverVersionData{
		Name: i.Metadata.Name.En,
		Version: Version{
			Api:    API_VERSION,
			Driver: API_VERSION,
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
func (i *Integration) handleGetDriverMetadataRequest(req *DriverMetadataReq) *DriverMetadataReponse {

	res := DriverMetadataReponse{
		CommonResp{
			Kind: "resp",
			Id:   req.Id,
			Msg:  "driver_metadata",
			Code: 200,
		},
		*i.Metadata,
	}

	return &res

}

// Called while configuring profiles and assigning entities to pages or groups in the web-configurator or the embedded editor of the remote UI.
// With the optional filter, only entities of a given type can be requested.
func (i *Integration) handleGetAvailableEntitiesRequest(req *AvailableEntityMessageReq) interface{} {

	log.WithFields(log.Fields{
		"Id":      req.Id,
		"Kind":    req.Kind,
		"Msg":     req.Msg,
		"MsgData": req.MsgData,
	}).Debug("Get available Entities")

	var entities []interface{}

	var res interface{}

	for _, e := range i.Entities {
		if req.MsgData.Filter.EntityType.Type == "" || i.getEntityType(e).Type == req.MsgData.Filter.EntityType.Type {
			entities = append(entities, e)
		}
	}

	if req.MsgData.Filter.EntityType.Type == "" {
		res = AvailableEntityNoFilterMessage{
			CommonResp{Kind: "resp", Id: req.Id, Msg: "available_entities", Code: 200},
			AvailableEntityNoFilterData{
				AvailableEntities: entities,
			},
		}
	} else {
		res = AvailableEntityMessage{
			CommonResp{Kind: "resp", Id: req.Id, Msg: "available_entities", Code: 200},
			AvailableEntityData{
				Filter:            req.MsgData.Filter,
				AvailableEntities: entities,
			},
		}
	}

	return &res

}

// start driver setup
// https://studio.asyncapi.com/?url=https://raw.githubusercontent.com/unfoldedcircle/core-api/main/integration-api/asyncapi.yaml#message-setup_driver
func (i *Integration) handleSetupDriverRequest(req *SetupDriverMessageReq) *ResponseMessage {

	i.SetupData = req.MsgData.Value

	i.PersistSetupData()

	if i.handleSetupFunction != nil {
		// The handleSetupFunction is where the driver specific implmenentation for driver setup is
		go i.handleSetupFunction(req.MsgData.Value)
	}

	res := ResponseMessage{
		CommonResp{
			Kind: "resp",
			Id:   req.Id,
			Msg:  "result",
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
	if req.MsgData.EntityIds == nil {
		// Subscribe to all available entities
		for _, e := range i.Entities {
			entity_id := i.getEntityId(e)
			if !slices.Contains(i.SubscribedEntities, entity_id) {
				log.WithField("entity_id", entity_id).Info("RT subscribed to entity")
				i.SubscribedEntities = append(i.SubscribedEntities, entity_id)
				i.callSubscribeCallback(e)

			}
		}

	} else {
		for _, entity_id := range req.MsgData.EntityIds {
			if !slices.Contains(i.SubscribedEntities, entity_id) {
				log.WithField("entity_id", entity_id).Info("RT subscribed to entity")
				i.SubscribedEntities = append(i.SubscribedEntities, entity_id)

				if entity, _, err := i.GetEntityById(entity_id); err != nil {
					i.callSubscribeCallback(entity)
				}
			}
		}
	}

	log.WithField("subscribedEtities", i.SubscribedEntities).Debug("Change in subscribed entities")

	res := SubscribeEventMessage{
		CommonResp{Kind: "resp", Id: req.Id, Msg: "result", Code: 200},
	}

	return &res

}

// If no entity IDs are specified then all events for all available entities are stopped.
// This message is sent by the Remote Two if a previously configured entity is no longer used and therefore no longer interested in entity updates. If the integration driver keeps sending events for the unsubscribed entities then they are simply discarded.
func (i *Integration) handleUnsubscribeEventsRequest(req *UnubscribeEventMessageReq) *UnubscribeEventMessage {

	for ix, e := range i.SubscribedEntities {
		if req.MsgData.EntityIds == nil || slices.Contains(req.MsgData.EntityIds, e) {
			log.WithField("entity_id", e).Info("RT subscribed from entity")

			i.SubscribedEntities[ix] = i.SubscribedEntities[len(i.SubscribedEntities)-1] // Copy last element to index i.
			i.SubscribedEntities[len(i.SubscribedEntities)-1] = ""                       // Erase last element (write zero value).
			i.SubscribedEntities = i.SubscribedEntities[:len(i.SubscribedEntities)-1]    // Truncate slice.

			if entity, _, err := i.GetEntityById(e); err != nil {
				i.callUnubscribeCallback(entity)
			}
		}
	}

	log.WithField("subscribedEtities", i.SubscribedEntities).Debug("Change in subscribed entities")

	res := UnubscribeEventMessage{
		CommonResp{Kind: "resp", Id: req.Id, Msg: "result", Code: 200},
	}

	return &res
}

// Called by the Remote Two when it needs to synchronize the dynamic entity attributes, e.g. after connection setup or waking up from standby.
func (i *Integration) handleGetEntityStatesRequest(req *GetEntityStatesMessageReq) *GetEntityStatesMessage {

	var entityStates []entities.EntityStateData

	for _, e := range i.Entities {

		entity_id := i.getEntityId(e)
		device_id := i.getDeviceId(e)
		entity_type := i.getEntityType(e)
		attributes := i.getEntityAttributes(e)

		entity_state := entities.EntityStateData{
			EntityId:   entity_id,
			DeviceId:   device_id,
			EntityType: entity_type,
			Attributes: attributes,
		}
		entityStates = append(entityStates, entity_state)
	}

	res := GetEntityStatesMessage{
		CommonResp{Kind: "resp", Id: req.Id, Msg: "entity_states", Code: 200},
		entityStates,
	}

	return &res
}

// Handle the entity command request sent by the remote
func (i *Integration) handleEntityCommandRequest(req *EntityCommandReq) *EntityCommandResponse {

	log.WithFields(log.Fields{
		"entity_id": req.MsgData.EntityId,
		"command":   req.MsgData.CmdId,
		"params":    req.MsgData.Params}).Debug("Entity Command")

	entity, _, err := i.GetEntityById(req.MsgData.EntityId)

	var returnCode int

	if err == nil {
		returnCode = i.handleCommand(entity, req)
	} else {
		returnCode = 404
	}

	res := EntityCommandResponse{
		CommonResp{Kind: "resp", Id: req.Id, Msg: "result", Code: returnCode},
	}

	return &res

}

func (i *Integration) handleSetDriverUserDataRequest(req *SetDriverUserDataRequest) *ResponseMessage {

	log.WithField("Message", req.MsgData).Debug("Set DriverUserData Request")

	if i.handleSetDriverUserDataFunction != nil {
		go i.handleSetDriverUserDataFunction(req.MsgData.InputValues, req.MsgData.Confirm)
	}

	res := ResponseMessage{
		CommonResp{
			Kind: "resp",
			Id:   req.Id,
			Msg:  "result",
			Code: 200,
		},
		nil,
	}

	return &res
}
