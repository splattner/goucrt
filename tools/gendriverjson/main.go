// Command gendriverjson writes the driver.json metadata file a custom-installed driver archive
// needs at its root, using integration.GenerateDriverJSON - see that function's doc comment for
// how and why.
//
// Run from the repository root:
//
//	go run ./tools/gendriverjson -client deconz -out driver.json
package main

import (
	"flag"
	"fmt"
	"os"

	deconzclient "github.com/splattner/goucrt/pkg/clients/deconz"
	shellyclient "github.com/splattner/goucrt/pkg/clients/shelly"
	tasmotaclient "github.com/splattner/goucrt/pkg/clients/tasmota"
	"github.com/splattner/goucrt/pkg/integration"
)

func main() {
	client := flag.String("client", "", "driver client to generate driver.json for: deconz, shelly, or tasmota")
	out := flag.String("out", "driver.json", "output file path")
	flag.Parse()

	var newClient func(*integration.Integration)
	switch *client {
	case "deconz":
		newClient = func(i *integration.Integration) { deconzclient.NewDeconzClient(i) }
	case "shelly":
		newClient = func(i *integration.Integration) { shellyclient.NewShellyClient(i) }
	case "tasmota":
		newClient = func(i *integration.Integration) { tasmotaclient.NewTasmotaClient(i) }
	default:
		fatalf("unknown -client %q: must be deconz, shelly, or tasmota", *client)
	}

	data, err := integration.GenerateDriverJSON(newClient)
	if err != nil {
		fatalf("%s: %v", *client, err)
	}

	if err := os.WriteFile(*out, data, 0o644); err != nil {
		fatalf("write %s: %v", *out, err)
	}

	fmt.Printf("wrote %s for %s\n", *out, *client)
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
