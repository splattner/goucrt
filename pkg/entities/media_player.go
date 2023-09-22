package entities

type MediaPlayerEntityState string
type MediaPlayerEntityFeatures string
type MediaPlayerEntityAttributes string
type MediaPlayerEntityCommand string
type MediaPlayerDeviceClass string

const (
	OnMediaPlayerEntityState        MediaPlayerEntityState = "ON"
	OffMediaPlayerEntityState                              = "OFF"
	PlayingPlayerEntityState                               = "PLAYING"
	PausedMediaPlayerEntityState                           = "PAUSED"
	StandbyMediaPlayerEntityState                          = "STANDBY"
	BufferingMediaPlayerEntityState                        = "BUFFERING"
)

const (
	OnOffMediaPlayerEntityFeatures           MediaPlayerEntityFeatures = "on_off"
	ToggleMediaPlayerEntityyFeatures                                   = "toggle"
	VolumeMediaPlayerEntityyFeatures                                   = "volume"
	VolumeUpDownMediaPlayerEntityFeatures                              = "volume_up_down"
	MuteToggleMediaPlayerEntityFeatures                                = "mute_toggle"
	MuteMediaPlayerEntityFeatures                                      = "mute"
	UnmuteMediaPlayerEntityFeatures                                    = "unmtue"
	PlayPauseMediaPlayerEntityFeatures                                 = "play_pause"
	StopMediaPlayerEntityFeatures                                      = "stop"
	NextMediaPlayerEntityFeatures                                      = "next"
	PreviusMediaPlayerEntityFeatures                                   = "previous"
	FastForwardMediaPlayerEntityFeatures                               = "fast_forward"
	RewindMediaPlayerEntityFeatures                                    = "rewind"
	RepeatMediaPlayerEntityFeatures                                    = "repeat"
	ShuffleMediaPlayerEntityFeatures                                   = "shuffle"
	SeekMediaPlayerEntityFeatures                                      = "seek"
	MediaDurationMediaPlayerEntityFeatures                             = "media_duration"
	MediaPositionMediaPlayerEntityFeatures                             = "media_position"
	MediaTitleMediaPlayerEntityFeatures                                = "media_title"
	MediaArtistMediaPlayerEntityFeatures                               = "media_artist"
	MediaAlbumMediaPlayerEntityFeatures                                = "media_album"
	MediaImageUrlMediaPlayerEntityFeatures                             = "media_image_url"
	MediaTypeMediaPlayerEntityFeatures                                 = "media_type"
	DPadMediaPlayerEntityFeatures                                      = "dpad"
	HomeMediaPlayerEntityFeatures                                      = "home"
	MenuMediaPlayerEntityFeatures                                      = "menu"
	ColorButtonsMediaPlayerEntityFeatures                              = "color_buttons"
	ChannelSwitcherMediaPlayerEntityFeatures                           = "channel_switcher"
	SelectSourceMediaPlayerEntityFeatures                              = "select_source"
	SelectSoundModeMediaPlayerEntityFeatures                           = "select_sound_mode"
)

const (
	StateMediaPlayerEntityAttribute         MediaPlayerEntityAttributes = "state"
	VolumeMediaPlayerEntityAttribute                                    = "volume"
	MutedMediaPlayeEntityAttribute                                      = "muted"
	MediaDurationMediaPlayerEntityAttribute                             = "media_duration"
	MediaPositionMediaPlayerEntityAttribute                             = "media_position"
	MediaTypeMediaPlayerEntityAttribute                                 = "media_type"
	MediaImageUrlMediaPlayerEntityAttribute                             = "media_image_url"
	MediaTitleMediaPlayerEntityAttribute                                = "media_title"
	MediaArtistMediaPlayerEntityAttribute                               = "media_artist"
	MediaAlbumMediaPlayerEntityAttribute                                = "media_album"
	RepeatMediaPlayerEntityAttribute                                    = "repeat"
	ShuffleMediaPlayerEntityAttribute                                   = "shuffle"
	SourceMediaPlayerEntityAttribute                                    = "source"
	SourceListMediaPlayerEntityAttribute                                = "source_list"
	SoundModeMediaPlayerEntityAttribute                                 = "sound_mode"
	SoundModeListMediaPlayerEntityAttribute                             = "sound_mode_list"
)

