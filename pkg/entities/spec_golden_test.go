package entities

// Checks the worked examples extracted from the vendored Core-API entity docs (see
// tools/extractgoldens and pkg/entities/testdata) against goucrt's own types and constants.
//
// The AsyncAPI spec doesn't define cmd_id, attribute names or state values as structured schema -
// only doc/entities/*.md's prose and worked examples do (see internal/spec's package doc). These
// fixtures are that prose turned into something a test can check goucrt against mechanically.
//
// Two different checks, deliberately different in strictness:
//   - Every fixture must decode and carry the entity_type its directory claims. This is a hard
//     failure: a fixture that doesn't decode, or claims the wrong entity_type, means the extraction
//     (or a fixture someone hand-edited) is broken, not that goucrt has a gap.
//   - A fixture's cmd_id/attribute keys are compared against goucrt's declared constants for that
//     entity, but a miss is only logged, not failed. goucrt's command dispatch is a plain map
//     (AddCommand/Commands), so an undeclared constant doesn't mean the command can't work - only
//     that there's no named Go constant a driver author would discover for it. A missing attribute
//     constant just means a newer spec addition (e.g. climate's hvac_mode naming, media_player's
//     sized media_image_url_* attributes) that goucrt hasn't caught up to yet. Both are real,
//     worth knowing about, and neither is this test's job to fix.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type goldenMessage struct {
	Kind    string `json:"kind"`
	Msg     string `json:"msg"`
	MsgData struct {
		EntityType string                 `json:"entity_type"`
		EntityID   string                 `json:"entity_id"`
		CmdID      string                 `json:"cmd_id"`
		Params     map[string]interface{} `json:"params"`
		Attributes map[string]interface{} `json:"attributes"`
	} `json:"msg_data"`
}

// knownCommands and knownAttributes list the string values of every X*EntityCommand /
// X*EntityAttribute[s] constant goucrt currently declares for the entity type named by the map key -
// the same source data as spec_features_test.go's per-entity constant lists, but for the properties
// the spec doesn't expose as a structured enum.
var knownCommands = map[string][]string{
	"button": {string(PushButtonEntityCommand)},
	"switch": {
		string(OnSwitchEntityCommand),
		string(OffSwitchEntityCommand),
		string(ToggleSwitchEntityCommand),
	},
	"light": {
		string(OnLightEntityCommand),
		string(OffLightEntityCommand),
		string(ToggleLightEntityCommand),
	},
	"cover": {
		string(OpenCoverEntityCommand),
		string(CloseCoverEntityCommand),
		string(StopCoverEntityCommand),
		string(PositionCoverEntityCommand),
		string(TiltCoverEntityCommand),
		string(TiltUpCoverEntityCommand),
		string(TiltDownCoverEntityCommand),
		string(TiltStopCoverEntityCommand),
	},
	"media_player": {
		string(OnMediaPlayerEntityCommand),
		string(OffMediaPlayerEntityCommand),
		string(ToggleMediaPlayerEntityCommand),
		string(PlayPauseMediaPlayerEntityCommand),
		string(StopMediaPlayerEntityCommand),
		string(PreviousMediaPlayerEntityCommand),
		string(NextMediaPlayerEntityCommand),
		string(FastForwardMediaPlayerEntityCommand),
		string(RewindMediaPlayerEntityCommand),
		string(SeekMediaPlayerEntityCommand),
		string(VolumeMediaPlayerEntityCommand),
		string(VolumeUpMediaPlayerEntityCommand),
		string(VolumeDownMediaPlayerEntityCommand),
		string(MuteToggleMediaPlayerEntityCommand),
		string(MuteMediaPlayerEntityCommand),
		string(UnmuteMediaPlayerEntityCommand),
		string(RepeatMediaPlayerEntityCommand),
		string(ShuffleMediaPlayerEntityCommand),
		string(ChannelUpMediaPlayerEntityCommand),
		string(ChannelDownMediaPlayerEntityCommand),
		string(CursorUpMediaPlayerEntityCommand),
		string(CursorDownMediaPlayerEntityCommand),
		string(CursorLeftMediaPlayerEntityCommand),
		string(CursorRightMediaPlayerEntityCommand),
		string(CursorEnterMediaPlayerEntityCommand),
		string(HomeMediaPlayerEntityCommand),
		string(MenuMediaPlayerEntityCommand),
		string(ContextMenuMediaPlayerEntityCommand),
		string(GuideMediaPlayerEntityCommand),
		string(InfoMediaPlayerEntityCommand),
		string(BackMediaPlayerEntityCommand),
		string(SelectSourceMediaPlayerEntityCommand),
		string(SelectSoundModeMediaPlayerEntityCommand),
		string(RecordMediaPlayerEntityCommand),
		string(EjectMediaPlayerEntityCommand),
		string(OpenCloseMediaPlayerEntityCommand),
		string(AudioTrackMediaPlayerEntityCommand),
		string(SubtitleMediaPlayerEntityCommand),
		string(SettingsMediaPlayerEntityCommand),
		string(SearchMediaPlayerEntityCommand),
		string(PlayMediaMediaPlayerEntityCommand),
		string(ClearPlaylistMediaPlayerEntityCommand),
	},
	"climate": {
		string(OnClimateEntityCommand),
		string(OffClimateEntityCommand),
		string(HVACModeClimateEntityCommand),
		string(TargetTemperatureClimateEntityCommand),
		string(TargetTemperatureRangeClimateEntityCommand),
		string(FanModeClimateEntityCommand),
	},
	"remote": {
		string(OnRemoteEntityCommand),
		string(OffRemoteEntityCommand),
		string(SendCmdRemoteEntityCommand),
		string(SendCmdSequenceRemoteEntityCommand),
	},
	"select": {
		string(SelectOptionSelectEntityCommand),
		string(SelectFirstSelectEntityCommand),
		string(SelectLastSelectEntityCommand),
		string(SelectNextSelectEntityCommand),
		string(SelectPreviousSelectEntityCommand),
	},
	"ir_emitter": {
		string(SendIrEmitterEntityCommand),
		string(StopIrEmitterEntityCommand),
	},
}

