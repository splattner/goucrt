package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// BindStandardFlags registers the persistent flags every goucrt-based driver binary needs -
// listen port, websocket path, mDNS/registration toggles, config/data home, and so on - on cmd,
// and binds each to viper together with its corresponding UC_* environment variable. Call this
// once on your root command before Execute(); integration.NewIntegration then reads the result
// back via the Config it's given (typically built with viper.Unmarshal).
//
// goucrt's own multi-driver ucrt binary calls this, and every single-driver binary (in this repo
// or an external one) should too, instead of hand-copying the registration: it's how the same
// ~15 flags stay identical, with identical env var names and defaults, across every driver.
func BindStandardFlags(cmd *cobra.Command) error {
	flags := cmd.PersistentFlags()

	flags.IntP("listenPort", "l", 8080, "the port this integration is listening for websocket connection from the remote")
	flags.String("websocketPath", "/ws", "path where this integration is available for websocket connections")
	flags.Bool("disableMDNS", false, "Disable integration advertisement via mDNS")
	flags.String("remoteTwoIP", "", "IP Address of your Remote Two instance (disables Remote Two discovery)")
	flags.Int("remoteTwoPort", 80, "Port of your Remote Two instance (disables Remote Two discovery)")
	flags.Bool("registration", false, "Enable driver registration on the Remote Two instead of mDNS advertisement")
	flags.String("registrationUsername", "web-configurator", "Username of the RemoteTwo for driver registration")
	flags.String("registrationPin", "", "Pin of the RemoteTwo for driver registration")
	flags.Bool("debug", false, "Enable debug log level")
	flags.String("ucconfighome", "./ucconfig/", "Configuration directory to save the user configuration from the driver setup")
	flags.String("datahome", "", "Directory for a driver's own application data; defaults to ucconfighome if unset")
	flags.String("bindInterface", "", "Host/IP address to listen on (default: all interfaces). Overridden by UC_INTEGRATION_INTERFACE when set.")

	// env is empty for flags with no corresponding UC_* environment variable: "debug" has none,
	// and "bindInterface" is deliberately not viper-bound to UC_INTEGRATION_INTERFACE here since
	// NewIntegration reads that env var itself and must win even when running as a
	// custom-installed driver, which can't pass CLI flags at all.
	bindings := []struct{ name, env string }{
		{"listenPort", "UC_INTEGRATION_LISTEN_PORT"},
		{"websocketPath", "UC_INTEGRATION_WEBSOCKET_PATH"},
		{"disableMDNS", "UC_DISABLE_MDNS_PUBLISH"},
		{"remoteTwoIP", "UC_RT_HOST"},
		{"remoteTwoPort", "UC_RT_PORT"},
		{"registration", "UC_ENABLE_REGISTRATION"},
		{"registrationUsername", "UC_REGISTRATION_USERNAME"},
		{"registrationPin", "UC_REGISTRATION_PIN"},
		{"debug", ""},
		{"ucconfighome", "UC_CONFIG_HOME"},
		{"datahome", "UC_DATA_HOME"},
		{"bindInterface", ""},
	}

	var errs []error
	for _, b := range bindings {
		if err := viper.BindPFlag(b.name, flags.Lookup(b.name)); err != nil {
			errs = append(errs, fmt.Errorf("bind flag %q: %w", b.name, err))
		}
		if b.env != "" {
			if err := viper.BindEnv(b.name, b.env); err != nil {
				errs = append(errs, fmt.Errorf("bind env for %q: %w", b.name, err))
			}
		}
	}

	return errors.Join(errs...)
}
