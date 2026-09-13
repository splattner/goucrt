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
func (i *Integration) registerIntegration() error {

	// Use configured IP for registration instead of Remote Two discovery
	if i.Config.RegistrationPin != "" && i.Config.RemoteTwoPort > 0 {
		return i.registerWithRemoteTwo(i.Config.RemoteTwoHost, i.Config.RemoteTwoPort)
	}

	entries := make(chan *zeroconf.ServiceEntry)

	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {

			log.WithField("MDNS Record", entry).Debug("Found Remote Two instance")

			if len(entry.AddrIPv4) == 0 {
				// TODO: IPv6?
				log.Debug("No IPv4 address available. Not using this record")
				continue
			}

			// Each discovered remote registers independently; one failing shouldn't stop the
			// driver from registering with any others found during this discovery window.
			if err := i.registerWithRemoteTwo(entry.AddrIPv4[0].String(), entry.Port); err != nil {
				log.WithError(err).Error("Cannot register with discovered Remote Two instance")
			}

		}
	}(entries)

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return fmt.Errorf("failed to initialize mDNS resolver: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)

	defer cancel()

	if err := resolver.Browse(ctx, "_uc-remote._tcp", "local.", entries); err != nil {
		return fmt.Errorf("failed to browse for Remote Two instances: %w", err)
	}

	<-ctx.Done()

	return nil
}

func (i *Integration) registerWithRemoteTwo(remoteTwoIP string, remoteTwoPort int) error {

	myip := GetLocalIP()
	driverURL := "ws://" + net.JoinHostPort(myip, fmt.Sprint(i.Config.ListenPort)) + i.Config.WebsocketPath
	remoteTwoURL := "http://" + remoteTwoIP + ":" + fmt.Sprint(remoteTwoPort)

	driverRegistration := DriverRegistration{
		DriverId:        i.SetupData["driver_id"],
		Name:            i.Metadata.Name,
		DriverURL:       driverURL,
		Version:         i.Metadata.Version,
		Icon:            i.Metadata.Icon,
		Enabled:         true,
		Description:     i.Metadata.Description,
		DeviceDiscovery: false,
		SetupDataSchema: i.Metadata.SetupDataSchema,
	}

	log.WithFields(log.Fields{
		"My IP":      myip,
		"Remote Two": remoteTwoURL,
		"IP":         remoteTwoIP,
		"DriverURL":  driverURL}).Info("Register Integration with Remote Two")

	data, err := json.Marshal(driverRegistration)
	if err != nil {
		return fmt.Errorf("cannot marshal driverRegistration: %w", err)
	}
	req, err := http.NewRequest("POST", remoteTwoURL+"/api/intg/drivers", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("cannot build registration request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Authentication wit the Remote Two
	credentials := b64.StdEncoding.EncodeToString([]byte(i.Config.RegistrationUsername + ":" + i.Config.RegistrationPin))
	req.Header.Set("Authorization", "Basic "+credentials)

	// send the request
	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send registration request: %w", err)
	}

	defer res.Body.Close()

	statusCode := res.StatusCode
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("cannot read registration response body: %w", err)
	}

	log.WithFields(log.Fields{
		"Status Code": statusCode,
		"Response":    string(resBody)}).Debug("Driver Registration")

	switch statusCode {
	//case http.StatusUnprocessableEntity:

	case http.StatusCreated:
		if err := json.Unmarshal(resBody, &driverRegistration); err != nil {
			return fmt.Errorf("cannot unmarshal driverRegistration response: %w", err)
		}

		i.SetupData["driver_id"] = driverRegistration.DriverId
		if err := i.PersistSetupData(); err != nil {
			return fmt.Errorf("cannot persist setup data: %w", err)
		}
	}

	return nil
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
