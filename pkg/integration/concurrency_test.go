package integration

// Exercises the entity registry (Entities, SubscribedEntities) concurrently, the way it's actually
// used: a driver's own goroutines (device discovery, MQTT/WebSocket callbacks) call AddEntity,
// RemoveEntity and SendEntityChangeEvent while the WebSocket read loop concurrently calls
// handleRequest for subscribe_events/unsubscribe_events/get_available_entities. Run with `-race`
// (as `make test` does), this is the regression test for entitiesMu: remove the locking and this
// test flags a data race, whether or not it happens to also produce a visibly wrong result on any
// particular run.

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/splattner/goucrt/pkg/entities"
)

func TestConcurrentEntityRegistryAccess(t *testing.T) {
	i := newTestIntegration(t)

	// A single long-lived drain, mirroring the real wsWriter loop: exactly one goroutine ever
	// reads i.Remote.messageChannel. Every handleRequest call in this test relies on *something*
	// reading its response off that unbuffered channel, but each reader spinning up its own
	// one-shot receiver (as elsewhere in this package's tests) would race for messages against
	// every other goroutine's receiver here and can steal a response meant for a different call.
	drainDone := make(chan struct{})
	go func() {
		for {
			select {
			case <-i.Remote.messageChannel:
			case <-drainDone:
				return
			}
		}
	}()
	defer close(drainDone)

	const workers = 20
	var wg sync.WaitGroup

	// Writers: add and immediately remove a uniquely-IDed entity each, racing against each other
	// and against the readers below.
	for n := 0; n < workers; n++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			id := fmt.Sprintf("switch-%d", n)
			e := entities.NewSwitchEntity(id, entities.LanguageText{En: id}, "")
			if err := i.AddEntity(e); err != nil {
				t.Errorf("AddEntity(%s): %v", id, err)
			}
			e.SetAttributes(map[string]interface{}{"state": "ON"}) // exercises SendEntityChangeEvent
			if err := i.RemoveEntity(e); err != nil {
				t.Errorf("RemoveEntity(%s): %v", id, err)
			}
		}(n)
	}

	// Readers: hammer every entity-registry read path while the writers above are active.
	for n := 0; n < workers; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = i.GetEntitiesByType(entities.EntityType{Type: "switch"})
			_, _, _ = i.GetEntityById("switch-0")

			availableRaw, _ := json.Marshal(AvailableEntityMessageReq{
				CommonReq: CommonReq{Kind: "req", Id: 1, Msg: "get_available_entities"},
			})
			i.handleRequest(&RequestMessage{CommonReq: CommonReq{Kind: "req", Id: 1, Msg: "get_available_entities"}}, availableRaw)

			subRaw, _ := json.Marshal(SubscribeEventMessageReq{
				CommonReq: CommonReq{Kind: "req", Id: 2, Msg: "subscribe_events"},
			})
			i.handleRequest(&RequestMessage{CommonReq: CommonReq{Kind: "req", Id: 2, Msg: "subscribe_events"}}, subRaw)

			unsubRaw, _ := json.Marshal(UnubscribeEventMessageReq{
				CommonReq: CommonReq{Kind: "req", Id: 3, Msg: "unsubscribe_events"},
			})
			i.handleRequest(&RequestMessage{CommonReq: CommonReq{Kind: "req", Id: 3, Msg: "unsubscribe_events"}}, unsubRaw)
		}()
	}

	wg.Wait()

	if got := len(i.Entities); got != 0 {
		t.Errorf("len(Entities) = %d after every writer removed its own entity, want 0", got)
	}
}
