package ucrt

import (
	"github.com/spf13/cobra"
	"github.com/splattner/goucrt/pkg/cmd"
	"github.com/splattner/goucrt/pkg/cmd/deconz"
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

	if err := cmd.BindStandardFlags(rootCmd); err != nil {
		log.WithError(err).Error("Cannot bind standard flags")
	}

	rootCmd.AddCommand(
		deconz.NewCommand(rootCmd),
		shelly.NewCommand(rootCmd),
		tasmota.NewCommand(rootCmd),
	)

	return rootCmd
}
