package denonavr

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/splattner/goucrt/pkg/client"
	"github.com/splattner/goucrt/pkg/cmd"
	"github.com/splattner/goucrt/pkg/integration"
)

func NewCommand(rootCmd *cobra.Command) *cobra.Command {

	var command = &cobra.Command{
		Use:   "denonavr",
		Short: "Denon AVR",
		Long:  "Denon AVR Integration for a Unfolded Circle Remote Two",
		Run: func(c *cobra.Command, args []string) {

			log.SetOutput(os.Stdout)

			debug, _ := rootCmd.Flags().GetBool("debug")
			if debug {
				log.SetLevel(log.DebugLevel)
			} else {
				log.SetLevel(log.InfoLevel)
			}

			var config = make(integration.Config)

			listenPort, _ := rootCmd.Flags().GetInt("listenPort")
			enableMDNS, _ := rootCmd.Flags().GetBool("mdns")
			enableRegistration, _ := rootCmd.Flags().GetBool("registration")
			registrationUsername, _ := rootCmd.Flags().GetString("registrationUsername")
			registrationPin, _ := rootCmd.Flags().GetString("registrationPin")
			websocketPath, _ := rootCmd.Flags().GetString("websocketPath")
			remoteTwoIP, _ := rootCmd.Flags().GetString("remoteTwoIP")
			remoteTwoPort, _ := rootCmd.Flags().GetInt("remoteTwoPort")

			config["listenport"] = listenPort
			config["enableMDNS"] = enableMDNS
			config["enableRegistration"] = enableRegistration
			config["registrationUsername"] = registrationUsername
			config["registrationPin"] = registrationPin
			config["remoteTwoIP"] = remoteTwoIP
			config["remoteTwoPort"] = remoteTwoPort
			config["websocketPath"] = websocketPath

			i, err := integration.NewIntegration(config)
			cmd.CheckError(err)

			myclient := client.NewDenonAVRClient(i)

			myclient.InitClient()

			cmd.CheckError(i.Run())

		},
	}

	return command
}
