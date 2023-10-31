package denonavr

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ziutek/telnet"
)

type TelnetEvent struct {
	RawData string
	Command string
	Payload string
}

func (d *DenonAVR) handleTelnetEvents() {

	for {
		select {
		case event := <-d.telnetEvents:
			log.WithFields(log.Fields{
				"cmd":     event.Command,
				"payload": event.Payload,
			}).Debug("received telnet event")

			parsedCommand := strings.Split(event.Command, "")
			command := parsedCommand[0] + parsedCommand[1]
			param := strings.Join(parsedCommand[2:], "")

			log.WithFields(log.Fields{
				"command": DenonCommand(command),
				"param":   param,
			}).Debug("parsed telnet event")

			switch DenonCommand(command) {
			case DenonCommandMainZoneVolume:
				if param != "MAX" {
					log.Debug("Main Zone Volume from telnet")
					volume, err := strconv.ParseFloat(param, 32)
					log.WithField("volume", volume).Debug("Got volume")
					if err != nil {
						log.WithError(err).Error("failed to parse volume")
					}

					if len(param) == 3 {
						volume = volume / 10
						log.WithField("volume", volume).Debug("Got volume after conversion")
					}

					log.WithField("volume", fmt.Sprintf("%0.1f", volume-80)).Debug("Got volume")

					d.SetAttribute("MainZoneVolume", fmt.Sprintf("%0.1f", volume-80))
				}

			case DenonCommandMainZoneMute:
				log.Debug("Main Zone Mute from telnet")
				d.SetAttribute("MainZoneMute", strings.ToLower(param))

			}

		}
	}
}

func (d *DenonAVR) listenTelnet() {

	go d.handleTelnetEvents()

	var err error

	for {
		d.telnet, err = telnet.DialTimeout("tcp", d.Host+":23", 5*time.Second)
		if err != nil {
			log.WithError(err).Info("failed to connect to telnet")
			continue
		}

		if err = d.telnet.Conn.(*net.TCPConn).SetKeepAlive(true); err != nil {
			log.WithError(err).Error("failed to enable tcp keep alive")
		}

		if err = d.telnet.Conn.(*net.TCPConn).SetKeepAlivePeriod(5 * time.Second); err != nil {
			log.WithError(err).Error("failed to set tcp keep alive period")
		}

		log.Debug("telnet connected")

		for {
			data, err := d.telnet.ReadString('\r')
			if err != nil {
				log.WithError(err).Errorf("failed to read form telnet")
				break
			}
			data = strings.Trim(data, " \n\r")

			parsedData := strings.Split(data, " ")
			event := TelnetEvent{}
			event.RawData = data
			event.Command = parsedData[0]
			if len(parsedData) > 1 {
				event.Payload = parsedData[1]
			}

			d.telnetEvents <- &event
		}
	}
}

// func (d *DenonAVR) sendTelnetCommand(cmd DenonCommand, payload string) {

// 	d.telnetMutex.Lock()

// 	defer d.telnetMutex.Unlock()

// 	log.WithFields(log.Fields{
// 		"cmd":     string(cmd),
// 		"payload": payload,
// 	}).Debug("send telnet command")

// 	d.telnet.Write([]byte(string(cmd) + payload + "\r"))
// }
