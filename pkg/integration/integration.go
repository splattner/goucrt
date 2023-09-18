package integration

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/grandcat/zeroconf"
	"github.com/splattner/goucrt/pkg/entities"
)

const API_VERSION = "0.8.1-alpha"

type Integration struct {
	DeviceId string

	Metadata *DriverMetadata

	authToken string

	deviceState DState

	config Config

	Remote remote

	Entities []interface{}

	SubscribedEntities []string

	handleSetupFunction             func(map[string]string)
	handleConnectionFunction        func(*ConnectEvent)
	handleSetDriverUserDataFunction func(map[string]string, bool)

	SetupState DriverSetupState

	mdns *zeroconf.Server
}

func NewIntegration(config Config) (*Integration, error) {

	i := Integration{
		config:      config,
		deviceState: DisconnectedDeviceState,
		DeviceId:    "", // I think device_id is not yet implemented in Remote TV, used for multi-device integrati

	}

	i.Remote.messageChannel = make(chan []byte)

	return &i, nil

}

func (i *Integration) SetMetadata(metadata *DriverMetadata) {
	log.WithField("Metadata", metadata).Debug("Set Metadata")
	i.Metadata = metadata
}

func (i *Integration) Run() error {

	if i.Metadata == nil {
		log.Error("Metadata not set, cannot run")
		return fmt.Errorf("Metadata not set")
	}

	http.HandleFunc("/ws", i.wsEndpoint)

	listenAddress := fmt.Sprintf(":%d", i.config["listenport"].(int))

	// MDNS
	i.startAdvertising()

	log.Fatal(http.ListenAndServe(listenAddress, nil))

	return nil

}

func (i *Integration) AddEntity(e interface{}) error {
	log.Debug("Add a new entity to the integration")

	// Search if entity is already added
	_, _, err := i.GetEntityById(i.getEntityId(e))
	if err != nil {
		// Entity not found, so add id
		i.Entities = append(i.Entities, e)
		// Send "entity_available" event to remote
		i.sendEntityAvailable(e)
		return nil
	}

	return fmt.Errorf("this entity is already added")
}

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

func (i *Integration) GetEntityById(id string) (interface{}, int, error) {
	for ix, entity := range i.Entities {
		entity_id := i.getEntityId(entity)
		log.Println(entity_id)

		if entity_id == id {
			log.Println("Found entity with type: " + fmt.Sprintf("%T", entity))
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

// Set the function which is called when the setup_driver request was sent by the remote
func (i *Integration) SetHandleSetupFunction(f func(map[string]string)) {
	i.handleSetupFunction = f
}

// Set the function which is called when the connect/disconnect request was sent by the remote
func (i *Integration) SetHandleConnectionFunction(f func(*ConnectEvent)) {
	i.handleConnectionFunction = f
}

// Set the function which is called when the connect/disconnect request was sent by the remote
func (i *Integration) SetHandleSetDriverUserDataFunction(f func(map[string]string, bool)) {
	i.handleSetDriverUserDataFunction = f
}

func (i *Integration) SetDriverSetupState(event_Type DriverSetupEventType, state DriverSetupState, err DriverSetupError, requireUserAction *RequireUserAction) {

	log.WithFields(log.Fields{
		"EventType": event_Type,
		"State":     state,
		"Error":     err,
	}).Info("Set DriverSetup State from Client")

	// Overwrite state if requireUserAction is set
	if requireUserAction != nil {
		state = WaitUserActionState
	}

	i.sendDriverSetupChangeEvent(event_Type, state, err, requireUserAction)

}
