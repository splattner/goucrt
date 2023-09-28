package integration

// Generic string key/value config map to store configuration option
type Config struct {
	ListenPort               int    `mapstructure:"listenPort"`
	DisableMDNS              bool   `mapstructure:"disableMDNS"`
	EnableRegistration       bool   `mapstructure:"enableRegistration"`
	RegistrationUsername     string `mapstructure:"registrationUsername"`
	RegistrationPin          string `mapstructure:"registrationPin"`
	WebsocketPath            string `mapstructure:"websocketPath"`
	ConfigHome               string `mapstructure:"ucconfighome"`
	RemoteTwoHost            string `mapstructure:"remoteTwoIP"`
	RemoteTwoPort            int    `mapstructure:"remoteTwoPort"`
	IgnoreEntitySubscription bool
}
