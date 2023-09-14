package entities

type EntityName struct {
	En string `json:"en"`
	De string `json:"de"`
}

type Entity struct {
	Id string `json:"entity_id"`
	EntityType
	DeviceId string     `json:"device_id"`
	Features []string   `json:"features"`
	Name     EntityName `json:"name"`
	Area     string     `json:"string"`
}

type EntityType struct {
	Type string `json:"entity_type"`
}

type EntityStateData struct {
	DeviceId string `json:"device_id,omitempty"`
	EntityType
	EntityId   string      `json:"entity_id"`
	Attributes interface{} `json:"attributes"`
}
