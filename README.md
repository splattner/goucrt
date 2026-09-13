# A Unfolded Circle Remote Two Integration Library in Go

![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/splattner/goucrt/main.yaml)
![GitHub](https://img.shields.io/github/license/splattner/goucrt)
![GitHub (Pre-)Release Date](https://img.shields.io/github/release-date-pre/splattner/goucrt)

> **⚠️ DISCLAIMER**
> This is Work in Progress and might not yet be functional.

This [Unfolded Circle Remote Two](https://www.unfoldedcircle.com/) integration driver written in Go implements is written as generic as possible to be used with any kind of device to control.

Currently this repository implements a driver for the following devices:

* [DeCONZ](https://dresden-elektronik.github.io/deconz-rest-doc/)
* [Shelly](https://www.shelly.com/)
* [Tasmota](https://tasmota.github.io/docs/)

Intergrations implemented using this library:

* [Denon Audio/Video Reveiver](https://github.com/splattner/remotetwo-integration-denonavr)

## Device / Clients

### Generic Client

The generic client is not functional. You need to implement your own Client for the device you want to control

### Deconz

Run with `ucrt deconz`

This client currently implements [`Light` entities](https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_light.md) for discovered DeCONZ Lights and Groups and [`Sensor` entitites](https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_sensor.md) for selected DeCONZ sensors.

### Shelly

Run with `ucrt shelly`

This client currently implements [`Switch` entities](https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/switch_light.md) for discovered Shelly Devices. It uses MQTT to discover and control Shelly devices.

### Tasmota

Run with `ucrt tasmota`

This client currently implements [`Switch` entities](https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/switch_light.md) and [`Light` entities](https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_light.md) It uses MQTT to discover and control Tasmota devices.

Currently on the following Sonoff device types are supported

* `0` Sonoff Basic results in a Switch entity
* `4` RGBW results in a Light entity

*Note* The light entity does not really suport RGBW. So currently, as a somehow working workaround, when setting the brithness is set to 0, the entity changes between RGB & W Settings.

## How to use

```bash
Unfolder Circle Remote Two integration

Usage:
  ucrt-amd64 [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  deconz      Start Deconz Ingegration
  help        Help about any command
  shelly      Start Shelly Ingegration
  tasmota     Start Tasmota Ingegration

Flags:
      --debug                         Enable debug log level
      --disableMDNS                   Disable integration advertisement via mDNS
  -h, --help                          help for ucrt-amd64
  -l, --listenPort int                the port this integration is listening for websocket connection from the remote (default 8080)
      --registration                  Enable driver registration on the Remote Two instead of mDNS advertisement
      --registrationPin string        Pin of the RemoteTwo for driver registration
      --registrationUsername string   Username of the RemoteTwo for driver registration (default "web-configurator")
      --remoteTwoIP string            IP Address of your Remote Two instance (disables Remote Two discovery)
      --remoteTwoPort int             Port of your Remote Two instance (disables Remote Two discovery) (default 80)
      --ucconfighome string           Configuration directory to save the user configuration from the driver setup (default "./ucconfig/")
      --websocketPath string          path where this integration is available for websocket connections (default "/ws")

Use "ucrt-amd64 [command] --help" for more information about a command.

```

### As a Container

You can start the Integration as a container by using the released container images:

Example to start the DenonAVR Integration:

```bash
docker run ghcr.io/splattner/goucrt:v0.1.7 denonavr
```

To keep the setup data persistet mount a volume to `/app/ucconfig`:

```bash
docker run -v ./localdir:/app/ucconfig ghcr.io/splattner/goucrt:v0.1.7 denonavr
```

For the mDNS adventisement to work correctly I suggest starting the integration in the `host` network. And you can set your websocket listening port with the environment variable `UC_INTEGRATION_HTTP_PORT`:

```bash
docker run --net=host -e UC_INTEGRATION_HTTP_PORT=10000 -v ./localdir:/app/ucconfig ghcr.io/splattner/goucrt:v0.1.7 denonavr
```

### As a custom-installed driver on the Remote

Since firmware v1.9.0, the remote can install an integration driver directly as a `.tar.gz` archive - no
separate host needed. goucrt's deCONZ, Shelly and Tasmota clients each build as their own installable
archive; see the [Core-API's custom driver installation docs](https://github.com/unfoldedcircle/core-api/blob/main/doc/integration-driver/driver-installation.md)
for the full format and the remote's sandbox environment.

Build all three archives:

```bash
make custom-driver-archives
```

This produces `custom-driver-dist/{deconz,shelly,tasmota}-custom-driver.tar.gz`, each a statically linked
`linux/arm64` binary at `bin/driver` plus a generated `driver.json` and the client's icon at the archive
root. Released versions are attached to each [GitHub release](https://github.com/splattner/goucrt/releases)
alongside the container images.

Install one in the remote's web-configurator: integrations, _Add new_, _Install custom_, and upload the
archive. Or via the REST API:

```bash
curl --location 'http://$REMOTE_IP/api/intg/install' \
  --user 'web-configurator:$PIN' \
  --form 'file=@"deconz-custom-driver.tar.gz"'
```

Setup data persists across restarts in the remote-managed `$UC_CONFIG_HOME`/`$UC_DATA_HOME` directories -
nothing to mount, unlike the container image. mDNS advertisement and self-registration are disabled
automatically in this mode (see [Environment Variables](#environment-variables) below): the remote already
knows how to reach the driver from the installation itself.

### Configuration

#### Environment Variables

The following environment variables exist in addition to the configuration file:

| Variable                     | Values               |Description |
|------------------------------|----------------------|--------------------------------------------------------------------------------|
| UC_CONFIG_HOME               | _directory path_     | Configuration directory to save the user configuration from the driver setup.<br>Default: `./ucconfig/` |
| UC_DATA_HOME                  | _directory path_     | Directory for a driver's own application data, separate from user setup data.<br>Default: same as `UC_CONFIG_HOME`. goucrt itself doesn't write anything here yet; it's exposed for driver code to use. |
| UC_DISABLE_MDNS_PUBLISH      | `true` / `false`     | Disables mDNS service advertisement.<br>Default: `false` |
| UC_INTEGRATION_HTTP_PORT | `int` | The port this integration is listening for websocket connections from the remote. Set automatically by the remote when running as a [custom-installed driver](#as-a-custom-installed-driver-on-the-remote); always takes precedence over `UC_INTEGRATION_LISTEN_PORT`/`--listenPort` when set.<br> Default: `8080` |
| UC_INTEGRATION_LISTEN_PORT | `int` | Same as `UC_INTEGRATION_HTTP_PORT` above, kept for backwards compatibility. Use `UC_INTEGRATION_HTTP_PORT` if you can - it's the name the remote itself uses. |
| UC_INTEGRATION_INTERFACE | _host/IP_ | Host/IP address to listen on. Set automatically by the remote when running as a custom-installed driver.<br> Default: all interfaces (`0.0.0.0`) |
| UC_INTEGRATION_WEBSOCKET_PATH | `string` | Path where this integration is available for websocket connections.<br> Default: `/ws` |
| UC_RT_HOST | `string` | IP Address of your Remote Two instance (disables Remote Two discovery via mDNS for registration) |
| UC_RT_PORT | `int` | Port of your Remote Two instance (disables Remote Two discovery via mDNS for registration) |
| UC_ENABLE_REGISTRATION | `string` | Enable driver registration on the Remote Two instead of mDNS advertisement.<br> Default: `false` |
| UC_REGISTRATION_USERNAME | `string` | Username of the RemoteTwo for driver registration.<br> Default: `web-configurator` |
| UC_REGISTRATION_PIN | `string` | Pin of the RemoteTwo for driver registration |

When running as a [custom-installed driver](#as-a-custom-installed-driver-on-the-remote), mDNS advertisement
and self-registration are disabled automatically (there's no way for the sandboxed process to receive CLI
flags, so `UC_DISABLE_MDNS_PUBLISH`/`UC_ENABLE_REGISTRATION` can't be set either) - neither serves a purpose
there, since the remote already knows about the driver from the installation archive.

## Development

The generic client is in `pgk/client/client.go`. In order to implement your own client, you have to implement the following functions:

```go
// Client specific functions
// Initialize the client
// Here you can add entities if they are already known
initFunc func()
// Called by RemoteTwo when the integration is added and setup started
setupFunc      func()
// Handles connect/disconnect calls from RemoteTwo
clientLoopFunc func()
```

See Denon Client Example in in `pkg/client/denonavrclient.go`

## Todo's

* [x] Implement all available entities
* [x] Handle command calls
* [x] Allow for attribute changes
* [x] Allow for driver regisration
  * [ ] Make it more robust
* [ ] Allow for driver authentication with token/header
* [ ] Documentation, how to use, how to implement your own device
* [ ] probably way more

## How to Build and Run

```bash
# in cmd/ucrt
go get -u
go build .
```

### Docker

```bash
docker build -f build/Dockerfile -t  ghcr.io/splattner/goucrt:latest
```

## Verifying

### Checksum

### Checksums

```shell
wget https://github.com/splattner/goucrt/releases/download/v0.1.3/goucrt_0.1.3_checksums.txt
cosign verify-blob \
  --certificate-identity 'https://github.com/splattner/goucrt/.github/workflows/release.yaml@refs/tags/v0.1.3' \
  --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
  --cert https://github.com/splattner/goucrt/releases/download/v0.1.3/goucrt_0.1.3_checksums.txt.pem \
  --signature https://github.com/splattner/goucrt/releases/download/v0.1.3/goucrt_0.1.3_checksums.txt.sig \
  ./goucrt_0.1.3_checksums.txt
```

You can then download any file you want from the release, and verify it with, for example:

```shell
wget https://github.com/splattner/goucrt/releases/download/v0.1.3/goucrt_0.1.3_linux_amd64.tar.gz.sbom
wget https://github.com/splattner/goucrt/releases/download/v0.1.3/goucrt_0.1.3_linux_amd64.tar.gz
sha256sum --ignore-missing -c checksums.txt
```

And both should say "OK".

You can then inspect the `.sbom` file to see the entire dependency tree of the binary.

### Docker image

```shell
cosign verify ghcr.io/splattner/goucrt:v0.1.3 --certificate-identity 'https://github.com/splattner/goucrt/.github/workflows/release.yaml@refs/tags/v0.1.3' --certificate-oidc-issuer 'https://token.actions.githubusercontent.com'
```

## License

This project is licensed under the [**Mozilla Public License 2.0**](https://choosealicense.com/licenses/mpl-2.0/).
See the [LICENSE](LICENSE) file for details.
