package ucrt

import (
	"github.com/spf13/cobra"
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
	rootCmd.PersistentFlags().String("websocketPath", "/ws", "path where this integration is available for websocket connections")

	rootCmd.PersistentFlags().Bool("mdns", true, "Enable integration advertisement via mDNS")
	rootCmd.PersistentFlags().Bool("registration", false, "Enable driver registration on the Remote Two instead of mDNS advertisement")
	rootCmd.PersistentFlags().String("registrationUsername", "web-configurator", "Username of the RemoteTwo for driver registration")
	rootCmd.PersistentFlags().String("registrationPin", "", "Pin of the RemoteTwo for driver registration")
	rootCmd.PersistentFlags().String("remoteTwoIP", "", "IP Address of your Remote Two instance (disables Remote Two discovery)")
	rootCmd.PersistentFlags().Int("remoteTwoPort", 80, "Port of your Remote Two instance (disables Remote Two discovery)")

	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug log level")

	rootCmd.AddCommand(
		denonavr.NewCommand(rootCmd),
		deconz.NewCommand(rootCmd),
		shelly.NewCommand(rootCmd),
	)

	return rootCmd
}
