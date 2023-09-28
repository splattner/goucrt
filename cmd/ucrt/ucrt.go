package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/splattner/goucrt/pkg/cmd"
	"github.com/splattner/goucrt/pkg/cmd/ucrt"

	log "github.com/sirupsen/logrus"
)

const RELEASEDATE string = "21.09.2023"

func main() {

	baseName := filepath.Base(os.Args[0])

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("UC_")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Info("Unable to read config")
	}

	err := ucrt.NewCommand(baseName).Execute()
	cmd.CheckError(err)

}
