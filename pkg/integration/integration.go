package integration

import (
	"fmt"
	"log"
	"net/http"

	"github.com/grandcat/zeroconf"
	"github.com/splattner/goucrt/pkg/entities"
)

type Integration struct {
	DeviceId string

	Metadata DriverMetadata

	authToken string

	deviceState DState

	config Config

	Remote remote

	Entities []interface{}

	// User input result of a SettingsPage as key values.
	// key: id of the field
	// value: entered user value as string. This is either the entered text or number, selected checkbox state or the selected dropdown item id.
	//⚠️ Non native string values as numbers or booleans are represented as string values!
	UserInputValues       map[string]interface{}
	UserInputConfirmation bool

	SubscribedEntities []string

	handleSetupFunction      func()
	handleConnectionFunction func(*ConnectEvent)

	SetupState DriverSetupState

	mdns *zeroconf.Server
}

func NewIntegration(config Config) (*Integration, error) {

	metadata := DriverMetadata{
		DriverId: "myintegration",
		Developer: Developer{
			Name: "Sebastian Plattner",
		},
		Name: LanguageText{
			En: "My UCRT Integration",
			De: "Meine UCRT Integration",
		},
		Version: "0.0.1",
	}

	i := Integration{
		config:      config,
		Metadata:    metadata,
		deviceState: DisconnectedDeviceState,
		DeviceId:    "", // I think device_id is not yet implemented in Remote TV, used for multi-device integration

	}

	return &i, nil

}

func (i *Integration) Run() error {

	http.HandleFunc("/ws", i.wsEndpoint)

	listenAddress := fmt.Sprintf(":%d", i.config["listenport"].(int))

	// MDNS
	i.startAdvertising()

	log.Fatal(http.ListenAndServe(listenAddress, nil))

	return nil

}

func (i *Integration) AddEntity(e interface{}) error {
	log.Println("Add a new entity to the integration")

	// Search if entity is already added
	_, _, err := i.GetEntityById(GetEntityId(e))
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

	_, ix, err := i.GetEntityById(GetEntityId(entity))
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
	for i, entity := range i.Entities {
		entity_id := GetEntityId(entity)
		log.Println(entity_id)

		if entity_id == id {
			log.Println("Found entity with type: " + fmt.Sprintf("%T", entity))
			return entity, i, nil
		}
	}

	return entities.Entity{}, 0, fmt.Errorf("entity with id %s not found", id)
}

// Return all available entities of a given type
func (i *Integration) GetEntitiesByType(entityType entities.EntityType) []interface{} {
	var es []interface{}

	for _, e := range i.Entities {
		if GetEntityType(e) == entityType {
			es = append(es, e)
		}
	}

	return es
}

// Set the function which is called when the setup_driver request was sent by the remote
func (i *Integration) SetHandleSetupFunction(f func()) {
	i.handleSetupFunction = f
}

// Set the function which is called when the connect/disconnect request was sent by the remote
func (i *Integration) SetHandleConnectionFunction(f func(*ConnectEvent)) {
	i.handleConnectionFunction = f
}

func (i *Integration) SetDriverSetupState(event_Type DriverSetupEventType, state DriverSetupState, err DriverSetupError, requiredUserAction *RequiredUserAction) {

	i.sendDriverSetupChangeEvent(event_Type, state, err, requiredUserAction)

}