const (
	OnMediaPlayerEntityCommand              MediaPlayerEntityCommand = "on"
	OffMediaPlayerEntityCommand                                      = "off"
	ToggleMediaPlayerEntityCommand                                   = "toggle"
	PlayPauseMediaPlayerEntityCommand                                = "play_pause"
	StopMediaPlayerEntityCommand                                     = "stop"
	PreviusMediaPlayerEntityCommand                                  = "previous"
	NextMediaPlayerEntityCommand                                     = "next"
	FastForwardMediaPlayerEntityCommand                              = "fast_forward"
	RewindMediaPlayerEntityCommand                                   = "rewind"
	SeekMediaPlayerEntityCommand                                     = "seek"
	VolumeMediaPlayerEntityCommand                                   = "volume"
	VolumeUpMediaPlayerEntityCommand                                 = "volume_up"
	VolumeDownMediaPlayerEntityCommand                               = "volume_down"
	MuteToggleMediaPlayerEntityCommand                               = "mute_toggle"
	MuteMediaPlayerEntityCommand                                     = "mute"
	UnmuteMediaPlayerEntityCommand                                   = "unmute"
	RepeatMediaPlayerEntityCommand                                   = "repeat"
	ShuffleMediaPlayerEntityCommand                                  = "shuffle"
	ChannelUpMediaPlayerEntityCommand                                = "channel_up"
	ChannelDownMediaPlayerEntityCommand                              = "channel_down"
	CursorUpMediaPlayerEntityCommand                                 = "cursor_up"
	CursorDownMediaPlayerEntityCommand                               = "cursor_down"
	CursorLeftMediaPlayerEntityCommand                               = "cursor_left"
	CursorRightMediaPlayerEntityCommand                              = "cursor_right"
	CursorEnterMediaPlayerEntityCommand                              = "cursor_enter"
	FunctionRedMediaPlayerEntityCommand                              = "function_red"
	FunctionGreenMediaPlayerEntityCommand                            = "function_green"
	FunctionYellowMediaPlayerEntityCommand                           = "function_yellow"
	FunctionBlueMediaPlayerEntityCommand                             = "function_blue"
	HomeMediaPlayerEntityCommand                                     = "home"
	MenuMediaPlayerEntityCommand                                     = "menu"
	BackMediaPlayerEntityCommand                                     = "back"
	SelectSourcMediaPlayerEntityCommand                              = "select_source"
	SelectSoundModeMediaPlayerEntityCommand                          = "select_sound_mode"
	SearchMediaPlayerEntityCommand                                   = "search"
)

const (
	ReceiverMediaPlayerDeviceClass      MediaPlayerDeviceClass = "receiver"
	SetTopBoxMediaPlayerDeviceClass                            = "set_top_box"
	SpeakerMediaPlayerDeviceClass                              = "speaker"
	StreamingBoxMMediaPlayerDeviceClass                        = "streaming_box"
	TVMediaPlayerDeviceClass                                   = "tv"
)

type MediaPlayerEntity struct {
	Entity
	DeviceClass MediaPlayerDeviceClass
	Commands    map[string]func(MediaPlayerEntity, map[string]interface{}) int `json:"-"`
}

