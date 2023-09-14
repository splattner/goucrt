package main

import (
	"os"
	"path/filepath"

	"github.com/splattner/goucrt/pkg/cmd"
	"github.com/splattner/goucrt/pkg/cmd/ucrt"
	"k8s.io/klog/v2"
)

func main() {
	defer klog.Flush()

	baseName := filepath.Base(os.Args[0])

	err := ucrt.NewCommand(baseName).Execute()
	cmd.CheckError(err)

}
