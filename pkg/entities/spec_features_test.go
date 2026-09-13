package entities

// Checks goucrt's hand-maintained `features` and `device_class` constants against the vendored
// Core-API spec (see internal/spec). This is exactly the class of bug Phase 1 fixed by hand
// (a mistyped or stale enum value the remote silently ignores): a failure here names the offending
// value directly instead of waiting for someone to notice by eye.
//
// Every value goucrt declares must exist in the spec's enum for that entity type. The reverse isn't
// required - the spec is free to define more than goucrt currently implements (see the t.Log output
// for known gaps; they're tracked as follow-up work, not failures here).

import (
	"testing"

	"github.com/splattner/goucrt/internal/spec"
	"github.com/splattner/goucrt/internal/spec/spectest"
)

func TestFeaturesAndDeviceClassMatchSpec(t *testing.T) {
	cases := []struct {
		entityType  string
		features    []string
		deviceClass []string
	}{
		{
			entityType: "button",
			features:   []string{string(PressButtonEntityFeatures)},
		},
		{
			entityType: "switch",
			features: []string{
				string(OnOffSwitchEntityyFeatures),
				string(ToggleSwitchEntityyFeatures),
			},
		},
		{
			entityType: "light",
			features: []string{
				string(OnOffLightEntityFeatures),
				string(ToggleLightEntityFeatures),
				string(DimLightEntityFeatures),
				string(ColorLightEntityFeatures),
				string(ColorTemperatureLightEntityFeatures),
			},
		},
		{
			entityType: "cover",
			features: []string{
				string(OpenCoverEntityFeatures),
				string(CloseCoverEntityFeatures),
				string(StopCoverEntityFeatures),
				string(PositionCoverEntityFeatures),
				string(TiltCoverEntityFeatures),
				string(TiltStopCoverEntityFeatures),
				string(TiltPositionCoverEntityFeatures),
			},
		},
		{
			entityType: "media_player",
			features: []string{
				string(OnOffMediaPlayerEntityFeatures),
				string(ToggleMediaPlayerEntityyFeatures),
				string(VolumeMediaPlayerEntityyFeatures),
				string(VolumeUpDownMediaPlayerEntityFeatures),
				string(MuteToggleMediaPlayerEntityFeatures),
				string(MuteMediaPlayerEntityFeatures),
				string(UnmuteMediaPlayerEntityFeatures),
				string(PlayPauseMediaPlayerEntityFeatures),
				string(StopMediaPlayerEntityFeatures),
				string(NextMediaPlayerEntityFeatures),
				string(PreviusMediaPlayerEntityFeatures),
				string(FastForwardMediaPlayerEntityFeatures),
				string(RewindMediaPlayerEntityFeatures),
				string(RepeatMediaPlayerEntityFeatures),
				string(ShuffleMediaPlayerEntityFeatures),
				string(SeekMediaPlayerEntityFeatures),
				string(MediaDurationMediaPlayerEntityFeatures),
				string(MediaPositionMediaPlayerEntityFeatures),
				string(MediaPositionUpdatedAtMediaPlayerEntityFeatures),
				string(MediaTitleMediaPlayerEntityFeatures),
				string(MediaArtistMediaPlayerEntityFeatures),
				string(MediaAlbumMediaPlayerEntityFeatures),
				string(MediaImageUrlMediaPlayerEntityFeatures),
				string(MediaTypeMediaPlayerEntityFeatures),
				string(DPadMediaPlayerEntityFeatures),
				string(NumPadMediaPlayerEntityFeatures),
				string(HomeMediaPlayerEntityFeatures),
				string(MenuMediaPlayerEntityFeatures),
				string(ContextMenuPlayerEntityFeatures),
				string(GuidePlayerEntityFeatures),
				string(InfoPlayerEntityFeatures),
				string(ColorButtonsMediaPlayerEntityFeatures),
				string(ChannelSwitcherMediaPlayerEntityFeatures),
				string(SelectSourceMediaPlayerEntityFeatures),
				string(SelectSoundModeMediaPlayerEntityFeatures),
				string(EjectMediaPlayerEntityFeatures),
				string(OpenCloseMediaPlayerEntityFeatures),
				string(AudioTrackMediaPlayerEntityFeatures),
				string(SubtitleMediaPlayerEntityFeatures),
				string(RecordMediaPlayerEntityFeatures),
				string(SettingsMediaPlayerEntityFeatures),
			},
			// media_player's device_class is deliberately not checked here: the vendored YAML's
			// media_player schema has no device_class property at all (confirmed against the
			// spec directly, not just a parsing gap), while doc/entities/entity_media_player.md's
			// "Device Classes" section documents exactly these five values. The machine-readable
			// schema lags its own prose docs here - goucrt's constants are correct, there's just
			// nothing in the YAML to check them against yet.
		},
		{
			entityType: "climate",
			features: []string{
				string(OnOffClimateEntityFeatures),
				string(HeatClimateEntityFeatures),
				string(CoolClimateEntityFeatures),
				string(CurrentTemperatureClimateEntityFeatures),
				string(TargetTemperatureClimateEntityFeatures),
				string(TargetTemperaturRangeClimateEntityFeatures),
				string(FanClimateEntityFeatures),
			},
		},
		{
			entityType: "sensor",
			// BinarySensorDeviceClass ("binary") is deliberately not checked here: like
			// media_player's device_class above, the vendored YAML's sensor device_class enum
			// doesn't include it, while doc/entities/entity_sensor.md's "Device Classes" section
			// (and its dedicated "Binary Device Class" subsection) documents it as a real,
			// finalized value. The schema lags its own prose here too.
			deviceClass: []string{
				string(CustomSensorDeviceClass),
				string(BatterySensorDeviceClass),
				string(CurrentSensorDeviceClass),
				string(EnergySensorDeviceClass),
				string(HumiditySensorDeviceClass),
				string(PowerSensorDeviceClass),
				string(TemperatureSensorDeviceClass),
				string(VoltageSensorDeviceClass),
			},
		},
		{
			entityType: "remote",
			features: []string{
				string(SendCmdRemoteEntityFeatures),
				string(OnOffRemoteEntityFeatures),
				string(ToggleRemoteEntityFeatures),
			},
		},
		{
			// select has no features and no device_class per the spec.
			entityType: "select",
		},
		{
			entityType: "ir_emitter",
			features: []string{
				string(SendIrEmitterEntityFeatures),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.entityType, func(t *testing.T) {
			schema, err := spec.EntityTypeSchema(c.entityType)
			if err != nil {
				t.Fatalf("spec.EntityTypeSchema(%q): %v", c.entityType, err)
			}

			spectest.AssertSubset(t, "features", c.features, schema.Features)
			spectest.AssertSubset(t, "device_class", c.deviceClass, schema.DeviceClass)

			if len(schema.Features) > len(c.features) {
				t.Logf("spec defines %d feature(s) for %q that goucrt doesn't implement yet: %v",
					len(schema.Features)-len(c.features), c.entityType, missing(schema.Features, c.features))
			}
		})
	}
}

// missing returns the values in a that aren't in b.
func missing(a, b []string) []string {
	inB := make(map[string]bool, len(b))
	for _, v := range b {
		inB[v] = true
	}
	var out []string
	for _, v := range a {
		if !inB[v] {
			out = append(out, v)
		}
	}
	return out
}
