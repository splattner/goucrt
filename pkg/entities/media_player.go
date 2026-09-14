package entities

import (
	"fmt"
	"slices"
)

type MediaPlayerEntityState EntityState
type MediaPlayerEntityFeatures EntityFeature
type MediaPlayerEntityAttributes EntityAttribute
type MediaPlayerEntityCommand EntityCommand
type MediaPlayerEntityOption EntityOption
type MediaPlayerDeviceClass string

const (
	OnMediaPlayerEntityState        MediaPlayerEntityState = "ON"
	OffMediaPlayerEntityState       MediaPlayerEntityState = "OFF"
	PlayingPlayerEntityState        MediaPlayerEntityState = "PLAYING"
	PausedMediaPlayerEntityState    MediaPlayerEntityState = "PAUSED"
	StandbyMediaPlayerEntityState   MediaPlayerEntityState = "STANDBY"
	BufferingMediaPlayerEntityState MediaPlayerEntityState = "BUFFERING"
)

const (
	OnOffMediaPlayerEntityFeatures                  MediaPlayerEntityFeatures = "on_off"
	ToggleMediaPlayerEntityFeatures                 MediaPlayerEntityFeatures = "toggle"
	VolumeMediaPlayerEntityFeatures                 MediaPlayerEntityFeatures = "volume"
	VolumeUpDownMediaPlayerEntityFeatures           MediaPlayerEntityFeatures = "volume_up_down"
	MuteToggleMediaPlayerEntityFeatures             MediaPlayerEntityFeatures = "mute_toggle"
	MuteMediaPlayerEntityFeatures                   MediaPlayerEntityFeatures = "mute"
	UnmuteMediaPlayerEntityFeatures                 MediaPlayerEntityFeatures = "unmute"
	PlayPauseMediaPlayerEntityFeatures              MediaPlayerEntityFeatures = "play_pause"
	StopMediaPlayerEntityFeatures                   MediaPlayerEntityFeatures = "stop"
	NextMediaPlayerEntityFeatures                   MediaPlayerEntityFeatures = "next"
	PreviousMediaPlayerEntityFeatures               MediaPlayerEntityFeatures = "previous"
	FastForwardMediaPlayerEntityFeatures            MediaPlayerEntityFeatures = "fast_forward"
	RewindMediaPlayerEntityFeatures                 MediaPlayerEntityFeatures = "rewind"
	RepeatMediaPlayerEntityFeatures                 MediaPlayerEntityFeatures = "repeat"
	ShuffleMediaPlayerEntityFeatures                MediaPlayerEntityFeatures = "shuffle"
	SeekMediaPlayerEntityFeatures                   MediaPlayerEntityFeatures = "seek"
	MediaDurationMediaPlayerEntityFeatures          MediaPlayerEntityFeatures = "media_duration"
	MediaPositionMediaPlayerEntityFeatures          MediaPlayerEntityFeatures = "media_position"
	MediaPositionUpdatedAtMediaPlayerEntityFeatures MediaPlayerEntityFeatures = "media_position_updated_at"
	MediaTitleMediaPlayerEntityFeatures             MediaPlayerEntityFeatures = "media_title"
	MediaArtistMediaPlayerEntityFeatures            MediaPlayerEntityFeatures = "media_artist"
	MediaAlbumMediaPlayerEntityFeatures             MediaPlayerEntityFeatures = "media_album"
	MediaImageUrlMediaPlayerEntityFeatures          MediaPlayerEntityFeatures = "media_image_url"
	MediaTypeMediaPlayerEntityFeatures              MediaPlayerEntityFeatures = "media_type"
	DPadMediaPlayerEntityFeatures                   MediaPlayerEntityFeatures = "dpad"
	NumPadMediaPlayerEntityFeatures                 MediaPlayerEntityFeatures = "numpad"
	HomeMediaPlayerEntityFeatures                   MediaPlayerEntityFeatures = "home"
	MenuMediaPlayerEntityFeatures                   MediaPlayerEntityFeatures = "menu"
	ContextMenuPlayerEntityFeatures                 MediaPlayerEntityFeatures = "context_menu"
	GuidePlayerEntityFeatures                       MediaPlayerEntityFeatures = "guide"
	InfoPlayerEntityFeatures                        MediaPlayerEntityFeatures = "info"
	ColorButtonsMediaPlayerEntityFeatures           MediaPlayerEntityFeatures = "color_buttons"
	ChannelSwitcherMediaPlayerEntityFeatures        MediaPlayerEntityFeatures = "channel_switcher"
	SelectSourceMediaPlayerEntityFeatures           MediaPlayerEntityFeatures = "select_source"
	SelectSoundModeMediaPlayerEntityFeatures        MediaPlayerEntityFeatures = "select_sound_mode"
	EjectMediaPlayerEntityFeatures                  MediaPlayerEntityFeatures = "eject"
	OpenCloseMediaPlayerEntityFeatures              MediaPlayerEntityFeatures = "open_close"
	AudioTrackMediaPlayerEntityFeatures             MediaPlayerEntityFeatures = "audio_track"
	SubtitleMediaPlayerEntityFeatures               MediaPlayerEntityFeatures = "subtitle"
	RecordMediaPlayerEntityFeatures                 MediaPlayerEntityFeatures = "record"
	SettingsMediaPlayerEntityFeatures               MediaPlayerEntityFeatures = "settings"
	// BrowseMediaMediaPlayerEntityFeatures: the entity supports the browse_media request. See
	// SetBrowseFunc.
	BrowseMediaMediaPlayerEntityFeatures MediaPlayerEntityFeatures = "browse_media"
	// SearchMediaMediaPlayerEntityFeatures: the entity supports the search_media request. See
	// SetSearchFunc.
	SearchMediaMediaPlayerEntityFeatures MediaPlayerEntityFeatures = "search_media"
	// PlayMediaMediaPlayerEntityFeatures: the entity supports the play_media command.
	PlayMediaMediaPlayerEntityFeatures MediaPlayerEntityFeatures = "play_media"
	// PlayMediaActionMediaPlayerEntityFeatures: the entity supports play_media's optional "action"
	// parameter to play now, play next or enqueue. See MediaPlayAction.
	PlayMediaActionMediaPlayerEntityFeatures MediaPlayerEntityFeatures = "play_media_action"
	// ClearPlaylistMediaPlayerEntityFeatures: the entity supports the clear_playlist command.
	ClearPlaylistMediaPlayerEntityFeatures MediaPlayerEntityFeatures = "clear_playlist"
	// SearchMediaClassesMediaPlayerEntityFeatures: the entity provides a list of MediaClass values
	// (in the search_media_classes attribute) to use as a filter for search_media.
	SearchMediaClassesMediaPlayerEntityFeatures MediaPlayerEntityFeatures = "search_media_classes"
)

