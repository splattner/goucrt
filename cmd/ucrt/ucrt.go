package main

import (
	"os"
	"path/filepath"

	"github.com/splattner/goucrt/pkg/cmd"
	"github.com/splattner/goucrt/pkg/cmd/ucrt"
)

func main() {

	baseName := filepath.Base(os.Args[0])

	err := ucrt.NewCommand(baseName).Execute()
	cmd.CheckError(err)

}
