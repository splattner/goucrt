package integration

// Driver-initiated metadata requests: get_version, get_supported_entity_types,
// get_configured_entities, get_localization_cfg and get_runtime_info. Per the Core-API spec these
// are the one family of messages the driver sends TO the Remote rather than the other way around -
// the driver asking about the Remote itself (its version, locale, which of this driver's entities
// are actually configured, ...), answered with a correlated "resp" frame carrying the same req_id.
//
// Every other request/response pair in this package is the Remote calling into the driver
// (handleRequest, response.go), so there was previously no code path that ever expected a "resp"
// kind message at all - wsReader dispatched "event" and "req" but silently dropped "resp". This
// file adds that: a pending-request table keyed by request id, filled in here before sending and
// drained by handleResponse when the matching frame arrives.

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// requestTimeout bounds how long a driver-initiated request (GetVersion, GetRuntimeInfo, ...)
// waits - both to hand its request frame to the WebSocket write loop and to receive the
// correlated response - before giving up. Without this, a driver goroutine calling e.g. GetVersion
// while no Remote is connected (or after it disconnected mid-request) would block forever: nothing
// else in this package ever times out a send on the unbuffered Remote.messageChannel.
//
// A var rather than a const so tests can shorten it instead of spending 10 real seconds proving a
// timeout path works.
var requestTimeout = 10 * time.Second

// nextRequestID returns the next outbound request id, unique among the driver's own in-flight
// requests. The Remote echoes it back as req_id so the response can be routed to the right waiter.
func (i *Integration) nextRequestID() int {
	i.requestsMu.Lock()
	defer i.requestsMu.Unlock()
	i.nextReqID++
	return i.nextReqID
}

func (i *Integration) registerPendingRequest(id int) chan []byte {
	ch := make(chan []byte, 1)
	i.requestsMu.Lock()
	i.pendingRequests[id] = ch
	i.requestsMu.Unlock()
	return ch
}

func (i *Integration) unregisterPendingRequest(id int) {
	i.requestsMu.Lock()
	delete(i.pendingRequests, id)
	i.requestsMu.Unlock()
}

// handleResponse routes a received "resp" kind message to whichever sendMetadataRequest call is
// waiting for it, matched by req_id. A response with no matching pending request (already timed
// out, or unsolicited) is logged and dropped - there's nothing to deliver it to.
func (i *Integration) handleResponse(p []byte) {
	var resp ResponseMessage
	if err := json.Unmarshal(p, &resp); err != nil {
		log.WithError(err).Error("Cannot unmarshal response message")
		return
	}

	i.requestsMu.Lock()
	ch, ok := i.pendingRequests[resp.Id]
	i.requestsMu.Unlock()

	if !ok {
		log.WithField("req_id", resp.Id).Debug("Received response for unknown or already-completed request")
		return
	}

	// Buffered size 1 and written to at most once (handleResponse only runs for req_ids still in
	// pendingRequests, and the request's own req_id is unregistered as soon as it's delivered or
	// times out), so this never blocks.
	ch <- p
}

// sendMetadataRequest sends a driver-initiated request carrying no msg_data (get_version and the
// other messages in this file all take none) and waits for its correlated response, returning the
// full raw response frame for the caller to unmarshal into the concrete *Message type for msg.
func (i *Integration) sendMetadataRequest(msg string) ([]byte, error) {
	id := i.nextRequestID()
	ch := i.registerPendingRequest(id)
	defer i.unregisterPendingRequest(id)

	data, err := json.Marshal(RequestMessage{CommonReq: CommonReq{Kind: "req", Id: id, Msg: msg}})
	if err != nil {
		return nil, fmt.Errorf("cannot marshal %q request: %w", msg, err)
	}

	select {
	case i.Remote.messageChannel <- data:
	case <-time.After(requestTimeout):
		return nil, fmt.Errorf("timed out sending %q request: no Remote connected", msg)
	}

	select {
	case raw := <-ch:
		return raw, nil
	case <-time.After(requestTimeout):
		return nil, fmt.Errorf("timed out waiting for %q response", msg)
	}
}

// GetVersion requests the Remote's version information (implemented in firmware 0.9.2+).
func (i *Integration) GetVersion() (*VersionInfo, error) {
	raw, err := i.sendMetadataRequest("get_version")
	if err != nil {
		return nil, err
	}

	var resp VersionMessage
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("cannot unmarshal version response: %w", err)
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get_version request failed with code %d", resp.Code)
	}

	return &resp.MsgData, nil
}

// GetSupportedEntityTypes requests the entity types the Remote supports (implemented in firmware
// 0.9.2+), letting a driver check compatibility since new releases can support new entity types or
// rename existing ones in major updates.
func (i *Integration) GetSupportedEntityTypes() ([]string, error) {
	raw, err := i.sendMetadataRequest("get_supported_entity_types")
	if err != nil {
		return nil, err
	}

	var resp SupportedEntityTypesMessage
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("cannot unmarshal supported_entity_types response: %w", err)
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get_supported_entity_types request failed with code %d", resp.Code)
	}

	return resp.MsgData, nil
}

// GetConfiguredEntities requests the IDs of this driver's entities that are actually configured in
// the Remote (implemented in firmware 0.9.2+) - not necessarily all of them assigned to a profile
// and shown in the UI, just configured at all. Lets a driver limit its own device communication to
// entities the Remote actually cares about, as an alternative or complement to subscribe_events.
func (i *Integration) GetConfiguredEntities() ([]string, error) {
	raw, err := i.sendMetadataRequest("get_configured_entities")
	if err != nil {
		return nil, err
	}

	var resp ConfiguredEntitiesMessage
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("cannot unmarshal configured_entities response: %w", err)
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get_configured_entities request failed with code %d", resp.Code)
	}

	return resp.MsgData, nil
}

// GetLocalizationCfg requests the Remote's active localization settings (implemented in firmware
// 0.9.2+), for a driver that needs localized text or units of measurement. Per the spec, a driver
// should always provide an English fallback in addition to any localized text.
func (i *Integration) GetLocalizationCfg() (*LocalizationSettings, error) {
	raw, err := i.sendMetadataRequest("get_localization_cfg")
	if err != nil {
		return nil, err
	}

	var resp LocalizationCfgMessage
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("cannot unmarshal localization_cfg response: %w", err)
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get_localization_cfg request failed with code %d", resp.Code)
	}

	return &resp.MsgData, nil
}

// GetRuntimeInfo requests this driver's runtime information from the Remote (implemented in
// firmware 0.9.2+): driver and integration instance identifiers for advanced use of the REST
// Core-API, e.g. to propagate to a home automation system that wants to interact with the Remote
// directly.
func (i *Integration) GetRuntimeInfo() (*RuntimeInfo, error) {
	raw, err := i.sendMetadataRequest("get_runtime_info")
	if err != nil {
		return nil, err
	}

	var resp RuntimeInfoMessage
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("cannot unmarshal runtime_info response: %w", err)
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get_runtime_info request failed with code %d", resp.Code)
	}

	return &resp.MsgData, nil
}
