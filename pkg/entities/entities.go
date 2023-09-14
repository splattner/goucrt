package entities

type EntityState string

const (
	UnavailableEntityState EntityState = "UNAVAILABLE"
	UnkownEntityState                  = "UNKNOWN"
)

type Entity struct {
	Id string `json:"entity_id"`
	EntityType
	DeviceId    string                 `json:"device_id"`
	Features    []interface{}          `json:"features"`
	Name        LanguageText           `json:"name"`
	Area        string                 `json:"area"`
	DeviceClass string                 `json:"-"`
	Attributes  map[string]interface{} `json:"-"`
}

type EntityType struct {
	Type string `json:"entity_type"`
}

type EntityStateData struct {
	DeviceId string `json:"device_id,omitempty"`
	EntityType
	EntityId   string                 `json:"entity_id"`
	Attributes map[string]interface{} `json:"attributes"`
}

// Add an attribute if not already available
func (e *Entity) AddAttribute(name string, value interface{}) {
	if e.Attributes[name] != nil {
		e.Attributes[name] = value
	}
}

func (e *Entity) GetEntityState() *EntityStateData {

	entityState := EntityStateData{
		DeviceId:   e.DeviceId,
		EntityType: e.EntityType,
		EntityId:   e.Id,
		Attributes: e.Attributes,
	}

	return &entityState
}
