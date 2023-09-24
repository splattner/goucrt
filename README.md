# A Unfolded Circle Remote Two Integration in Go

![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/splattner/goucrt/docker-publish.yml)


> **⚠️ DISCLAIMER**
> This is Work in Progress and might not yet be functional.

This [Unfolded Circle Remote Two](https://www.unfoldedcircle.com/) integration driver written in Go implements is written as generic as possible to be used with any kind of device to control.

Currently this repository implements a driver for the following devices:

* Denon Audio/Video Reveiver
* [DeCONZ](https://dresden-elektronik.github.io/deconz-rest-doc/)

## Device / Clients

### Generic Client

The generic client is not functional. You need to implement your own Client for the device you want to control

### Denon Audio/Video Reveiver

Run with `ucrt denon`

This client currently implements a [`MediaPlayer` entity](https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_media_player.md) and some [`Button` entities](https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_button.md)

### Deconz

Run with `ucrt deconz`

This client currently implements [`Light` entities](https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_light.md) for discovered DeCONZ Lights and [`Sensor` entitites](https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_sensor.md) for selected DeCONZ sensors.

## How to use

```bash
Unfolder Circle Remote Two integration

Usage:
  ucrt [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  deconz      Deconz
  denonavr    Denon AVR
  help        Help about any command

Flags:
      --debug                         Enable debug log level
  -h, --help                          help for ucrt
  -l, --listenPort int                the port this integration is listening for websocket connection from the remote (default 8080)
      --mdns                          Enable integration advertisement via mDNS (default true)
      --registration                  Enable driver registration on the Remote Two instead of mDNS advertisement
      --registrationPin string        Pin of the RemoteTwo for driver registration
      --registrationUsername string   Username of the RemoteTwo for driver registration (default "web-configurator")
      --remoteTwoIP string            IP Address of your Remote Two instance (disables Remote Two discovery)
      --remoteTwoPort int             Port of your Remote Two instance (disables Remote Two discovery) (default 80)
      --websocketPath string          path where this integration is available for websocket connections (default "/ws")

Use "ucrt [command] --help" for more information about a command.

```

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
docker build -f build/Dockerfile -t splattner/goucrt
```

## License

This project is licensed under the [**Mozilla Public License 2.0**](https://choosealicense.com/licenses/mpl-2.0/).
See the [LICENSE](LICENSE) file for details.