// Deprecated: misspelled aliases kept for backward compatibility, will be removed in a future release.
const (
	// Deprecated: use ToggleMediaPlayerEntityFeatures instead.
	ToggleMediaPlayerEntityyFeatures = ToggleMediaPlayerEntityFeatures
	// Deprecated: use VolumeMediaPlayerEntityFeatures instead.
	VolumeMediaPlayerEntityyFeatures = VolumeMediaPlayerEntityFeatures
	// Deprecated: use PreviousMediaPlayerEntityFeatures instead.
	PreviusMediaPlayerEntityFeatures = PreviousMediaPlayerEntityFeatures
)

const (
	StateMediaPlayerEntityAttribute         MediaPlayerEntityAttributes = "state"
	VolumeMediaPlayerEntityAttribute        MediaPlayerEntityAttributes = "volume"
	MutedMediaPlayerEntityAttribute         MediaPlayerEntityAttributes = "muted"
	MediaDurationMediaPlayerEntityAttribute MediaPlayerEntityAttributes = "media_duration"
	MediaPositionMediaPlayerEntityAttribute MediaPlayerEntityAttributes = "media_position"
	// MediaPositionUpdatedAtMediaPlayerEntityAttribute: string, ISO 8601 datetime. Optional
	// timestamp when media_position was last updated. Requires integration support.
	MediaPositionUpdatedAtMediaPlayerEntityAttribute MediaPlayerEntityAttributes = "media_position_updated_at"
	MediaTypeMediaPlayerEntityAttribute              MediaPlayerEntityAttributes = "media_type"
	MediaImageUrlMediaPlayerEntityAttribute          MediaPlayerEntityAttributes = "media_image_url"
	MediaTitleMediaPlayerEntityAttribute             MediaPlayerEntityAttributes = "media_title"
	MediaArtistMediaPlayerEntityAttribute            MediaPlayerEntityAttributes = "media_artist"
	MediaAlbumMediaPlayerEntityAttribute             MediaPlayerEntityAttributes = "media_album"
	RepeatMediaPlayerEntityAttribute                 MediaPlayerEntityAttributes = "repeat"
	ShuffleMediaPlayerEntityAttribute                MediaPlayerEntityAttributes = "shuffle"
	SourceMediaPlayerEntityAttribute                 MediaPlayerEntityAttributes = "source"
	SourceListMediaPlayerEntityAttribute             MediaPlayerEntityAttributes = "source_list"
	SoundModeMediaPlayerEntityAttribute              MediaPlayerEntityAttributes = "sound_mode"
	SoundModeListMediaPlayerEntityAttribute          MediaPlayerEntityAttributes = "sound_mode_list"
	// MediaIdMediaPlayerEntityAttribute: the content ID of the media being played. Enabled by the
	// media_type feature, same as MediaTypeMediaPlayerEntityAttribute.
	MediaIdMediaPlayerEntityAttribute MediaPlayerEntityAttributes = "media_id"
	// MediaPlaylistMediaPlayerEntityAttribute: title of the playlist currently playing. Enabled by
	// play_pause, next or previous, same as the other now-playing attributes.
	MediaPlaylistMediaPlayerEntityAttribute MediaPlayerEntityAttributes = "media_playlist"
	// PlayMediaActionMediaPlayerEntityAttribute: the MediaPlayAction values this entity supports
	// for play_media's optional "action" parameter.
	PlayMediaActionMediaPlayerEntityAttribute MediaPlayerEntityAttributes = "play_media_action"
	// SearchMediaClassesMediaPlayerEntityAttribute: the MediaClass values this entity accepts as a
	// search_media filter.
	SearchMediaClassesMediaPlayerEntityAttribute MediaPlayerEntityAttributes = "search_media_classes"
)

