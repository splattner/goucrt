// Command tasmota-driver is the Tasmota client built as its own single-purpose binary, for
// packaging as a custom-installed driver on the remote. See cmd/deconz-driver/main.go's doc
// comment - this is the same shape, just for the "tasmota" subcommand.
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
	root.SetArgs([]string{"tasmota"})

	cmd.CheckError(root.Execute())

}
