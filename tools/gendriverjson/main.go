// Command gendriverjson writes the driver.json metadata file a custom-installed driver archive
// needs at its root (see doc/integration-driver/driver-installation.md), generated from the exact
// same DriverMetadata the corresponding client sets on its Integration at runtime - so the two
// can never drift apart. It works by constructing a throwaway Integration and calling the client's
// own NewXClient constructor, the same call every driver binary makes at startup, then reading
// back the DriverMetadata that call set and marshalling it directly: DriverMetadata's fields and
// JSON tags already match driver.json's schema (the connection-only fields it doesn't need -
// driver_url, auth_method - are never set by any client here, so they're simply absent via
// omitempty).
//
// Run from the repository root:
//
//	go run ./tools/gendriverjson -client deconz -out driver.json
package main

import (
	"encoding/json"
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

	i, err := integration.NewIntegration(integration.Config{})
	if err != nil {
		fatalf("NewIntegration: %v", err)
	}

	// Each NewXClient constructor only builds its DriverMetadata and registers function pointers -
	// it doesn't start any network activity, so it's safe to call here without InitClient/Run.
	switch *client {
	case "deconz":
		deconzclient.NewDeconzClient(i)
	case "shelly":
		shellyclient.NewShellyClient(i)
	case "tasmota":
		tasmotaclient.NewTasmotaClient(i)
	default:
		fatalf("unknown -client %q: must be deconz, shelly, or tasmota", *client)
	}

	if i.Metadata == nil {
		fatalf("%s's client constructor didn't set driver metadata", *client)
	}

	data, err := json.MarshalIndent(i.Metadata, "", "  ")
	if err != nil {
		fatalf("marshal driver metadata: %v", err)
	}
	data = append(data, '\n')

	if err := os.WriteFile(*out, data, 0o644); err != nil {
		fatalf("write %s: %v", *out, err)
	}

	fmt.Printf("wrote %s for %s\n", *out, *client)
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
