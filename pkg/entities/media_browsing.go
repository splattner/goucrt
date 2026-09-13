package entities

// Types for the media_player entity's media browsing and searching support (the browse_media/
// search_media protocol messages). Defined here rather than in pkg/integration so that package can
// reference them directly in its request/response message types without a conversion layer -
// there's already precedent for that (EntityStateData is used the same way).
//
// See https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_media_player.md's
// "Media Browsing" and "Media Searching" sections.

// MediaClass is for browser/structure semantics: how a media item should be presented and
// organized in the media browser hierarchy. The named constants are the spec's recommended values
// (so the UI can show type-specific icons); an integration may also use its own custom value.
type MediaClass string

const (
	AlbumMediaClass     MediaClass = "album"
	AppMediaClass       MediaClass = "app"
	ArtistMediaClass    MediaClass = "artist"
	ChannelMediaClass   MediaClass = "channel"
	ComposerMediaClass  MediaClass = "composer"
	DirectoryMediaClass MediaClass = "directory"
	EpisodeMediaClass   MediaClass = "episode"
	GameMediaClass      MediaClass = "game"
	GenreMediaClass     MediaClass = "genre"
	ImageMediaClass     MediaClass = "image"
	MovieMediaClass     MediaClass = "movie"
	MusicMediaClass     MediaClass = "music"
	PlaylistMediaClass  MediaClass = "playlist"
	PodcastMediaClass   MediaClass = "podcast"
	RadioMediaClass     MediaClass = "radio"
	SeasonMediaClass    MediaClass = "season"
	TrackMediaClass     MediaClass = "track"
	TvShowMediaClass    MediaClass = "tv_show"
	UrlMediaClass       MediaClass = "url"
	VideoMediaClass     MediaClass = "video"
)

// MediaContentType is for playback/content semantics: the type of media content to play or that's
// currently playing. Distinct from MediaClass, which is about browser presentation. As with
// MediaClass, the named constants are recommended values; a custom value is also valid.
type MediaContentType string

const (
	AlbumMediaContentType    MediaContentType = "album"
	AppMediaContentType      MediaContentType = "app"
	AppsMediaContentType     MediaContentType = "apps"
	ArtistMediaContentType   MediaContentType = "artist"
	ChannelMediaContentType  MediaContentType = "channel"
	ChannelsMediaContentType MediaContentType = "channels"
	ComposerMediaContentType MediaContentType = "composer"
	EpisodeMediaContentType  MediaContentType = "episode"
	GameMediaContentType     MediaContentType = "game"
	GenreMediaContentType    MediaContentType = "genre"
	ImageMediaContentType    MediaContentType = "image"
	MovieMediaContentType    MediaContentType = "movie"
	MusicMediaContentType    MediaContentType = "music"
	PlaylistMediaContentType MediaContentType = "playlist"
	PodcastMediaContentType  MediaContentType = "podcast"
	RadioMediaContentType    MediaContentType = "radio"
	SeasonMediaContentType   MediaContentType = "season"
	TrackMediaContentType    MediaContentType = "track"
	TvShowMediaContentType   MediaContentType = "tv_show"
	UrlMediaContentType      MediaContentType = "url"
	VideoMediaContentType    MediaContentType = "video"
)

// MediaPaging is a browse_media/search_media request's optional paging parameters.
type MediaPaging struct {
	// Limit is the number of items per page. Spec default: 10.
	Limit int `json:"limit,omitempty"`
	// Page is the requested page number, 1-based. Spec default: 1.
	Page int `json:"page,omitempty"`
}

// MediaPagination is a browse_media/search_media response's actual paging result.
type MediaPagination struct {
	// Count is the total number of items, if known. Nil (omitted) if unknown.
	Count *int `json:"count,omitempty"`
	Limit int  `json:"limit,omitempty"`
	Page  int  `json:"page,omitempty"`
}

