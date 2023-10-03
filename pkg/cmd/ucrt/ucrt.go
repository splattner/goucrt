package ucrt

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/splattner/goucrt/pkg/cmd/deconz"
	"github.com/splattner/goucrt/pkg/cmd/denonavr"
	"github.com/splattner/goucrt/pkg/cmd/shelly"
	"github.com/splattner/goucrt/pkg/cmd/tasmota"

	log "github.com/sirupsen/logrus"
)

func NewCommand(name string) *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   name,
		Short: "Unfolder Circle Remote Two integration",
		Long:  `Unfolder Circle Remote Two integration`,
	}

	rootCmd.PersistentFlags().IntP("listenPort", "l", 8080, "the port this integration is listening for websocket connection from the remote")
	if err := viper.BindPFlag("listenPort", rootCmd.PersistentFlags().Lookup("listenPort")); err != nil {
		log.WithError(err).Error(("Cannot bindPFplag"))
	}
	if err := viper.BindEnv("listenPort", "UC_INTEGRATION_LISTEN_PORT"); err != nil {
		log.WithError(err).Error(("Cannot BindEnv"))
	}

	rootCmd.PersistentFlags().String("websocketPath", "/ws", "path where this integration is available for websocket connections")
	if err := viper.BindPFlag("websocketPath", rootCmd.PersistentFlags().Lookup("websocketPath")); err != nil {
		log.WithError(err).Error(("Cannot bindPFplag"))
	}
	if err := viper.BindEnv("websocketPath", "UC_INTEGRATION_WEBSOCKET_PATH"); err != nil {
		log.WithError(err).Error(("Cannot BindEnv"))
	}

	rootCmd.PersistentFlags().Bool("disableMDNS", false, "Disable integration advertisement via mDNS")
	if err := viper.BindPFlag("disableMDNS", rootCmd.PersistentFlags().Lookup("disableMDNS")); err != nil {
		log.WithError(err).Error(("Cannot bindPFplag"))
	}
	if err := viper.BindEnv("disableMDNS", "UC_DISABLE_MDNS_PUBLISH"); err != nil {
		log.WithError(err).Error(("Cannot BindEnv"))
	}

	rootCmd.PersistentFlags().String("remoteTwoIP", "", "IP Address of your Remote Two instance (disables Remote Two discovery)")
	if err := viper.BindPFlag("remoteTwoIP", rootCmd.PersistentFlags().Lookup("remoteTwoIP")); err != nil {
		log.WithError(err).Error(("Cannot bindPFplag"))
	}
	if err := viper.BindEnv("remoteTwoIP", "UC_RT_HOST"); err != nil {
		log.WithError(err).Error(("Cannot BindEnv"))
	}

	rootCmd.PersistentFlags().Int("remoteTwoPort", 80, "Port of your Remote Two instance (disables Remote Two discovery)")
	if err := viper.BindPFlag("remoteTwoPort", rootCmd.PersistentFlags().Lookup("remoteTwoPort")); err != nil {
		log.WithError(err).Error(("Cannot bindPFplag"))
	}
	if err := viper.BindEnv("remoteTwoPort", "UC_RT_PORT"); err != nil {
		log.WithError(err).Error(("Cannot BindEnv"))
	}

	rootCmd.PersistentFlags().Bool("registration", false, "Enable driver registration on the Remote Two instead of mDNS advertisement")
	if err := viper.BindPFlag("registration", rootCmd.PersistentFlags().Lookup("registration")); err != nil {
		log.WithError(err).Error(("Cannot bindPFplag"))
	}
	if err := viper.BindEnv("registration", "UC_ENABLE_REGISTRATION"); err != nil {
		log.WithError(err).Error(("Cannot BindEnv"))
	}

	rootCmd.PersistentFlags().String("registrationUsername", "web-configurator", "Username of the RemoteTwo for driver registration")
	if err := viper.BindPFlag("registrationUsername", rootCmd.PersistentFlags().Lookup("registrationUsername")); err != nil {
		log.WithError(err).Error(("Cannot bindPFplag"))
	}
	if err := viper.BindEnv("registrationUsername", "UC_REGISTRATION_USERNAME"); err != nil {
		log.WithError(err).Error(("Cannot BindEnv"))
	}

	rootCmd.PersistentFlags().String("registrationPin", "", "Pin of the RemoteTwo for driver registration")
	if err := viper.BindPFlag("registrationPin", rootCmd.PersistentFlags().Lookup("registrationPin")); err != nil {
		log.WithError(err).Error(("Cannot bindPFplag"))
	}
	if err := viper.BindEnv("registrationUsername", "UC_REGISTRATION_PIN"); err != nil {
		log.WithError(err).Error(("Cannot BindEnv"))
	}

	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug log level")
	if err := viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")); err != nil {
		log.WithError(err).Error(("Cannot bindPFplag"))
	}

	rootCmd.PersistentFlags().String("ucconfighome", "./ucconfig/", "Configuration directory to save the user configuration from the driver setup")
	if err := viper.BindPFlag("ucconfighome", rootCmd.PersistentFlags().Lookup("ucconfighome")); err != nil {
		log.WithError(err).Error(("Cannot bindPFplag"))
	}
	if err := viper.BindEnv("ucconfighome", "UC_CONFIG_HOME"); err != nil {
		log.WithError(err).Error(("Cannot BindEnv"))
	}

	rootCmd.AddCommand(
		denonavr.NewCommand(rootCmd),
		deconz.NewCommand(rootCmd),
		shelly.NewCommand(rootCmd),
		tasmota.NewCommand(rootCmd),
	)

	return rootCmd
}
