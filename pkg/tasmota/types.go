package tasmota

type TasmotaResultMsg struct {
	Power1     string `json:"POWER1,omitempty"`
	Power      string `json:"POWER,omitempty"`
	Dimmer     int    `json:"Dimmer,omitempty"`
	Color      string `json:"Color,omitempty"`
	HSBCOlor   string `json:"HSBColor,omitempty"`
	White      int    `json:"White,omitempty"`
	Channel    []int  `json:"Channel,omitempty"`
	CustomSend string `json:"CustomSend,omitempty"`
}

type TasmotaTeleMsg struct {
	Time     string              `json:"time,omitempty"`
	TempUnit string              `json:"TempUnit,omitempty"`
	SI7021   TasmotaTeleSI721Msg `json:"SI7021,omitempty"`
}

type TasmotaTeleSI721Msg struct {
	Temperature float32 `json:"Temperature,omitempty"`
	Humidity    float32 `json:"Humidity,omitempty"`
	DewPoint    float32 `json:"DewPoint,omitempty"`
}
