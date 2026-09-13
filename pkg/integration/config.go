package integration

// Generic string key/value config map to store configuration option
type Config struct {
	ListenPort int `mapstructure:"listenPort"`
	// BindInterface is the host/IP address to listen on. Empty means all interfaces (0.0.0.0).
	// Set by the remote itself via UC_INTEGRATION_INTERFACE when running as a custom-installed
	// driver (see doc/integration-driver/driver-installation.md's "Runtime environment" section).
	BindInterface        string `mapstructure:"bindInterface"`
	DisableMDNS          bool   `mapstructure:"disableMDNS"`
	EnableRegistration   bool   `mapstructure:"enableRegistration"`
	RegistrationUsername string `mapstructure:"registrationUsername"`
	RegistrationPin      string `mapstructure:"registrationPin"`
	WebsocketPath        string `mapstructure:"websocketPath"`
	ConfigHome           string `mapstructure:"ucconfighome"`
	// DataHome is where a driver can persist its own application data (device caches, etc.),
	// distinct from ConfigHome's user-entered setup data. goucrt itself doesn't write anything
	// here yet; it's read from UC_DATA_HOME and exposed for driver code to use. See
	// doc/integration-driver/driver-installation.md.
	DataHome                 string `mapstructure:"datahome"`
	RemoteTwoHost            string `mapstructure:"remoteTwoIP"`
	RemoteTwoPort            int    `mapstructure:"remoteTwoPort"`
	IgnoreEntitySubscription bool
}