func NewMediaPlayerEntity(id string, name LanguageText, area string, deviceClass MediaPlayerDeviceClass) *MediaPlayerEntity {

	mediaPlayerEntity := MediaPlayerEntity{}
	mediaPlayerEntity.Id = id
	mediaPlayerEntity.Name = name
	mediaPlayerEntity.Area = area
	mediaPlayerEntity.DeviceClass = deviceClass

	mediaPlayerEntity.EntityType.Type = "media_player"

	mediaPlayerEntity.Commands = make(map[string]func(MediaPlayerEntity, map[string]interface{}) int)
	mediaPlayerEntity.Attributes = make(map[string]interface{})

	return &mediaPlayerEntity
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

	case ToggleMediaPlayerEntityyFeatures:
		e.AddAttribute(string(StateMediaPlayerEntityAttribute), OffMediaPlayerEntityState)

	case PlayPauseMediaPlayerEntityFeatures:
		e.AddAttribute(string(StateMediaPlayerEntityAttribute), OffMediaPlayerEntityState)
		e.AddAttribute(string(MediaPositionMediaPlayerEntityAttribute), 0)
		e.AddAttribute(string(MediaImageUrlMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(MediaTitleMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(MediaArtistMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(MediaAlbumMediaPlayerEntityAttribute), "")

	case StopMediaPlayerEntityFeatures:
		e.AddAttribute(string(StateMediaPlayerEntityAttribute), OffMediaPlayerEntityState)
		e.AddAttribute(string(MediaPositionMediaPlayerEntityAttribute), 0)

	case VolumeMediaPlayerEntityyFeatures:
		e.AddAttribute(string(VolumeMediaPlayerEntityAttribute), 0)

	case VolumeUpDownMediaPlayerEntityFeatures:
		e.AddAttribute(string(VolumeMediaPlayerEntityAttribute), 0)

	case MuteToggleMediaPlayerEntityFeatures:
		e.AddAttribute(string(MutedMediaPlayeEntityAttribute), false)

	case MuteMediaPlayerEntityFeatures:
		e.AddAttribute(string(MutedMediaPlayeEntityAttribute), false)

	case UnmuteMediaPlayerEntityFeatures:
		e.AddAttribute(string(MutedMediaPlayeEntityAttribute), false)

	case MediaDurationMediaPlayerEntityFeatures:
		e.AddAttribute(string(MediaDurationMediaPlayerEntityAttribute), 0)

	case MediaPositionMediaPlayerEntityFeatures:
		e.AddAttribute(string(MediaPositionMediaPlayerEntityAttribute), 0)

	case FastForwardMediaPlayerEntityFeatures:
		e.AddAttribute(string(MediaPositionMediaPlayerEntityAttribute), 0)

	case RewindMediaPlayerEntityFeatures:
		e.AddAttribute(string(MediaPositionMediaPlayerEntityAttribute), 0)

	case SeekMediaPlayerEntityFeatures:
		e.AddAttribute(string(MediaPositionMediaPlayerEntityAttribute), 0)

	case MediaTypeMediaPlayerEntityFeatures:
		e.AddAttribute(string(MediaTypeMediaPlayerEntityAttribute), 0)

	case MediaImageUrlMediaPlayerEntityFeatures:
		e.AddAttribute(string(MediaImageUrlMediaPlayerEntityAttribute), "")

	case NextMediaPlayerEntityFeatures:
		e.AddAttribute(string(MediaImageUrlMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(MediaTitleMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(MediaArtistMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(MediaAlbumMediaPlayerEntityAttribute), "")

	case PreviusMediaPlayerEntityFeatures:
		e.AddAttribute(string(MediaImageUrlMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(MediaTitleMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(MediaArtistMediaPlayerEntityAttribute), "")
		e.AddAttribute(string(MediaAlbumMediaPlayerEntityAttribute), "")

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

	}
}

// Register a function for the Entity command
func (e *MediaPlayerEntity) AddCommand(command MediaPlayerEntityCommand, function func(MediaPlayerEntity, map[string]interface{}) int) {
	e.Commands[string(command)] = function

}

// Call the registred function for this entity_command
func (e *MediaPlayerEntity) HandleCommand(cmd_id string, params map[string]interface{}) int {
	if e.Commands[cmd_id] != nil {
		return e.Commands[cmd_id](*e, params)
	}

	return 404
}
