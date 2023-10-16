package deconz

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	deconzclient "github.com/splattner/goucrt/pkg/clients/deconz"
	"github.com/splattner/goucrt/pkg/cmd"
	"github.com/splattner/goucrt/pkg/integration"
)

func NewCommand(rootCmd *cobra.Command) *cobra.Command {

	var command = &cobra.Command{
		Use:   "deconz",
		Short: "Start Deconz Ingegration",
		Long:  "Deconz Integration for a Unfolded Circle Remote Two",
		Run: func(c *cobra.Command, args []string) {

			log.SetOutput(os.Stdout)

			debug := viper.GetBool("debug")
			if debug {
				log.SetLevel(log.DebugLevel)
			} else {
				log.SetLevel(log.InfoLevel)
			}

			var config integration.Config
			if err := viper.Unmarshal(&config); err != nil {
				log.WithError(err).Error("Cannot unmarshal config with viper")
			}

			i, err := integration.NewIntegration(config)
			cmd.CheckError(err)

			myclient := deconzclient.NewDeconzClient(i)

			myclient.InitClient()

			cmd.CheckError(i.Run())

		},
	}

	return command
}
