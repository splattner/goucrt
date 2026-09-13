package integration

// Round-trip-decodes a representative sample of the worked examples extracted from the vendored
// Core-API entity docs (see tools/extractgoldens, fixtures under pkg/entities/testdata) through
// goucrt's actual wire types - EntityCommandReq and EntityChangeEvent - rather than a generic
// struct. pkg/entities' spec_golden_test.go checks the same fixtures more broadly (cmd_id/attribute
// coverage across every extracted example); this test's job is narrower and complementary: prove
// that the literal types handleRequest/SendEntityChangeEvent use in production actually decode real
// spec examples without error, so a struct tag regression (the class of bug the missing
// `device_class` tag was, fixed in an earlier pass) would fail a test instead of shipping silently.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGoldenExamplesDecodeAsWireTypes(t *testing.T) {
	t.Run("entity_command", func(t *testing.T) {
		// One representative command fixture per implemented entity type that has commands.
		cases := []struct {
			entityType string
			file       string
			wantCmdID  string
		}{
			{"switch", "on.json", "on"},
			{"light", "on.json", "on"},
			{"cover", "open.json", "open"},
			{"media_player", "on.json", "on"},
			{"climate", "hvac_mode.json", "hvac_mode"},
			{"remote", "on.json", "on"},
		}

		for _, c := range cases {
			t.Run(c.entityType, func(t *testing.T) {
				raw := readGoldenFixture(t, c.entityType, "commands", c.file)

				var req RequestMessage
				if err := json.Unmarshal(raw, &req); err != nil {
					t.Fatalf("decode as RequestMessage: %v", err)
				}
				if req.Kind != "req" || req.Msg != "entity_command" {
					t.Fatalf("envelope = {kind: %q, msg: %q}, want {req, entity_command}", req.Kind, req.Msg)
				}

				var full EntityCommandReq
				if err := json.Unmarshal(raw, &full); err != nil {
					t.Fatalf("decode as EntityCommandReq: %v", err)
				}
				if full.MsgData.CmdId != c.wantCmdID {
					t.Errorf("MsgData.CmdId = %q, want %q", full.MsgData.CmdId, c.wantCmdID)
				}
				if full.MsgData.EntityId == "" {
					t.Error("MsgData.EntityId is empty")
				}
			})
		}
	})

	t.Run("entity_change", func(t *testing.T) {
		// One representative event fixture per implemented entity type.
		cases := []struct {
			entityType string
			file       string
		}{
			{"switch", "switched_on.json"},
			{"light", "switched_on.json"},
			{"cover", "state_change_event.json"},
			{"media_player", "state_change_event.json"},
			{"climate", "state_change_event.json"},
			{"sensor", "state_change_event.json"},
			{"remote", "device_was_turned_on.json"},
		}

		for _, c := range cases {
			t.Run(c.entityType, func(t *testing.T) {
				raw := readGoldenFixture(t, c.entityType, "events", c.file)

				var event EntityChangeEvent
				if err := json.Unmarshal(raw, &event); err != nil {
					t.Fatalf("decode as EntityChangeEvent: %v", err)
				}
				if event.Kind != "event" || event.Msg != "entity_change" {
					t.Fatalf("envelope = {kind: %q, msg: %q}, want {event, entity_change}", event.Kind, event.Msg)
				}
				if event.MsgData.EntityType != c.entityType {
					t.Errorf("MsgData.EntityType = %q, want %q", event.MsgData.EntityType, c.entityType)
				}
				if len(event.MsgData.Attributes) == 0 {
					t.Error("MsgData.Attributes is empty")
				}
			})
		}
	})
}

func readGoldenFixture(t *testing.T, entityType, sub, file string) []byte {
	t.Helper()
	path := filepath.Join("..", "entities", "testdata", entityType, sub, file)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s (run `go run ./tools/extractgoldens` from the repository root if this is missing): %v", path, err)
	}
	return raw
}