// Deprecated: use MutedMediaPlayerEntityAttribute instead.
const MutedMediaPlayeEntityAttribute = MutedMediaPlayerEntityAttribute

const (
	OnMediaPlayerEntityCommand               MediaPlayerEntityCommand = "on"
	OffMediaPlayerEntityCommand              MediaPlayerEntityCommand = "off"
	ToggleMediaPlayerEntityCommand           MediaPlayerEntityCommand = "toggle"
	PlayPauseMediaPlayerEntityCommand        MediaPlayerEntityCommand = "play_pause"
	StopMediaPlayerEntityCommand             MediaPlayerEntityCommand = "stop"
	PreviousMediaPlayerEntityCommand         MediaPlayerEntityCommand = "previous"
	NextMediaPlayerEntityCommand             MediaPlayerEntityCommand = "next"
	FastForwardMediaPlayerEntityCommand      MediaPlayerEntityCommand = "fast_forward"
	RewindMediaPlayerEntityCommand           MediaPlayerEntityCommand = "rewind"
	SeekMediaPlayerEntityCommand             MediaPlayerEntityCommand = "seek"
	VolumeMediaPlayerEntityCommand           MediaPlayerEntityCommand = "volume"
	VolumeUpMediaPlayerEntityCommand         MediaPlayerEntityCommand = "volume_up"
	VolumeDownMediaPlayerEntityCommand       MediaPlayerEntityCommand = "volume_down"
	MuteToggleMediaPlayerEntityCommand       MediaPlayerEntityCommand = "mute_toggle"
	MuteMediaPlayerEntityCommand             MediaPlayerEntityCommand = "mute"
	UnmuteMediaPlayerEntityCommand           MediaPlayerEntityCommand = "unmute"
	RepeatMediaPlayerEntityCommand           MediaPlayerEntityCommand = "repeat"
	ShuffleMediaPlayerEntityCommand          MediaPlayerEntityCommand = "shuffle"
	ChannelUpMediaPlayerEntityCommand        MediaPlayerEntityCommand = "channel_up"
	ChannelDownMediaPlayerEntityCommand      MediaPlayerEntityCommand = "channel_down"
	CursorUpMediaPlayerEntityCommand         MediaPlayerEntityCommand = "cursor_up"
	CursorDownMediaPlayerEntityCommand       MediaPlayerEntityCommand = "cursor_down"
	CursorLeftMediaPlayerEntityCommand       MediaPlayerEntityCommand = "cursor_left"
	CursorRightMediaPlayerEntityCommand      MediaPlayerEntityCommand = "cursor_right"
	CursorEnterMediaPlayerEntityCommand      MediaPlayerEntityCommand = "cursor_enter"
	Digit0MediaPlayerEntityCommand           MediaPlayerEntityCommand = "digit_0"
	Digit1MediaPlayerEntityCommand           MediaPlayerEntityCommand = "digit_1"
	Digit2MediaPlayerEntityCommand           MediaPlayerEntityCommand = "digit_2"
	Digit3MediaPlayerEntityCommand           MediaPlayerEntityCommand = "digit_3"
	Digit4MediaPlayerEntityCommand           MediaPlayerEntityCommand = "digit_4"
	Digit5MediaPlayerEntityCommand           MediaPlayerEntityCommand = "digit_5"
	Digit6MediaPlayerEntityCommand           MediaPlayerEntityCommand = "digit_6"
	Digit7MediaPlayerEntityCommand           MediaPlayerEntityCommand = "digit_7"
	Digit8MediaPlayerEntityCommand           MediaPlayerEntityCommand = "digit_8"
	Digit9MediaPlayerEntityCommand           MediaPlayerEntityCommand = "digit_9"
	FunctionRedMediaPlayerEntityCommand      MediaPlayerEntityCommand = "function_red"
	FunctionGreenMediaPlayerEntityCommand    MediaPlayerEntityCommand = "function_green"
	FunctionYellowMediaPlayerEntityCommand   MediaPlayerEntityCommand = "function_yellow"
	FunctionBlueMediaPlayerEntityCommand     MediaPlayerEntityCommand = "function_blue"
	HomeMediaPlayerEntityCommand             MediaPlayerEntityCommand = "home"
	MenuMediaPlayerEntityCommand             MediaPlayerEntityCommand = "menu"
	ContextMenuMediaPlayerEntityCommand      MediaPlayerEntityCommand = "context_menu"
	GuideMediaPlayerEntityCommand            MediaPlayerEntityCommand = "guide"
	InfoMediaPlayerEntityCommand             MediaPlayerEntityCommand = "info"
	BackMediaPlayerEntityCommand             MediaPlayerEntityCommand = "back"
	SelectSourceMediaPlayerEntityCommand     MediaPlayerEntityCommand = "select_source"
	SelectSoundModeMediaPlayerEntityCommand  MediaPlayerEntityCommand = "select_sound_mode"
	RecordMediaPlayerEntityCommand           MediaPlayerEntityCommand = "record"
	MyRecordingsMenuMediaPlayerEntityCommand MediaPlayerEntityCommand = "my_recordings"
	LiveMediaPlayerEntityCommand             MediaPlayerEntityCommand = "live"
	EjectMediaPlayerEntityCommand            MediaPlayerEntityCommand = "eject"
	OpenCloseMediaPlayerEntityCommand        MediaPlayerEntityCommand = "open_close"
	AudioTrackMediaPlayerEntityCommand       MediaPlayerEntityCommand = "audio_track"
	SubtitleMediaPlayerEntityCommand         MediaPlayerEntityCommand = "subtitle"
	SettingsMediaPlayerEntityCommand         MediaPlayerEntityCommand = "settings"
	SearchMediaPlayerEntityCommand           MediaPlayerEntityCommand = "search"
	// PlayMediaMediaPlayerEntityCommand: play or enqueue a media item. Takes required "media_id"
	// and "media_type" params, plus an optional "action" (MediaPlayAction) param if the
	// PlayMediaActionMediaPlayerEntityFeatures feature is supported - PlayNowMediaPlayAction is the
	// default if omitted.
	PlayMediaMediaPlayerEntityCommand MediaPlayerEntityCommand = "play_media"
	// ClearPlaylistMediaPlayerEntityCommand: remove all items from the playback queue. What happens
	// to the currently playing item (kept playing or also cleared) is integration-dependent.
	ClearPlaylistMediaPlayerEntityCommand MediaPlayerEntityCommand = "clear_playlist"
)

