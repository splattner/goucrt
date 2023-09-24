package deconz

type DeconzWebSocketMessage struct {
	Type       string               `json:"t,omitempty"`
	Event      string               `json:"e,omitempty"`
	Resource   string               `json:"r,omitempty"`
	ID         string               `json:"id,omitempty"`
	UniqueID   string               `json:"uniqueid,omitempty"`
	GroupID    string               `json:"gid,omitempty"`
	SceneID    string               `json:"scid,omitempty"`
	Name       string               `json:"name,omitempty"`
	Attributes DeconzLightAttribute `json:"attr,omitempty"`
	State      DeconzState          `json:"state,omitempty"`
}

type DeconzLightAttribute struct {
	Id                string `json:"id,omitempty"`
	LastAnnounced     string `json:"lastannounced,omitempty"`
	LastSeen          string `json:"lastseen,omitempty"`
	ManufacturerName  string `json:"manufacturername,omitempty"`
	ModelId           string `json:"modelid,omitempty"`
	Name              string `json:"name,omitempty"`
	SWVersion         string `json:"swversion,omitempty"`
	Type              string `json:"type,omitempty"`
	UniqueID          string `json:"uniqueid,omitempty"`
	ColorCapabilities int    `json:"colorcapabilities,omitempty"`
	Ctmax             int    `json:"ctmax,omitempty"`
	Ctmin             int    `json:"ctmin,omitempty"`
}

type DeconzState struct {

	// Light & Group
	On     *bool     `json:"on,omitempty"`     //
	Hue    *uint16   `json:"hue,omitempty"`    //
	Effect string    `json:"effect,omitempty"` //
	Bri    *uint8    `json:"bri,omitempty"`    // min = 1, max = 254
	Sat    *uint8    `json:"sat,omitempty"`    //
	CT     *uint16   `json:"ct,omitempty"`     // min = 154, max = 500
	XY     []float32 `json:"xy,omitempty"`
	Alert  string    `json:"alert,omitempty"`

	// Light
	Reachable      *bool   `json:"reachable,omitempty"`
	ColorMode      string  `json:"colormode,omitempty"`
	ColorLoopSpeed *uint8  `json:"colorloopspeed,omitempty"`
	TransitionTime *uint16 `json:"transitiontime,omitempty"`

	// Group
	AllOn bool `json:"all_on,omitempty"`
	AnyOn bool `json:"any_on,omitempty"`

	// Sensor
	ButtonEvent int `json:"buttonevent,omitempty"`
}
