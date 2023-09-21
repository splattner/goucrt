package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	b64 "encoding/base64"

	log "github.com/sirupsen/logrus"

	"github.com/grandcat/zeroconf"
)

type DriverRegistration struct {
	DriverId        string          `json:"driver_id,omitempty"`
	Name            LanguageText    `json:"name"`
	DriverURL       string          `json:"driver_url"`
	Version         string          `json:"version"`
	Icon            string          `json:"icon"`
	Enabled         bool            `json:"enabled"`
	Description     LanguageText    `json:"description"`
	DeviceDiscovery bool            `json:"device_discovery"`
	SetupDataSchema SetupDataSchema `json:"setup_data_schema"`
	ReleaseDate     string          `json:"release_date,omitempty"`
}

// Register the integration with Remote Two
// TODO: make this more robust and nicer
func (i *Integration) registerIntegration() {

	myip := GetLocalIP()
	log.WithField("MyIP", myip).Info("Register driver at availabe Remote Two instances")

	driverRegistration := DriverRegistration{
		DriverId:        i.SetupData["driver_id"],
		Name:            i.Metadata.Name,
		DriverURL:       "ws://" + myip + i.listenAddress + i.config["websocketPath"].(string),
		Version:         i.Metadata.Version,
		Icon:            i.Metadata.Icon,
		Enabled:         true,
		Description:     i.Metadata.Description,
		DeviceDiscovery: false,
		SetupDataSchema: i.Metadata.SetupDataSchema,
	}

	entries := make(chan *zeroconf.ServiceEntry)

	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {

			remoteTwoURL := "http://" + entry.AddrIPv4[0].String() + ":" + fmt.Sprint(entry.Port)

			log.WithFields(log.Fields{
				"Remote Two": remoteTwoURL,
				"IP":         entry.AddrIPv4[0].String()}).Info("Register Integration with Remote Two")

			data, err := json.Marshal(driverRegistration)
			req, err := http.NewRequest("POST", remoteTwoURL+"/api/intg/drivers", bytes.NewReader(data))
			if err != nil {
				log.WithError(err).Fatal("impossible to build request")
			}
			req.Header.Set("Content-Type", "application/json")

			credentials := b64.StdEncoding.EncodeToString([]byte("web-configurator:" + i.config["registrationPin"].(string)))

			req.Header.Set("Authorization", "Basic "+credentials)
			client := http.Client{Timeout: 10 * time.Second}

			// send the request
			res, err := client.Do(req)
			if err != nil {
				log.WithError(err).Fatal("impossible to send request")
			}

			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				log.Fatalf("impossible to read all body of response: %s", err)
			}
			json.Unmarshal(resBody, &driverRegistration)

			log.WithField("Response", string(resBody)).Debug("Driver Registration")

			i.SetupData["driver_id"] = driverRegistration.DriverId
			i.persistSetupData()
		}
		log.Info("No more entries.")
	}(entries)

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)

	defer cancel()

	err = resolver.Browse(ctx, "_uc-remote._tcp", "local.", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	<-ctx.Done()
}

// GetLocalIP returns the non loopback local IP of the host
// TODO: make this more robust, what if more ifaces are available
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
