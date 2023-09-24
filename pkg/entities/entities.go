package entities

import (
	log "github.com/sirupsen/logrus"
)

type EntityState string

const (
	UnavailableEntityState EntityState = "UNAVAILABLE"
	UnkownEntityState                  = "UNKNOWN"
)

type Entity struct {
	Id string `json:"entity_id"`
	EntityType
	DeviceId               string                 `json:"device_id,omitempty"`
	Features               []interface{}          `json:"features"`
	Name                   LanguageText           `json:"name"`
	Area                   string                 `json:"area,omitempty"`
	DeviceClass            string                 `json:"-"`
	Attributes             map[string]interface{} `json:"-"`
	handleEntityChangeFunc func(interface{})      `json:"-"`
}

type EntityType struct {
	Type string `json:"entity_type,omitempty"`
}

type EntityStateData struct {
	DeviceId string `json:"device_id,omitempty"`
	EntityType
	EntityId   string                 `json:"entity_id"`
	Attributes map[string]interface{} `json:"attributes"`
}

// Add an attribute if not already available
func (e *Entity) AddAttribute(name string, value interface{}) {

	if _, ok := e.Attributes[name]; !ok {
		log.WithFields(log.Fields{
			"entity_id": e.Id,
			"attribute": name,
		}).Debug("Add Attribute to entitiy")
		e.Attributes[name] = value
	}
}

// Retun the Entity State fr this entity
func (e *Entity) GetEntityState() *EntityStateData {

	entityState := EntityStateData{
		DeviceId:   e.DeviceId,
		EntityType: e.EntityType,
		EntityId:   e.Id,
		Attributes: e.Attributes,
	}

	return &entityState
}

// Register the function that is called when a Attribute change
// This normally is set by the integration when the entity is added
// To send entity_change events to Remote two
func (e *Entity) SetHandleEntityChangeFunc(f func(interface{})) {
	e.handleEntityChangeFunc = f
}

// Set attributes for the Entity and then call the EntityChange Function
func (e *Entity) SetAttributes(attributes map[string]interface{}) {

	log.WithFields(log.Fields{
		"entity_id":  e.Id,
		"attributes": attributes}).Info("Handle attribute change")

	for k, v := range attributes {
		e.Attributes[k] = v
	}

	// Handle the entity Change
	if e.handleEntityChangeFunc != nil {
		e.handleEntityChangeFunc(e)
	}

}
