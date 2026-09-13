package integration

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/splattner/goucrt/pkg/entities"
	"k8s.io/utils/strings/slices"
)

// Add a new Entity to the list of Entities (if not already added)
// Also make sure the EntityChange Function is set so Entity Change Events are emitted when a Entity Attribute changes
// Send Entity Available Event to RT
func (i *Integration) AddEntity(e entities.Entity) error {
	entity_id := e.GetID()
	log.WithField("entity_id", entity_id).Debug("Add a new entity to the integration")

	i.entitiesMu.Lock()
	existing, _ := i.findEntityByIdLocked(entity_id)
	if existing == nil {
		e.SetHandleEntityChangeFunc(i.SendEntityChangeEvent)
		i.Entities = append(i.Entities, e)
	}
	i.entitiesMu.Unlock()

	if existing != nil {
		// Entity already registered: update it in place instead of adding a duplicate.
		return existing.UpdateEntity(e)
	}

	// Send "entity_available" event to remote
	i.sendEntityAvailable(e)

	// if RT already subscribed, call the Subscribe callback for this entity
	if i.isSubscribed(e) {
		e.CallSubscribeCallback()
	}

	return nil
}

func (i *Integration) isSubscribed(entity entities.Entity) bool {
	i.entitiesMu.RLock()
	defer i.entitiesMu.RUnlock()

	return slices.Contains(i.SubscribedEntities, entity.GetID())
}

// Remove an Entity from the Integration
// Send Entity Removed Event to RT
func (i *Integration) RemoveEntity(entity entities.Entity) error {
	return i.RemoveEntityByID(entity.GetID())
}

// Remove an Entity from the Integration
// Send Entity Removed Event to RT
func (i *Integration) RemoveEntityByID(entity_id string) error {
	i.entitiesMu.Lock()
	entity, ix := i.findEntityByIdLocked(entity_id)
	if entity == nil {
		i.entitiesMu.Unlock()
		return fmt.Errorf("entity to remove not found")
	}

	i.Entities[ix] = i.Entities[len(i.Entities)-1] // Copy last element to index i.
	i.Entities[len(i.Entities)-1] = nil            // Erase last element (write zero value).
	i.Entities = i.Entities[:len(i.Entities)-1]    // Truncate slice.
	i.entitiesMu.Unlock()

	entity.CallUnsubscribeCallback()

	// Send "entity_removed" event to remote
	i.sendEntityRemoved(entity)
	return nil
}

// Return an Entity by its Name
// Also return the current index in the Entities Array (TODO: do we need this?)
// Error when Entity not found
func (i *Integration) GetEntityById(id string) (entities.Entity, int, error) {
	i.entitiesMu.RLock()
	defer i.entitiesMu.RUnlock()

	entity, ix := i.findEntityByIdLocked(id)
	if entity == nil {
		return nil, 0, fmt.Errorf("entity with id %s not found", id)
	}
	return entity, ix, nil
}

// findEntityByIdLocked is GetEntityById's lock-free core, for callers that already hold
// entitiesMu (in either mode - it only reads i.Entities). Returns (nil, 0) on a miss.
func (i *Integration) findEntityByIdLocked(id string) (entities.Entity, int) {
	for ix, entity := range i.Entities {
		if entity.GetID() == id {
			return entity, ix
		}
	}
	return nil, 0
}

// Return all available entities of a given type
func (i *Integration) GetEntitiesByType(entityType entities.EntityType) []entities.Entity {
	i.entitiesMu.RLock()
	defer i.entitiesMu.RUnlock()

	var es []entities.Entity
	for _, e := range i.Entities {
		if e.GetEntityType() == entityType {
			es = append(es, e)
		}
	}

	return es
}

// handleCommand dispatches an entity_command request to the entity's own HandleCommand.
func (i *Integration) handleCommand(entity entities.Entity, req *EntityCommandReq) int {
	return entity.HandleCommand(req.MsgData.CmdId, req.MsgData.Params)
}
