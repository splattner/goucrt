package entities

// Regression tests for AddFeature's attribute registration on newly-added and newly-fixed
// media_player features. spec_features_test.go only checks that a feature's string value exists in
// the vendored spec's enum - it says nothing about what AddFeature actually registers for it, which
// is where the real bugs are: media_type's default value used to be the int 0 even though the spec
// types it as a string, and play_media/search_media's list attributes are new.

import "testing"

func TestMediaPlayerAddFeature_MediaTypeRegistersStringDefaults(t *testing.T) {
	e := NewMediaPlayerEntity("mp-1", LanguageText{En: "Player"}, "", "")
	e.AddFeature(MediaTypeMediaPlayerEntityFeatures)

	mediaType, ok := e.Attributes[string(MediaTypeMediaPlayerEntityAttribute)].(string)
	if !ok {
		t.Fatalf("media_type attribute = %#v (%T), want a string default", e.Attributes[string(MediaTypeMediaPlayerEntityAttribute)], e.Attributes[string(MediaTypeMediaPlayerEntityAttribute)])
	}
	if mediaType != "" {
		t.Errorf("media_type default = %q, want \"\"", mediaType)
	}

	// media_id shares the media_type feature per the spec's Attributes table (both list
	// "media_type" as their enabling feature) - AddFeature must register it too.
	mediaId, ok := e.Attributes[string(MediaIdMediaPlayerEntityAttribute)].(string)
	if !ok {
		t.Fatalf("media_id attribute = %#v (%T), want a string default", e.Attributes[string(MediaIdMediaPlayerEntityAttribute)], e.Attributes[string(MediaIdMediaPlayerEntityAttribute)])
	}
	if mediaId != "" {
		t.Errorf("media_id default = %q, want \"\"", mediaId)
	}
}

func TestMediaPlayerAddFeature_PlayMediaActionRegistersListAttribute(t *testing.T) {
	e := NewMediaPlayerEntity("mp-1", LanguageText{En: "Player"}, "", "")
	e.AddFeature(PlayMediaActionMediaPlayerEntityFeatures)

	actions, ok := e.Attributes[string(PlayMediaActionMediaPlayerEntityAttribute)].([]MediaPlayAction)
	if !ok {
		t.Fatalf("play_media_action attribute = %#v (%T), want []MediaPlayAction", e.Attributes[string(PlayMediaActionMediaPlayerEntityAttribute)], e.Attributes[string(PlayMediaActionMediaPlayerEntityAttribute)])
	}
	if len(actions) != 0 {
		t.Errorf("play_media_action default = %v, want empty", actions)
	}
}

func TestMediaPlayerAddFeature_SearchMediaClassesRegistersListAttribute(t *testing.T) {
	e := NewMediaPlayerEntity("mp-1", LanguageText{En: "Player"}, "", "")
	e.AddFeature(SearchMediaClassesMediaPlayerEntityFeatures)

	classes, ok := e.Attributes[string(SearchMediaClassesMediaPlayerEntityAttribute)].([]MediaClass)
	if !ok {
		t.Fatalf("search_media_classes attribute = %#v (%T), want []MediaClass", e.Attributes[string(SearchMediaClassesMediaPlayerEntityAttribute)], e.Attributes[string(SearchMediaClassesMediaPlayerEntityAttribute)])
	}
	if len(classes) != 0 {
		t.Errorf("search_media_classes default = %v, want empty", classes)
	}
}

// TestMediaPlayerAddFeature_PlayPauseNextPreviousRegisterMediaPlaylist is the regression test for
// media_playlist: the spec's Attributes table lists it under play_pause, next AND previous, but it
// was missing from AddFeature entirely before this change.
func TestMediaPlayerAddFeature_PlayPauseNextPreviousRegisterMediaPlaylist(t *testing.T) {
	for _, feature := range []MediaPlayerEntityFeatures{
		PlayPauseMediaPlayerEntityFeatures,
		NextMediaPlayerEntityFeatures,
		PreviousMediaPlayerEntityFeatures,
	} {
		e := NewMediaPlayerEntity("mp-1", LanguageText{En: "Player"}, "", "")
		e.AddFeature(feature)

		if _, ok := e.Attributes[string(MediaPlaylistMediaPlayerEntityAttribute)]; !ok {
			t.Errorf("AddFeature(%s) did not register media_playlist attribute", feature)
		}
	}
}
