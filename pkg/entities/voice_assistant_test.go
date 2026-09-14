package entities

import "testing"

func TestNewVoiceAssistantEntity_HasStateAttribute(t *testing.T) {
	e := NewVoiceAssistantEntity("va-1", LanguageText{En: "Assistant"}, "")

	state, ok := e.Attributes[string(StateVoiceAssistantEntityAttribute)]
	if !ok {
		t.Fatal("state attribute not set")
	}
	if state != OffVoiceAssistantEntityState {
		t.Errorf("state = %v, want %v", state, OffVoiceAssistantEntityState)
	}
}

// TestVoiceAssistantAddFeature_RegistersNoAttribute is the regression test for AddFeature's
// documented behavior: unlike every other entity type in this package, none of voice_assistant's
// features gate an attribute (they're all reported live via assistant_event instead) - so declaring
// one must change only Features, never Attributes.
func TestVoiceAssistantAddFeature_RegistersNoAttribute(t *testing.T) {
	for _, feature := range []VoiceAssistantEntityFeatures{
		TranscriptionVoiceAssistantEntityFeatures,
		ResponseTextVoiceAssistantEntityFeatures,
		ResponseSpeechVoiceAssistantEntityFeatures,
	} {
		e := NewVoiceAssistantEntity("va-1", LanguageText{En: "Assistant"}, "")
		before := len(e.Attributes)

		e.AddFeature(feature)

		if len(e.Features) != 1 || e.Features[0] != feature {
			t.Errorf("AddFeature(%s): Features = %v, want [%s]", feature, e.Features, feature)
		}
		if len(e.Attributes) != before {
			t.Errorf("AddFeature(%s): Attributes grew from %d to %d entries, want unchanged",
				feature, before, len(e.Attributes))
		}
	}
}

func TestVoiceAssistantEntity_VoiceStartCommandDispatch(t *testing.T) {
	e := NewVoiceAssistantEntity("va-1", LanguageText{En: "Assistant"}, "")

	var gotSessionId float64
	var gotAudioCfg map[string]interface{}
	e.MapCommandWithParams(VoiceStartVoiceAssistantEntityCommand, func(params map[string]interface{}) error {
		gotSessionId, _ = params["session_id"].(float64)
		gotAudioCfg, _ = params["audio_cfg"].(map[string]interface{})
		return nil
	})

	code := e.HandleCommand("voice_start", map[string]interface{}{
		"session_id": float64(8),
		"audio_cfg": map[string]interface{}{
			"channels":      float64(1),
			"sample_rate":   float64(8000),
			"sample_format": "I16",
		},
	})

	if code != 200 {
		t.Errorf("HandleCommand(voice_start) = %d, want 200", code)
	}
	if gotSessionId != 8 {
		t.Errorf("session_id = %v, want 8", gotSessionId)
	}
	if gotAudioCfg["sample_rate"] != float64(8000) {
		t.Errorf("audio_cfg.sample_rate = %v, want 8000", gotAudioCfg["sample_rate"])
	}
}

func TestVoiceAssistantEntity_UnknownCommandReturns404(t *testing.T) {
	e := NewVoiceAssistantEntity("va-1", LanguageText{En: "Assistant"}, "")

	if code := e.HandleCommand("voice_start", nil); code != 404 {
		t.Errorf("HandleCommand(voice_start) with no registered handler = %d, want 404", code)
	}
}
