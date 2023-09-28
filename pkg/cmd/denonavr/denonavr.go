package denonavr

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/splattner/goucrt/pkg/client"
	"github.com/splattner/goucrt/pkg/cmd"
	"github.com/splattner/goucrt/pkg/integration"
)

func NewCommand(rootCmd *cobra.Command) *cobra.Command {

	var command = &cobra.Command{
		Use:   "denonavr",
		Short: "Start Denon AVR Ingegration",
		Long:  "Denon AVR Integration for a Unfolded Circle Remote Two",
		Run: func(c *cobra.Command, args []string) {

			log.SetOutput(os.Stdout)

			debug := viper.GetBool("debug")
			if debug {
				log.SetLevel(log.DebugLevel)
			} else {
				log.SetLevel(log.InfoLevel)
			}

			var config integration.Config
			viper.Unmarshal(&config)

			i, err := integration.NewIntegration(config)
			cmd.CheckError(err)

			myclient := client.NewDenonAVRClient(i)

			myclient.InitClient()

			cmd.CheckError(i.Run())

		},
	}

	return command
}