// MediaPlayAction is play_media's optional "action" parameter: how to handle a newly played item
// relative to whatever's already in the queue. An integration may use its own values beyond the
// three predefined here, though the UI may not have locale-aware labels for them.
type MediaPlayAction string

const (
	PlayNowMediaPlayAction    MediaPlayAction = "PLAY_NOW"
	PlayNextMediaPlayAction   MediaPlayAction = "PLAY_NEXT"
	AddToQueueMediaPlayAction MediaPlayAction = "ADD_TO_QUEUE"
)

// Deprecated: misspelled aliases kept for backward compatibility, will be removed in a future release.
const (
	// Deprecated: use PreviousMediaPlayerEntityCommand instead.
	PreviusMediaPlayerEntityCommand = PreviousMediaPlayerEntityCommand
	// Deprecated: use SelectSourceMediaPlayerEntityCommand instead.
	SelectSourcMediaPlayerEntityCommand = SelectSourceMediaPlayerEntityCommand
)

const (
	ReceiverMediaPlayerDeviceClass     MediaPlayerDeviceClass = "receiver"
	SetTopBoxMediaPlayerDeviceClass    MediaPlayerDeviceClass = "set_top_box"
	SpeakerMediaPlayerDeviceClass      MediaPlayerDeviceClass = "speaker"
	StreamingBoxMediaPlayerDeviceClass MediaPlayerDeviceClass = "streaming_box"
	TVMediaPlayerDeviceClass           MediaPlayerDeviceClass = "tv"
)

