package ucrt

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/splattner/goucrt/pkg/cmd/deconz"
	"github.com/splattner/goucrt/pkg/cmd/denonavr"
	"github.com/splattner/goucrt/pkg/cmd/shelly"
)

func NewCommand(name string) *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   name,
		Short: "Unfolder Circle Remote Two integration",
		Long:  `Unfolder Circle Remote Two integration`,
	}

	rootCmd.PersistentFlags().IntP("listenPort", "l", 8080, "the port this integration is listening for websocket connection from the remote")
	viper.BindPFlag("listenPort", rootCmd.PersistentFlags().Lookup("listenPort"))
	viper.BindEnv("listenPort", "UC_INTEGRATION_LISTEN_PORT")

	rootCmd.PersistentFlags().String("websocketPath", "/ws", "path where this integration is available for websocket connections")
	viper.BindPFlag("websocketPath", rootCmd.PersistentFlags().Lookup("websocketPath"))
	viper.BindEnv("websocketPath", "UC_INTEGRATION_WEBSOCKET_PATH")

	rootCmd.PersistentFlags().Bool("disableMDNS", false, "Disable integration advertisement via mDNS")
	viper.BindPFlag("disableMDNS", rootCmd.PersistentFlags().Lookup("disableMDNS"))
	viper.BindEnv("disableMDNS", "UC_DISABLE_MDNS_PUBLISH")

	rootCmd.PersistentFlags().String("remoteTwoIP", "", "IP Address of your Remote Two instance (disables Remote Two discovery)")
	viper.BindPFlag("remoteTwoIP", rootCmd.PersistentFlags().Lookup("remoteTwoIP"))
	viper.BindEnv("remoteTwoIP", "UC_RT_HOST")

	rootCmd.PersistentFlags().Int("remoteTwoPort", 80, "Port of your Remote Two instance (disables Remote Two discovery)")
	viper.BindPFlag("remoteTwoPort", rootCmd.PersistentFlags().Lookup("remoteTwoPort"))
	viper.BindEnv("remoteTwoPort", "UC_RT_PORT")

	rootCmd.PersistentFlags().Bool("registration", false, "Enable driver registration on the Remote Two instead of mDNS advertisement")
	viper.BindPFlag("registration", rootCmd.PersistentFlags().Lookup("registration"))
	viper.BindEnv("registration", "UC_ENABLE_REGISTRATION")

	rootCmd.PersistentFlags().String("registrationUsername", "web-configurator", "Username of the RemoteTwo for driver registration")
	viper.BindPFlag("registrationUsername", rootCmd.PersistentFlags().Lookup("registrationUsername"))
	viper.BindEnv("registrationUsername", "UC_REGISTRATION_USERNAME")

	rootCmd.PersistentFlags().String("registrationPin", "", "Pin of the RemoteTwo for driver registration")
	viper.BindPFlag("registrationPin", rootCmd.PersistentFlags().Lookup("registrationPin"))
	viper.BindEnv("registrationUsername", "UC_REGISTRATION_PIN")

	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug log level")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.PersistentFlags().String("ucconfighome", "./ucconfig/", "Configuration directory to save the user configuration from the driver setup")
	viper.BindPFlag("ucconfighome", rootCmd.PersistentFlags().Lookup("ucconfighome"))
	viper.BindEnv("ucconfighome", "UC_CONFIG_HOME")

	rootCmd.AddCommand(
		denonavr.NewCommand(rootCmd),
		deconz.NewCommand(rootCmd),
		shelly.NewCommand(rootCmd),
	)

	return rootCmd
}
