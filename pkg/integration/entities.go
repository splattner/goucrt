package integration

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/splattner/goucrt/pkg/entities"
	"k8s.io/utils/strings/slices"
)

// Return the ID of an entity
func (i *Integration) getEntityId(entity interface{}) string {
	var id string

	// Ugly.. I guess but I don't know how better
	switch e := entity.(type) {
	case *entities.Entity:
		id = e.Id
	case *entities.ButtonEntity:
		id = e.Id

	case *entities.LightEntity:
		id = e.Id

	case *entities.SwitchsEntity:
		id = e.Id

	case *entities.MediaPlayerEntity:
		id = e.Id

	case *entities.SensorEntity:
		id = e.Id

	case *entities.ClimateEntity:
		id = e.Id

	case *entities.CoverEntity:
		id = e.Id
	}

	return id
}

// Return the DeviceId of an entity
func (i *Integration) getDeviceId(entity interface{}) string {
	var device_id string

	// Ugly.. I guess but I don't know how better
	switch e := entity.(type) {
	case *entities.Entity:
		device_id = e.DeviceId
	case *entities.ButtonEntity:
		device_id = e.DeviceId

	case *entities.LightEntity:
		device_id = e.DeviceId

	case *entities.SwitchsEntity:
		device_id = e.DeviceId

	case *entities.MediaPlayerEntity:
		device_id = e.DeviceId

	case *entities.SensorEntity:
		device_id = e.DeviceId

	case *entities.ClimateEntity:
		device_id = e.DeviceId

	case *entities.CoverEntity:
		device_id = e.DeviceId
	}

	return device_id
}

// Return the EntityType of an entity
func (i *Integration) getEntityType(entity interface{}) entities.EntityType {
	var entity_type entities.EntityType

	// Ugly.. I guess but I don't know how better
	switch e := entity.(type) {
	case *entities.Entity:
		entity_type = e.EntityType
	case *entities.ButtonEntity:
		entity_type = e.EntityType

	case *entities.LightEntity:
		entity_type = e.EntityType

	case *entities.SwitchsEntity:
		entity_type = e.EntityType

	case *entities.MediaPlayerEntity:
		entity_type = e.EntityType

	case *entities.SensorEntity:
		entity_type = e.EntityType

	case *entities.ClimateEntity:
		entity_type = e.EntityType

	case *entities.CoverEntity:
		entity_type = e.EntityType
	}

	return entity_type
}

// Return the EntityType of an entity
func (i *Integration) getEntityAttributes(entity interface{}) map[string]interface{} {
	var attributes map[string]interface{}

	// Ugly.. I guess but I don't know how better
	switch e := entity.(type) {
	case *entities.Entity:
		attributes = e.GetAttribute()
	case *entities.ButtonEntity:
		attributes = e.GetAttribute()
	case *entities.LightEntity:
		attributes = e.GetAttribute()
	case *entities.SwitchsEntity:
		attributes = e.GetAttribute()
	case *entities.MediaPlayerEntity:
		attributes = e.GetAttribute()
	case *entities.SensorEntity:
		attributes = e.GetAttribute()
	case *entities.ClimateEntity:
		attributes = e.GetAttribute()
	case *entities.CoverEntity:
		attributes = e.GetAttribute()
	}

	return attributes
}

// Add a new Entity to the list of Entities (if not already added)
// Also make sure the EntityChange Function is set so Entity Change Events are emitted when a Entity Attribute changes
// Send Entity Available Event to RT
func (i *Integration) AddEntity(e interface{}) error {
	entity_id := i.getEntityId(e)
	log.WithField("entity_id", entity_id).Debug("Add a new entity to the integration")

	// Search if entity is already added
	existingEntity, _, err := i.GetEntityById(entity_id)
	if err != nil {
		// Entity not found, so add id
		i.setEntityChangeFunc(e, i.SendEntityChangeEvent)
		i.Entities = append(i.Entities, e)
		// Send "entity_available" event to remote
		i.sendEntityAvailable(e)

		// if RT already subscribed, call the Subscribe callback for this entity
		if i.isSubscribed(e) {
			i.callSubscribeCallback(e)
		}

		return nil
	}

	// else update the existing entity
	return i.UpdateEntity(existingEntity, e)
}

func (i *Integration) isSubscribed(entity interface{}) bool {
	entity_id := i.getEntityId(entity)

	return slices.Contains(i.SubscribedEntities, entity_id)
}

// Update an existing entity with a new entity
func (i *Integration) UpdateEntity(entity interface{}, newEntity interface{}) error {
	switch e := entity.(type) {
	case *entities.ButtonEntity:
		return e.UpdateEntity(newEntity.(entities.ButtonEntity))

	case *entities.LightEntity:
		return e.UpdateEntity(newEntity.(entities.LightEntity))

	case *entities.SwitchsEntity:
		return e.UpdateEntity(newEntity.(entities.SwitchsEntity))

	case *entities.MediaPlayerEntity:
		return e.UpdateEntity(newEntity.(entities.MediaPlayerEntity))

	case *entities.SensorEntity:
		return e.UpdateEntity(newEntity.(entities.SensorEntity))

	case *entities.ClimateEntity:
		return e.UpdateEntity(newEntity.(entities.ClimateEntity))

	case *entities.CoverEntity:
		return e.UpdateEntity(newEntity.(entities.CoverEntity))
	}

	return nil
}

// Remove an Entity from the Integration
// Send Entity Removed Event to RT
func (i *Integration) RemoveEntity(entity interface{}) error {
	// Search if entity is available

	return i.RemoveEntityByID(i.getEntityId(entity))

}

