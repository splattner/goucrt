package server

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/splattner/goucrt/pkg/client"
	"github.com/splattner/goucrt/pkg/cmd"
	"github.com/splattner/goucrt/pkg/integration"
)

var (
	listenPort int
)

func NewCommand() *cobra.Command {

	var command = &cobra.Command{
		Use:   "server",
		Short: "Run a Unfolded Circle Remote Two integratin",
		Long:  "Run a Unfolded Circle Remote Two integratin",
		Run: func(c *cobra.Command, args []string) {
			log.SetOutput(os.Stdout)
			log.Println("Integration run")

			var config = make(integration.Config)

			config["listenport"] = listenPort

			i, err := integration.NewIntegration(config)
			cmd.CheckError(err)

			myclient := client.NewClient(i)

			myclient.SetupClient()

			cmd.CheckError(i.Run())

		},
	}

	command.PersistentFlags().IntVarP(&listenPort, "listenport", "l", 8080, "the port this integration is listening for websocket connection from the remote")

	return command
}
