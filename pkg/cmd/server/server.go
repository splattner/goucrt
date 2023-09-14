package server

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/splattner/goucrt/pkg/cmd"
	"github.com/splattner/goucrt/pkg/integration"
)

func NewCommand() *cobra.Command {

	var command = &cobra.Command{
		Use:   "server",
		Short: "Run a Unfolded Circle Remote Two integratin",
		Long:  "Run a Unfolded Circle Remote Two integratin",
		Run: func(c *cobra.Command, args []string) {
			log.SetOutput(os.Stdout)
			log.Println("Integration run")

			i, err := integration.NewIntegration()
			cmd.CheckError(err)

			cmd.CheckError(i.Run())

		},
	}

	return command
}