// Deprecated: use StreamingBoxMediaPlayerDeviceClass instead.
const StreamingBoxMMediaPlayerDeviceClass = StreamingBoxMediaPlayerDeviceClass

const (
	SimpleCommandsMediaPlayerEntityOption MediaPlayerEntityOption = "simple_commands"
	VolumeStepsMediaPlayerEntityOption    MediaPlayerEntityOption = "volume_steps"
)

type MediaPlayerEntity struct {
	BaseEntity
	DeviceClass MediaPlayerDeviceClass                                                           `json:"device_class,omitempty"`
	Commands    map[MediaPlayerEntityCommand]func(MediaPlayerEntity, map[string]interface{}) int `json:"-"`
	Options     map[MediaPlayerEntityOption]interface{}                                          `json:"options,omitempty"`
	// BrowseFunc/SearchFunc handle browse_media/search_media requests for this entity. Set via
	// SetBrowseFunc/SetSearchFunc; nil means unsupported, even if the corresponding feature is
	// declared.
	BrowseFunc func(BrowseMediaRequest) (*BrowseMediaResult, error) `json:"-"`
	SearchFunc func(SearchMediaRequest) (*SearchMediaResult, error) `json:"-"`
}

func NewMediaPlayerEntity(id string, name LanguageText, area string, deviceClass MediaPlayerDeviceClass) *MediaPlayerEntity {

	mediaPlayerEntity := MediaPlayerEntity{}
	mediaPlayerEntity.Id = id
	mediaPlayerEntity.Name = name
	mediaPlayerEntity.Area = area
	mediaPlayerEntity.DeviceClass = deviceClass

	mediaPlayerEntity.Type = "media_player"

	mediaPlayerEntity.Commands = make(map[MediaPlayerEntityCommand]func(MediaPlayerEntity, map[string]interface{}) int)
	mediaPlayerEntity.Attributes = make(map[string]interface{})

	mediaPlayerEntity.Options = make(map[MediaPlayerEntityOption]interface{})

	return &mediaPlayerEntity
}

