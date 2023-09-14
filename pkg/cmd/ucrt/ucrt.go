package ucrt

import (
	"github.com/spf13/cobra"
	"github.com/splattner/goucrt/pkg/cmd/server"
)

func NewCommand(name string) *cobra.Command {

	c := &cobra.Command{
		Use:   name,
		Short: "Unfolder Circle Remote Two integration",
		Long:  `Unfolder Circle Remote Two integration`,
	}

	c.AddCommand(
		server.NewCommand(),
	)

	return c
}