var knownAttributes = map[string][]string{
	"switch": {string(StateSwitchEntityAttribute)},
	"light": {
		string(StateLightEntityAttribute),
		string(HueLightEntityAttribute),
		string(SaturationLightEntityAttribute),
		string(BrightnessLightEntityAttribute),
		string(ColorTemperatureLightEntityAttribute),
	},
	"cover": {
		string(StateCoverEntityAttribute),
		string(PositionCoverEntityAttribute),
		string(TiltPositionCoverEntityAttribute),
	},
	"media_player": {
		string(StateMediaPlayerEntityAttribute),
		string(VolumeMediaPlayerEntityAttribute),
		string(MutedMediaPlayerEntityAttribute),
		string(MediaDurationMediaPlayerEntityAttribute),
		string(MediaPositionMediaPlayerEntityAttribute),
		string(MediaTypeMediaPlayerEntityAttribute),
		string(MediaImageUrlMediaPlayerEntityAttribute),
		string(MediaTitleMediaPlayerEntityAttribute),
		string(MediaArtistMediaPlayerEntityAttribute),
		string(MediaAlbumMediaPlayerEntityAttribute),
		string(RepeatMediaPlayerEntityAttribute),
		string(ShuffleMediaPlayerEntityAttribute),
		string(SourceMediaPlayerEntityAttribute),
		string(SourceListMediaPlayerEntityAttribute),
		string(SoundModeMediaPlayerEntityAttribute),
		string(SoundModeListMediaPlayerEntityAttribute),
		string(MediaIdMediaPlayerEntityAttribute),
		string(MediaPlaylistMediaPlayerEntityAttribute),
		string(PlayMediaActionMediaPlayerEntityAttribute),
		string(SearchMediaClassesMediaPlayerEntityAttribute),
	},
	"climate": {
		string(StateClimateEntityAttribute),
		string(CurrentTemperatureClimateEntityAttribute),
		string(TargetTemperatureClimateEntityAttribute),
		string(TargetTemperatureHighClimateEntityAttribute),
		string(TargetTemperatureLowClimateEntityAttribute),
		string(FanModeClimateEntityAttribute),
	},
	"sensor": {
		string(StateSensorEntityAttribute),
		string(ValueSensorEntityAttribute),
		string(UnitSensorEntityAttribute),
	},
	"remote": {string(StateRemoteEntityAttribute)},
	"select": {
		string(StateSelectEntityAttribute),
		string(CurrentOptionSelectEntityAttribute),
		string(OptionsSelectEntityAttribute),
	},
	"ir_emitter": {string(StateIrEmitterEntityAttribute)},
}

func TestGoldenExamplesMatchSpec(t *testing.T) {
	entityDirs, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatalf("read testdata: %v", err)
	}
	if len(entityDirs) == 0 {
		t.Fatal("testdata is empty - run `go run ./tools/extractgoldens` from the repository root")
	}

	for _, entityDir := range entityDirs {
		if !entityDir.IsDir() {
			continue
		}
		entityType := entityDir.Name()

		t.Run(entityType, func(t *testing.T) {
			checkFixtures(t, entityType, "commands", func(t *testing.T, name string, msg goldenMessage) {
				if msg.MsgData.CmdID == "" {
					t.Logf("%s: no cmd_id field (spec doc quirk - e.g. entity_button.md's one example uses the stale key \"command_id\") - skipping", name)
					return
				}
				if !contains(knownCommands[entityType], msg.MsgData.CmdID) {
					t.Logf("%s: spec example uses cmd_id %q, which goucrt has no named %sEntityCommand constant for (may still work via AddCommand's map-based dispatch)",
						name, msg.MsgData.CmdID, entityType)
				}
			})

			checkFixtures(t, entityType, "events", func(t *testing.T, name string, msg goldenMessage) {
				known := knownAttributes[entityType]
				var unknown []string
				for key := range msg.MsgData.Attributes {
					if !contains(known, key) {
						unknown = append(unknown, key)
					}
				}
				if len(unknown) > 0 {
					t.Logf("%s: spec example uses attribute(s) %v, which goucrt has no matching constant for yet", name, unknown)
				}
			})
		})
	}
}

// checkFixtures decodes every fixture under testdata/<entityType>/<sub>/*.json, hard-fails on a
// decode error or an entity_type mismatch, and otherwise hands the decoded message to check.
func checkFixtures(t *testing.T, entityType, sub string, check func(t *testing.T, name string, msg goldenMessage)) {
	t.Helper()

	dir := filepath.Join("testdata", entityType, sub)
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return // not every entity type has both commands and events examples (e.g. sensor has no commands).
	}
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}

	for _, e := range entries {
		name := filepath.Join(sub, e.Name())
		t.Run(name, func(t *testing.T) {
			raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				t.Fatalf("read fixture: %v", err)
			}

			var msg goldenMessage
			if err := json.Unmarshal(raw, &msg); err != nil {
				t.Fatalf("decode fixture: %v", err)
			}

			if msg.MsgData.EntityType != entityType {
				t.Fatalf("fixture's entity_type = %q, want %q (extraction put this file in the wrong directory)",
					msg.MsgData.EntityType, entityType)
			}

			check(t, name, msg)
		})
	}
}

func contains(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return false
}
