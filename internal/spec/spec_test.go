package spec

import (
	"slices"
	"testing"
)

func TestEntityTypeSchema_KnownEntities(t *testing.T) {
	// Spot-checks against the values read directly out of the vendored YAML during development
	// (see internal/spec/core-api/VENDORED_COMMIT.txt for the pinned commit). If these fail after
	// updating the vendored spec, the spec's shape has changed and EntityTypeSchema's YAML-walking
	// logic needs to be revisited, not just the expected values here.
	cases := []struct {
		entityType      string
		wantFeature     string // one feature that must be present
		wantDeviceClass string // one device_class that must be present, "" if the type has none
	}{
		{"button", "press", ""},
		{"switch", "on_off", "outlet"},
		{"light", "color_temperature", ""},
		{"cover", "tilt_position", "blind"},
		{"media_player", "play_pause", ""},
		{"climate", "target_temperature", ""},
		{"sensor", "", "temperature"},
		{"remote", "send_cmd", ""},
	}

	for _, c := range cases {
		t.Run(c.entityType, func(t *testing.T) {
			got, err := EntityTypeSchema(c.entityType)
			if err != nil {
				t.Fatalf("EntityTypeSchema(%q): %v", c.entityType, err)
			}
			if c.wantFeature != "" && !slices.Contains(got.Features, c.wantFeature) {
				t.Errorf("Features = %v, want to contain %q", got.Features, c.wantFeature)
			}
			if c.wantDeviceClass != "" && !slices.Contains(got.DeviceClass, c.wantDeviceClass) {
				t.Errorf("DeviceClass = %v, want to contain %q", got.DeviceClass, c.wantDeviceClass)
			}
		})
	}
}

func TestEntityTypeSchema_UnknownEntity(t *testing.T) {
	if _, err := EntityTypeSchema("not_a_real_entity_type"); err == nil {
		t.Fatal("expected an error for an unknown entity_type, got nil")
	}
}

func TestEntityTypeSchema_SensorHasNoFeatures(t *testing.T) {
	got, err := EntityTypeSchema("sensor")
	if err != nil {
		t.Fatalf("EntityTypeSchema(\"sensor\"): %v", err)
	}
	if len(got.Features) != 0 {
		t.Errorf("Features = %v, want empty (the sensor entity has no features per the spec)", got.Features)
	}
}
