package integration

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/splattner/goucrt/pkg/entities"
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
		attributes = e.Attributes
	case *entities.ButtonEntity:
		attributes = e.Attributes

	case *entities.LightEntity:
		attributes = e.Attributes

	case *entities.SwitchsEntity:
		attributes = e.Attributes

	case *entities.MediaPlayerEntity:
		attributes = e.Attributes

	case *entities.SensorEntity:
		attributes = e.Attributes

	case *entities.ClimateEntity:
		attributes = e.Attributes

	case *entities.CoverEntity:
		attributes = e.Attributes
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
	_, _, err := i.GetEntityById(entity_id)
	if err != nil {
		// Entity not found, so add id
		i.setEntityChangeFunc(e, i.SendEntityChangeEvent)
		i.Entities = append(i.Entities, e)
		// Send "entity_available" event to remote
		i.sendEntityAvailable(e)
		return nil
	}

	return fmt.Errorf("this entity is already added")
}

// Remove an Entity from the Integration
// Send Entity Removed Event to RT
func (i *Integration) RemoveEntity(entity interface{}) error {
	// Search if entity is available

	_, ix, err := i.GetEntityById(i.getEntityId(entity))
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
func (i *Integration) handleCommand(entity interface{}, req *EntityCommandReq) {
	cmd_id := req.MsgData.CmdId
	params := req.MsgData.Params

	// Ugly.. I guess but I don't know how better
	switch e := entity.(type) {
	case *entities.ButtonEntity:
		e.HandleCommand(cmd_id, params)

	case *entities.LightEntity:
		e.HandleCommand(cmd_id, params)

	case *entities.SwitchsEntity:
		e.HandleCommand(cmd_id, params)

	case *entities.MediaPlayerEntity:
		e.HandleCommand(cmd_id, params)

	case *entities.ClimateEntity:
		e.HandleCommand(cmd_id, params)

	case *entities.CoverEntity:
		e.HandleCommand(cmd_id, params)
	}
}

// Call the correct HandleCommand function depending on the entity type
func (i *Integration) setEntityChangeFunc(entity interface{}, f func(interface{})) {

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