// Remove an Entity from the Integration
// Send Entity Removed Event to RT
func (i *Integration) RemoveEntityByID(entity_id string) error {
	// Search if entity is available

	entity, ix, err := i.GetEntityById(entity_id)
	if err == nil {

		i.Entities[ix] = i.Entities[len(i.Entities)-1] // Copy last element to index i.
		i.Entities[len(i.Entities)-1] = ""             // Erase last element (write zero value).
		i.Entities = i.Entities[:len(i.Entities)-1]    // Truncate slice.

		// Send "entity_removed" event to remote
		i.sendEntityRemoved(entity)
		return nil
	}

	return fmt.Errorf("entity to remove not found")
}

// Return an Entity by its Name
// Also return the current index in the Entities Array (TODO: do we need this?)
// Error when Entity not found
func (i *Integration) GetEntityById(id string) (interface{}, int, error) {
	for ix, entity := range i.Entities {
		entity_id := i.getEntityId(entity)

		if entity_id == id {
			return entity, ix, nil
		}
	}

	return entities.Entity{}, 0, fmt.Errorf("entity with id %s not found", id)
}

// Return all available entities of a given type
func (i *Integration) GetEntitiesByType(entityType entities.EntityType) []interface{} {
	var es []interface{}

	for _, e := range i.Entities {
		if i.getEntityType(e) == entityType {
			es = append(es, e)
		}
	}

	return es
}

// Call the correct HandleCommand function depending on the entity type
func (i *Integration) handleCommand(entity interface{}, req *EntityCommandReq) int {
	cmd_id := req.MsgData.CmdId
	params := req.MsgData.Params

	// Ugly.. I guess but I don't know how better
	switch e := entity.(type) {
	case *entities.ButtonEntity:
		return e.HandleCommand(cmd_id)

	case *entities.LightEntity:
		return e.HandleCommand(cmd_id, params)

	case *entities.SwitchsEntity:
		return e.HandleCommand(cmd_id, params)

	case *entities.MediaPlayerEntity:
		return e.HandleCommand(cmd_id, params)

	case *entities.ClimateEntity:
		return e.HandleCommand(cmd_id, params)

	case *entities.CoverEntity:
		return e.HandleCommand(cmd_id, params)

	case *entities.SensorEntity:
		// Sensor do not have commands
	}

	return 404
}

// Call the correct HandleCommand function depending on the entity type
func (i *Integration) setEntityChangeFunc(entity interface{}, f func(interface{}, *map[string]interface{})) {

	// Ugly.. I guess but I don't know how better
	switch e := entity.(type) {
	case *entities.Entity:
		e.SetHandleEntityChangeFunc(f)
	case *entities.SensorEntity:
		e.SetHandleEntityChangeFunc(f)
	case *entities.ButtonEntity:
		e.SetHandleEntityChangeFunc(f)
	case *entities.LightEntity:
		e.SetHandleEntityChangeFunc(f)
	case *entities.SwitchsEntity:
		e.SetHandleEntityChangeFunc(f)
	case *entities.MediaPlayerEntity:
		e.SetHandleEntityChangeFunc(f)
	case *entities.ClimateEntity:
		e.SetHandleEntityChangeFunc(f)
	case *entities.CoverEntity:
		e.SetHandleEntityChangeFunc(f)
	}
}

// Call the correct subscribe Callback Func depending on the entity type
func (i *Integration) callSubscribeCallback(entity interface{}) {

	// Ugly.. I guess but I don't know how better
	switch e := entity.(type) {
	case *entities.ButtonEntity:
		if e.SubscribeCallbackFunc != nil {
			e.SubscribeCallbackFunc()
		}

	case *entities.LightEntity:
		if e.SubscribeCallbackFunc != nil {
			e.SubscribeCallbackFunc()
		}

	case *entities.SwitchsEntity:
		if e.SubscribeCallbackFunc != nil {
			e.SubscribeCallbackFunc()
		}

	case *entities.MediaPlayerEntity:
		if e.SubscribeCallbackFunc != nil {
			e.SubscribeCallbackFunc()
		}

	case *entities.ClimateEntity:
		if e.SubscribeCallbackFunc != nil {
			e.SubscribeCallbackFunc()
		}

	case *entities.CoverEntity:
		if e.SubscribeCallbackFunc != nil {
			e.SubscribeCallbackFunc()
		}

	case *entities.SensorEntity:
		if e.SubscribeCallbackFunc != nil {
			e.SubscribeCallbackFunc()
		}
	}
}

// Call the correct unsubscribe Callback Func depending on the entity type
func (i *Integration) callUnubscribeCallback(entity interface{}) {

	// Ugly.. I guess but I don't know how better
	switch e := entity.(type) {
	case *entities.ButtonEntity:
		if e.UnsubscribeCallbackFunc != nil {
			e.UnsubscribeCallbackFunc()
		}

	case *entities.LightEntity:
		if e.UnsubscribeCallbackFunc != nil {
			e.UnsubscribeCallbackFunc()
		}

	case *entities.SwitchsEntity:
		if e.UnsubscribeCallbackFunc != nil {
			e.UnsubscribeCallbackFunc()
		}

	case *entities.MediaPlayerEntity:
		if e.UnsubscribeCallbackFunc != nil {
			e.UnsubscribeCallbackFunc()
		}

	case *entities.ClimateEntity:
		if e.UnsubscribeCallbackFunc != nil {
			e.UnsubscribeCallbackFunc()
		}

	case *entities.CoverEntity:
		if e.SubscribeCallbackFunc != nil {
			e.SubscribeCallbackFunc()
		}

	case *entities.SensorEntity:
		if e.UnsubscribeCallbackFunc != nil {
			e.UnsubscribeCallbackFunc()
		}
	}
}