func (e *MediaPlayerEntity) UpdateEntity(newEntity interface{}) error {
	updated, ok := newEntity.(MediaPlayerEntity)
	if !ok {
		return fmt.Errorf("cannot update MediaPlayerEntity from %T", newEntity)
	}

	e.Name = updated.Name
	e.Area = updated.Area
	e.Commands = updated.Commands
	e.Features = updated.Features
	e.Attributes = updated.Attributes
	e.Options = updated.Options
	e.BrowseFunc = updated.BrowseFunc
	e.SearchFunc = updated.SearchFunc

	return nil
}

// Register a function for the Entity command
// Based on the Feature, the correct Attributes will be added
func (e *MediaPlayerEntity) AddFeature(feature MediaPlayerEntityFeatures) {
	e.Features = append(e.Features, feature)

	// Add Attributes based on enabled features
	// https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_media_player.md
	switch feature {
	case OnOffMediaPlayerEntityFeatures:
		e.AddAttribute(string(StateMediaPlayerEntityAttribute), OffMediaPlayerEntityState)

	case ToggleMediaPlayerEntityFeatures:
		e.AddAttribute(string(StateMediaPlayerEntityAttribute), OffMediaPlayerEntityState)

	case PlayPauseMediaPlayerEntityFeatures:
		e.AddAttribute(string(StateMediaPlayerEntityAttribute), OffMediaPlayerEntityState)
		e.AddAttribute(string(MediaPositionMediaPlayerEntityAttribute), 0)
		e.AddAttribute(string(MediaImageUrlMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(MediaTitleMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(MediaArtistMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(MediaAlbumMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(MediaPlaylistMediaPlayerEntityAttribute), "")

	case StopMediaPlayerEntityFeatures:
		e.AddAttribute(string(StateMediaPlayerEntityAttribute), OffMediaPlayerEntityState)
		e.AddAttribute(string(MediaPositionMediaPlayerEntityAttribute), 0)

	case VolumeMediaPlayerEntityFeatures:
		e.AddAttribute(string(VolumeMediaPlayerEntityAttribute), 0)

	case VolumeUpDownMediaPlayerEntityFeatures:
		e.AddAttribute(string(VolumeMediaPlayerEntityAttribute), 0)

	case MuteToggleMediaPlayerEntityFeatures:
		e.AddAttribute(string(MutedMediaPlayerEntityAttribute), false)

	case MuteMediaPlayerEntityFeatures:
		e.AddAttribute(string(MutedMediaPlayerEntityAttribute), false)

	case UnmuteMediaPlayerEntityFeatures:
		e.AddAttribute(string(MutedMediaPlayerEntityAttribute), false)

	case MediaDurationMediaPlayerEntityFeatures:
		e.AddAttribute(string(MediaDurationMediaPlayerEntityAttribute), 0)

	case MediaPositionMediaPlayerEntityFeatures:
		e.AddAttribute(string(MediaPositionMediaPlayerEntityAttribute), 0)

	case MediaPositionUpdatedAtMediaPlayerEntityFeatures:
		e.AddAttribute(string(MediaPositionUpdatedAtMediaPlayerEntityAttribute), "")

	case FastForwardMediaPlayerEntityFeatures:
		e.AddAttribute(string(MediaPositionMediaPlayerEntityAttribute), 0)

	case RewindMediaPlayerEntityFeatures:
		e.AddAttribute(string(MediaPositionMediaPlayerEntityAttribute), 0)

	case SeekMediaPlayerEntityFeatures:
		e.AddAttribute(string(MediaPositionMediaPlayerEntityAttribute), 0)

	case MediaTypeMediaPlayerEntityFeatures:
		// media_id shares this feature (see the spec's Attributes table: both media_id and
		// media_type list "media_type" as their enabling feature, there's no separate media_id
		// feature). media_type's default was previously the int 0 rather than a string, which
		// doesn't match its spec type (string, one of the Media Content Types).
		e.AddAttribute(string(MediaTypeMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(MediaIdMediaPlayerEntityAttribute), "")

	case MediaImageUrlMediaPlayerEntityFeatures:
		e.AddAttribute(string(MediaImageUrlMediaPlayerEntityAttribute), "")

	case NextMediaPlayerEntityFeatures:
		e.AddAttribute(string(MediaImageUrlMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(MediaTitleMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(MediaArtistMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(MediaAlbumMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(MediaPlaylistMediaPlayerEntityAttribute), "")

	case PreviousMediaPlayerEntityFeatures:
		e.AddAttribute(string(MediaImageUrlMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(MediaTitleMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(MediaArtistMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(MediaAlbumMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(MediaPlaylistMediaPlayerEntityAttribute), "")

	case MediaTitleMediaPlayerEntityFeatures:
		e.AddAttribute(string(MediaTitleMediaPlayerEntityAttribute), "")

	case MediaArtistMediaPlayerEntityFeatures:
		e.AddAttribute(string(MediaArtistMediaPlayerEntityAttribute), "")

	case MediaAlbumMediaPlayerEntityFeatures:
		e.AddAttribute(string(MediaAlbumMediaPlayerEntityAttribute), "")

	case RepeatMediaPlayerEntityFeatures:
		e.AddAttribute(string(RepeatMediaPlayerEntityAttribute), 0)

	case ShuffleMediaPlayerEntityFeatures:
		e.AddAttribute(string(ShuffleMediaPlayerEntityAttribute), false)

	case SelectSourceMediaPlayerEntityFeatures:
		e.AddAttribute(string(SourceMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(SourceListMediaPlayerEntityAttribute), []string{})

	case SelectSoundModeMediaPlayerEntityFeatures:
		e.AddAttribute(string(SoundModeMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(SoundModeListMediaPlayerEntityAttribute), []string{})

	case PlayMediaActionMediaPlayerEntityFeatures:
		e.AddAttribute(string(PlayMediaActionMediaPlayerEntityAttribute), []MediaPlayAction{})

	case SearchMediaClassesMediaPlayerEntityFeatures:
		e.AddAttribute(string(SearchMediaClassesMediaPlayerEntityAttribute), []MediaClass{})

	}
}

// Register a function for the Entity command
func (e *MediaPlayerEntity) AddCommand(command MediaPlayerEntityCommand, function func(MediaPlayerEntity, map[string]interface{}) int) {
	e.Commands[command] = function

}

// Map a Light EntityCommand to a function call with params
func (e *MediaPlayerEntity) MapCommandWithParams(command MediaPlayerEntityCommand, f func(map[string]interface{}) error) {

	e.AddCommand(command, func(entity MediaPlayerEntity, params map[string]interface{}) int {

		if err := f(params); err != nil {
			return 404
		}
		return 200
	})
}

// Map a Light EntityCommand to a function call without params
func (e *MediaPlayerEntity) MapCommand(command MediaPlayerEntityCommand, f func() error) {

	e.AddCommand(command, func(entity MediaPlayerEntity, params map[string]interface{}) int {

		if err := f(); err != nil {
			return 404
		}
		return 200
	})
}

// Call the registred function for this entity_command
func (e *MediaPlayerEntity) HandleCommand(cmd_id string, params map[string]interface{}) int {
	if e.Commands[MediaPlayerEntityCommand(cmd_id)] != nil {
		return e.Commands[MediaPlayerEntityCommand(cmd_id)](*e, params)
	}

	// When simple_commands are enabled and the command exists, call the regstisered function if one is set
	if e.Options[SimpleCommandsMediaPlayerEntityOption] != nil &&
		slices.Contains(e.Options[SimpleCommandsMediaPlayerEntityOption].([]string), cmd_id) {

		if e.Commands[MediaPlayerEntityCommand(cmd_id)] != nil {
			return e.Commands[MediaPlayerEntityCommand(cmd_id)](*e, params)
		}
	}

	return 404
}

// Add an option to the MediaPlayer Entity
func (e *MediaPlayerEntity) AddOption(option MediaPlayerEntityOption, value interface{}) {

	e.Options[option] = value

}