// MediaSearchFilter narrows a search_media request.
type MediaSearchFilter struct {
	MediaClasses []MediaClass `json:"media_classes,omitempty"`
	Artist       string       `json:"artist,omitempty"`
	Album        string       `json:"album,omitempty"`
}

// BrowseMediaItem is one entry in a browse_media/search_media response: either a container (
// CanBrowse, with up to one level of child Items) or a playable/searchable leaf. A search result
// item never has children, per the spec.
type BrowseMediaItem struct {
	// MediaId identifies the item; empty only for a non-playable item such as a root directory.
	MediaId    string           `json:"media_id"`
	Title      string           `json:"title"`
	Subtitle   string           `json:"subtitle,omitempty"`
	Artist     string           `json:"artist,omitempty"`
	Album      string           `json:"album,omitempty"`
	MediaClass MediaClass       `json:"media_class,omitempty"`
	MediaType  MediaContentType `json:"media_type,omitempty"`
	// CanBrowse: the item is a container, browsable via this MediaId/MediaType.
	CanBrowse bool `json:"can_browse,omitempty"`
	// CanPlay: the item can be played directly via the play_media command with this MediaId/MediaType.
	CanPlay bool `json:"can_play,omitempty"`
	// CanSearch: a search can be scoped to this item via search_media with this MediaId/MediaType.
	CanSearch bool `json:"can_search,omitempty"`
	// Thumbnail: a URL, a base64-encoded PNG/JPG, or "icon://uc:<name>" for a built-in icon.
	// Prefer a URL - keep payloads small. Preferred image size: 480x480.
	Thumbnail string `json:"thumbnail,omitempty"`
	// Duration in seconds.
	Duration int `json:"duration,omitempty"`
	// Items: child items if this is a container. Only one level of nesting - a child item's own
	// Items must be empty; browse again for deeper levels.
	Items []BrowseMediaItem `json:"items,omitempty"`
}

// BrowseMediaRequest is a browse_media request's parameters, passed to a MediaPlayerEntity's
// registered BrowseFunc.
type BrowseMediaRequest struct {
	// MediaId, MediaType: which container to browse. Both empty means browse the root.
	MediaId   string
	MediaType MediaContentType
	// StableIds is a hint that the caller needs stable identifiers across requests - some media
	// providers (e.g. Roon) generate new keys per request, in which case a driver may need to
	// return a resolvable path instead (see the spec's browse_media description for the rationale).
	StableIds bool
	Paging    *MediaPaging
}

// BrowseMediaResult is a MediaPlayerEntity's BrowseFunc's return value.
type BrowseMediaResult struct {
	// Media is the browsed container (its Items are the container's contents), or nil for special
	// cases like an empty root.
	Media      *BrowseMediaItem
	Pagination MediaPagination
}

// SearchMediaRequest is a search_media request's parameters, passed to a MediaPlayerEntity's
// registered SearchFunc.
type SearchMediaRequest struct {
	Query string
	// MediaId, MediaType: optionally scope the search to a specific container.
	MediaId   string
	MediaType MediaContentType
	StableIds bool
	Filter    *MediaSearchFilter
	Paging    *MediaPaging
}

// SearchMediaResult is a MediaPlayerEntity's SearchFunc's return value.
type SearchMediaResult struct {
	Media      []BrowseMediaItem
	Pagination MediaPagination
}

// SetBrowseFunc registers the function called to handle a browse_media request for this entity.
// Requires the BrowseMediaMediaPlayerEntityFeatures feature to be declared.
func (e *MediaPlayerEntity) SetBrowseFunc(f func(BrowseMediaRequest) (*BrowseMediaResult, error)) {
	e.BrowseFunc = f
}

// SetSearchFunc registers the function called to handle a search_media request for this entity.
// Requires the SearchMediaMediaPlayerEntityFeatures feature to be declared.
func (e *MediaPlayerEntity) SetSearchFunc(f func(SearchMediaRequest) (*SearchMediaResult, error)) {
	e.SearchFunc = f
}
