# goucrt

[![Build](https://img.shields.io/github/actions/workflow/status/splattner/goucrt/main.yaml)](https://github.com/splattner/goucrt/actions/workflows/main.yaml)
[![Release](https://img.shields.io/github/v/release/splattner/goucrt?include_prereleases)](https://github.com/splattner/goucrt/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/splattner/goucrt.svg)](https://pkg.go.dev/github.com/splattner/goucrt)
[![License](https://img.shields.io/github/license/splattner/goucrt)](LICENSE)

A Go library for writing integration drivers for the [Unfolded Circle](https://www.unfoldedcircle.com/) Remote Two
(and other remotes speaking the same [Core-API](https://github.com/unfoldedcircle/core-api)), plus a few ready-to-run
drivers built on top of it.

goucrt takes care of the protocol side — the WebSocket server, mDNS advertisement or driver registration, the setup
flow, persisting setup data, entity subscriptions and change events — so a driver only has to describe its entities
and talk to its devices.

> [!NOTE]
> goucrt is pre-1.0. It is used in real setups, but the Go API may still change between minor releases.

## Contents

- [Included drivers](#included-drivers)
- [Running a driver](#running-a-driver)
  - [As a container](#as-a-container)
  - [As a custom-installed driver on the remote](#as-a-custom-installed-driver-on-the-remote)
  - [Configuration](#configuration)
- [Writing your own driver](#writing-your-own-driver)
- [Development](#development)
- [Verifying releases](#verifying-releases)
- [License](#license)

## Included drivers

All three drivers ship in a single `ucrt` binary (one subcommand each), in the container image, and as separate
custom-driver archives for installing directly on the remote.

| Driver | Command | Talks to the device via | Entities |
| --- | --- | --- | --- |
| [deCONZ](https://dresden-elektronik.github.io/deconz-rest-doc/) | `ucrt deconz` | deCONZ REST + WebSocket API | [Light](https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_light.md) for lights and groups, [Sensor](https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_sensor.md) for temperature and humidity sensors |
| [Shelly](https://www.shelly.com/) | `ucrt shelly` | MQTT (discovery and control) | [Switch](https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_switch.md) |
| [Tasmota](https://tasmota.github.io/docs/) | `ucrt tasmota` | MQTT (discovery and control) | [Switch](https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_switch.md) and [Light](https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_light.md) |

**Tasmota notes:** only two Tasmota device types are supported so far:

- `0` Sonoff Basic → Switch entity
- `4` RGBW → Light entity

The remote's light entity has no real RGBW mode. As a workaround, setting the brightness to 0 switches the entity
between its RGB and white channels.

### Built with goucrt elsewhere

- [Denon AV receiver](https://github.com/splattner/remotetwo-integration-denonavr)

Built one yourself? PRs adding it here are welcome.

## Running a driver

```text
$ ucrt --help
Unfolded Circle Remote Two integration

Usage:
  ucrt [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  deconz      Start Deconz Integration
  help        Help about any command
  shelly      Start Shelly Integration
  tasmota     Start Tasmota Integration

Flags:
      --bindInterface string          Host/IP address to listen on (default: all interfaces). Overridden by UC_INTEGRATION_INTERFACE when set.
      --datahome string               Directory for a driver's own application data; defaults to ucconfighome if unset
      --debug                         Enable debug log level
      --disableMDNS                   Disable integration advertisement via mDNS
  -h, --help                          help for ucrt
  -l, --listenPort int                the port this integration is listening for websocket connection from the remote (default 8080)
      --registration                  Enable driver registration on the Remote Two instead of mDNS advertisement
      --registrationPin string        Pin of the RemoteTwo for driver registration
      --registrationUsername string   Username of the RemoteTwo for driver registration (default "web-configurator")
      --remoteTwoIP string            IP Address of your Remote Two instance (disables Remote Two discovery)
      --remoteTwoPort int             Port of your Remote Two instance (disables Remote Two discovery) (default 80)
      --ucconfighome string           Configuration directory to save the user configuration from the driver setup (default "./ucconfig/")
      --websocketPath string          path where this integration is available for websocket connections (default "/ws")
```

Prebuilt `linux/amd64` and `linux/arm64` binaries are attached to every
[GitHub release](https://github.com/splattner/goucrt/releases).

By default a driver advertises itself via mDNS, and you add it from the remote's web configurator
(_Integrations_ → _Add new_). If mDNS doesn't work on your network, use `--registration` with `--registrationPin`
(and optionally `--remoteTwoIP`) to have the driver register itself with the remote instead.

### As a container

Multi-arch images are published to `ghcr.io/splattner/goucrt` with the tags `vX.Y.Z`, `vX` and `latest`.
The image's entrypoint is `ucrt`, so pass the driver name as the command:

```bash
docker run --net=host \
  -v ./ucconfig:/app/ucconfig \
  ghcr.io/splattner/goucrt:latest deconz
```

- `--net=host` lets mDNS advertisement reach the remote.
- The volume on `/app/ucconfig` keeps the setup data across container restarts.
- To change the listen port, set `UC_INTEGRATION_LISTEN_PORT` (e.g. `-e UC_INTEGRATION_LISTEN_PORT=10000`).
  Don't use `UC_INTEGRATION_HTTP_PORT` here: it makes the driver believe it is running on the remote and turns
  mDNS off (see [Environment variables](#environment-variables)).

### As a custom-installed driver on the remote

Since firmware v1.9.0, the remote can install an integration driver directly from a `.tar.gz` archive, so you
don't need a separate host. Each included driver is released as its own installable archive
(`deconz-custom-driver.tar.gz`, `shelly-custom-driver.tar.gz`, `tasmota-custom-driver.tar.gz`) on the
[GitHub releases](https://github.com/splattner/goucrt/releases). See the Core-API's
[custom driver installation docs](https://github.com/unfoldedcircle/core-api/blob/main/doc/integration-driver/driver-installation.md)
for the archive format and the remote's sandbox environment.

Install one from the web configurator (_Integrations_ → _Add new_ → _Install custom_, then upload the archive),
or through the REST API:

```bash
curl --location "http://$REMOTE_IP/api/intg/install" \
  --user "web-configurator:$PIN" \
  --form 'file=@"deconz-custom-driver.tar.gz"'
```

In this mode the remote manages the `$UC_CONFIG_HOME`/`$UC_DATA_HOME` directories, so setup data persists
without any volumes to mount. mDNS advertisement and self-registration are switched off automatically, because
the remote already knows about the driver from the installation.

To build the archives yourself, see [Development](#development).

### Configuration

Every flag can also be set through an environment variable. The remote itself sets some of them when it runs a
custom-installed driver.

#### Environment variables

| Variable | Flag | Value | Description |
| --- | --- | --- | --- |
| `UC_CONFIG_HOME` | `--ucconfighome` | directory | Where setup data from the driver setup is saved.<br>Default: `./ucconfig/` |
| `UC_DATA_HOME` | `--datahome` | directory | Directory for a driver's own application data, separate from setup data.<br>Default: same as `UC_CONFIG_HOME`. goucrt itself doesn't write anything here yet; it's there for driver code to use. |
| `UC_INTEGRATION_LISTEN_PORT` | `--listenPort` | int | Port the driver listens on for the remote's WebSocket connection.<br>Default: `8080` |
| `UC_INTEGRATION_HTTP_PORT` | — | int | **Set by the remote** for custom-installed drivers. Takes precedence over `UC_INTEGRATION_LISTEN_PORT`/`--listenPort`, and its presence tells goucrt it is running on the remote: mDNS advertisement and self-registration are turned off. Don't set it yourself when running standalone. |
| `UC_INTEGRATION_INTERFACE` | `--bindInterface` | host/IP | Address to listen on. Set by the remote for custom-installed drivers, and then takes precedence over the flag.<br>Default: all interfaces (`0.0.0.0`) |
| `UC_INTEGRATION_WEBSOCKET_PATH` | `--websocketPath` | string | Path of the WebSocket endpoint.<br>Default: `/ws` |
| `UC_DISABLE_MDNS_PUBLISH` | `--disableMDNS` | `true`/`false` | Turns off mDNS advertisement.<br>Default: `false` |
| `UC_ENABLE_REGISTRATION` | `--registration` | `true`/`false` | Register the driver on the remote instead of relying on mDNS advertisement.<br>Default: `false` |
| `UC_REGISTRATION_USERNAME` | `--registrationUsername` | string | Username used for registration.<br>Default: `web-configurator` |
| `UC_REGISTRATION_PIN` | `--registrationPin` | string | PIN of the remote used for registration. |
| `UC_RT_HOST` | `--remoteTwoIP` | host/IP | Address of the remote for registration. Setting it skips discovering the remote via mDNS. |
| `UC_RT_PORT` | `--remoteTwoPort` | int | Port of the remote for registration. Setting it skips discovering the remote via mDNS.<br>Default: `80` |

`--debug` has no environment variable.

## Writing your own driver

A driver is a `Client` plus an `Integration`. The client declares the driver's metadata and entities and implements
a few hooks; the integration does the protocol work.

| Hook | When it's called |
| --- | --- |
| `InitFunc` | Once at start-up. Add entities you already know about here. |
| `SetupFunc` | When the user adds the driver on the remote and the setup flow starts. Setup data is already persisted when this runs. |
| `SetDriverUserDataFunc` | When the user submits a page your setup asked for (optional). |
| `AbortSetupFunc` | When the remote aborts setup, e.g. to stop a running discovery (optional). |
| `ClientLoopFunc` | When the remote connects. Runs until it receives `"disconnect"` on `c.Messages`. |

A minimal driver with one switch:

```go
package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/splattner/goucrt/pkg/cmd"
	"github.com/splattner/goucrt/pkg/entities"
	"github.com/splattner/goucrt/pkg/integration"
)

type MyClient struct {
	integration.Client
}

func NewMyClient(i *integration.Integration) *MyClient {
	c := &MyClient{Client: *integration.NewClient(i)}

	i.SetMetadata(&integration.DriverMetadata{
		DriverId:  "mydevice",
		Developer: integration.Developer{Name: "Jane Doe"},
		Name:      integration.LanguageText{En: "My Device"},
		Version:   "0.1.0",
	})

	c.InitFunc = c.init
	c.SetupFunc = c.setup
	c.ClientLoopFunc = c.loop

	return c
}

func (c *MyClient) init() {
	sw := entities.NewSwitchEntity("my_switch", entities.LanguageText{En: "My Switch"}, "")
	sw.AddFeature(entities.OnOffSwitchEntityFeatures)
	sw.MapCommand(entities.OnSwitchEntityCommand, func() error { /* talk to your device */ return nil })
	sw.MapCommand(entities.OffSwitchEntityCommand, func() error { return nil })

	_ = c.IntegrationDriver.AddEntity(sw)
}

func (c *MyClient) setup(integration.SetupData) {
	// Setup data is already persisted by the library; just finish the flow.
	c.IntegrationDriver.SetDriverSetupState(integration.StopEvent, integration.OkState, "", nil)
}

func (c *MyClient) loop() {
	c.SetDeviceState(integration.ConnectedDeviceState)
	defer c.SetDeviceState(integration.DisconnectedDeviceState)

	for msg := range c.Messages {
		if msg == "disconnect" {
			return
		}
	}
}

func main() {
	root := &cobra.Command{
		Use: "mydevice-driver",
		Run: func(*cobra.Command, []string) {
			var config integration.Config
			cmd.CheckError(viper.Unmarshal(&config))

			i, err := integration.NewIntegration(config)
			cmd.CheckError(err)

			NewMyClient(i).InitClient()
			cmd.CheckError(i.Run())
		},
	}
	cmd.CheckError(cmd.BindStandardFlags(root)) // --listenPort, UC_CONFIG_HOME, ... like ucrt itself
	cmd.CheckError(root.Execute())
}
```

When the device reports a change, update the entity with `SetAttributes` and goucrt forwards it to the remote if
the entity is subscribed.

Supported entity types, all in [`pkg/entities`](pkg/entities): button, climate, cover, IR emitter, light,
media player (including browsing and search), remote, select, sensor, switch and voice assistant.

Useful references in this repo:

- [`pkg/clients/shelly`](pkg/clients/shelly/shellyclient.go) is a small real driver: setup settings, MQTT
  discovery and entities added at runtime.
- [`pkg/clients/deconz`](pkg/clients/deconz/deconzclient.go) shows a multi-step setup flow with user input.
- [`cmd/shelly-driver`](cmd/shelly-driver/main.go) and [`tools/gendriverjson`](tools/gendriverjson/main.go) show
  how to build a single-driver binary and the `driver.json` needed for a custom-install archive.
- `entities.CommandParams` gives typed access to command parameters. Use it instead of asserting on the raw map:
  JSON numbers always arrive as `float64`.

## Development

Requires Go (see [`go.mod`](go.mod) for the version) and `make`.

```bash
make help                    # list all targets
go build -o ucrt ./cmd/ucrt  # build the multi-driver binary
make test                    # go test -race ./...
make lint                    # go fmt, go vet and golangci-lint
make custom-driver-archives  # build the custom-install archives into custom-driver-dist/
make docker-build-amd64      # build the binary and a linux/amd64 image
```

`make custom-driver-archives` produces `custom-driver-dist/{deconz,shelly,tasmota}-custom-driver.tar.gz`. Each
contains a statically linked `linux/arm64` binary at `bin/driver`, a generated `driver.json` and the driver's
icon.

Releases are built by [GoReleaser](.goreleaser.yaml) when a `v*` tag is pushed.

### Core-API spec

goucrt follows the [Unfolded Circle Core-API](https://github.com/unfoldedcircle/core-api). A copy of the spec is
vendored in [`internal/spec/core-api`](internal/spec/core-api), and the tests check goucrt's entity features,
commands and attributes against it. `make check-spec-freshness` tells you whether the copy is behind upstream;
[`internal/spec/core-api/README.md`](internal/spec/core-api/README.md) explains how to update it.

### Roadmap

- [x] All entity types from the Core-API
- [x] Entity commands and attribute changes
- [x] Driver registration on the remote
- [x] Custom-installed driver archives
- [ ] More robust driver registration
- [ ] Driver authentication with token/header
- [ ] More documentation on writing your own driver

## Verifying releases

Release artifacts and container images are signed with [cosign](https://github.com/sigstore/cosign) keyless
signing from the GitHub Actions release workflow. Set the release you want to check:

```bash
VERSION=0.6.1
```

### Binaries

Verify the signature on the checksum file:

```bash
BASE=https://github.com/splattner/goucrt/releases/download/v${VERSION}
wget "$BASE/goucrt_${VERSION}_checksums.txt"
cosign verify-blob \
  --certificate-identity "https://github.com/splattner/goucrt/.github/workflows/release.yaml@refs/tags/v${VERSION}" \
  --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
  --cert "$BASE/goucrt_${VERSION}_checksums.txt.pem" \
  --signature "$BASE/goucrt_${VERSION}_checksums.txt.sig" \
  "./goucrt_${VERSION}_checksums.txt"
```

Then download any file from the release and check it against the verified checksums:

```bash
wget "$BASE/goucrt_${VERSION}_linux_amd64.tar.gz" "$BASE/goucrt_${VERSION}_linux_amd64.tar.gz.sbom.json"
sha256sum --ignore-missing -c "goucrt_${VERSION}_checksums.txt"
```

Both should print `OK`. The SBOM lists the binary's full dependency tree.

### Container image

```bash
cosign verify "ghcr.io/splattner/goucrt:v${VERSION}" \
  --certificate-identity "https://github.com/splattner/goucrt/.github/workflows/release.yaml@refs/tags/v${VERSION}" \
  --certificate-oidc-issuer 'https://token.actions.githubusercontent.com'
```

## License

goucrt is licensed under the [Mozilla Public License 2.0](LICENSE).
