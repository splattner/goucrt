// Command deconz-driver is the deCONZ client built as its own single-purpose binary, for
// packaging as a custom-installed driver on the remote (see
// doc/integration-driver/driver-installation.md). The remote's sandbox runs the archive's
// ./bin/driver binary directly with no arguments, so unlike the main ucrt binary this doesn't
// dispatch on a subcommand - it just runs "deconz" unconditionally, reusing the exact same
// pkg/cmd/ucrt command tree (and therefore the exact same flag/env var handling) as
// `ucrt deconz` would.
package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/splattner/goucrt/pkg/cmd"
	"github.com/splattner/goucrt/pkg/cmd/ucrt"

	log "github.com/sirupsen/logrus"
)

func main() {

	baseName := filepath.Base(os.Args[0])

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("UC_")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Info("Unable to read config")
	}

	root := ucrt.NewCommand(baseName)
	root.SetArgs([]string{"deconz"})

	cmd.CheckError(root.Execute())

}
