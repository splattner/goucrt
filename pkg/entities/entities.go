package entities

import (
	log "github.com/sirupsen/logrus"
)

type EntityState string
type EntityCommand string
type EntityAttribute string
type EntityFeature string
type EntityOption string

const (
	UnavailableEntityState EntityState = "UNAVAILABLE"
	UnkownEntityState      EntityState = "UNKNOWN"
)

// EntityInfo is the read-only subset of Entity that BaseEntity implements directly, without needing
// a concrete entity type's own HandleCommand/UpdateEntity. It exists because of a Go embedding
// wrinkle: BaseEntity.SetAttributes fires its change callback as `e.handleEntityChangeFunc(e, ...)`,
// where `e` is BaseEntity's own receiver - a promoted method has no way to know or reach whatever
// concrete type (LightEntity, SwitchEntity, ...) embeds it. So that callback can only ever be
// handed something satisfying this narrower interface, never the full Entity below.
type EntityInfo interface {
	GetID() string
	GetDeviceID() string
	GetEntityType() EntityType
	GetAttribute() map[string]interface{}
}

// Entity is satisfied by every concrete entity type (ButtonEntity, LightEntity, SwitchEntity, ...).
// It's what lets package integration dispatch to a concrete entity's behavior - looking up its ID,
// forwarding a command, firing its subscribe callback - without a type switch listing every entity
// type by name. Every concrete type gets most of this for free by embedding BaseEntity; each only
// needs its own HandleCommand and UpdateEntity, since those genuinely differ per entity type.
type Entity interface {
	EntityInfo
	SetHandleEntityChangeFunc(f func(EntityInfo, *map[string]interface{}))
	CallSubscribeCallback()
	CallUnsubscribeCallback()
	// HandleCommand dispatches an entity_command's cmd_id/params to a registered command handler,
	// returning the response code to send back (200 on success, 404 if cmd_id isn't recognized).
	HandleCommand(cmdID string, params map[string]interface{}) int
	// UpdateEntity replaces this entity's mutable fields with those of newEntity, which must be the
	// same concrete type (e.g. a *ButtonEntity's UpdateEntity requires a ButtonEntity). Returns an
	// error if it isn't.
	UpdateEntity(newEntity interface{}) error
}

// BaseEntity is the common state and behavior every concrete entity type embeds. See the Entity
// interface for what package integration dispatches through; BaseEntity implements all of it except
// HandleCommand and UpdateEntity, which each concrete type provides itself.
//
// See https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/README.md for the wire format.
type BaseEntity struct {
	Id string `json:"entity_id"`
	EntityType
	DeviceId                string                                    `json:"device_id,omitempty"`
	Features                []interface{}                             `json:"features,omitempty"`
	Name                    LanguageText                              `json:"name"`
	Area                    string                                    `json:"area,omitempty"`
	DeviceClass             string                                    `json:"-"`
	Attributes              map[string]interface{}                    `json:"-"`
	handleEntityChangeFunc  func(EntityInfo, *map[string]interface{}) `json:"-"`
	SubscribeCallbackFunc   func()                                    `json:"-"`
	UnsubscribeCallbackFunc func()                                    `json:"-"`
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

func (e *BaseEntity) GetID() string {
	return e.Id
}

func (e *BaseEntity) GetDeviceID() string {
	return e.DeviceId
}

func (e *BaseEntity) GetEntityType() EntityType {
	return e.EntityType
}

func (e *BaseEntity) HasFeature(feature interface{}) bool {
	for _, f := range e.Features {
		if f == feature {
			return true
		}
	}
	return false
}

func (e *BaseEntity) GetAttribute() map[string]interface{} {
	return e.Attributes
}

// Add an attribute if not already available
func (e *BaseEntity) AddAttribute(name string, value interface{}) {

	if _, ok := e.Attributes[name]; !ok {
		log.WithFields(log.Fields{
			"entity_id": e.Id,
			"attribute": name,
		}).Debug("Add Attribute to entitiy")
		e.Attributes[name] = value
	}
}

// Retun the Entity State fr this entity
func (e *BaseEntity) GetEntityState() *EntityStateData {

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
func (e *BaseEntity) SetHandleEntityChangeFunc(f func(EntityInfo, *map[string]interface{})) {
	e.handleEntityChangeFunc = f
}

// Set the function that is called when RT subscribes to this entity
func (e *BaseEntity) SetSubscribeCallbackFunc(f func()) {
	e.SubscribeCallbackFunc = f
}

// Set the function that is called when RT unsubscribes to this entity
func (e *BaseEntity) SetUnsubscribeCallbackFunc(f func()) {
	e.UnsubscribeCallbackFunc = f
}

// CallSubscribeCallback calls SubscribeCallbackFunc if one is set. Safe to call unconditionally.
func (e *BaseEntity) CallSubscribeCallback() {
	if e.SubscribeCallbackFunc != nil {
		e.SubscribeCallbackFunc()
	}
}

// CallUnsubscribeCallback calls UnsubscribeCallbackFunc if one is set. Safe to call unconditionally.
func (e *BaseEntity) CallUnsubscribeCallback() {
	if e.UnsubscribeCallbackFunc != nil {
		e.UnsubscribeCallbackFunc()
	}
}

// Set one attribute for the Entity
func (e *BaseEntity) SetAttribute(attribute string, value interface{}) {

	attributes := make(map[string]interface{})
	attributes[attribute] = value
	e.SetAttributes(attributes)
}

// Set attributes for the Entity and then call the EntityChange Function
func (e *BaseEntity) SetAttributes(attributes map[string]interface{}) {

	log.WithFields(log.Fields{
		"entity_id":  e.Id,
		"attributes": attributes}).Info("Handle attribute change")

	for k, v := range attributes {
		e.Attributes[k] = v
	}

	// Handle the entity Change
	if e.handleEntityChangeFunc != nil {
		e.handleEntityChangeFunc(e, &attributes)
	}
}
