package deconz

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type DeconzAPIKeyReqResponse struct {
	Success DeconzAPIKeyReqResponseData `json:"success"`
}

type DeconzAPIKeyReqResponseData struct {
	Username string `json:"username"`
}

// Get a new API key from DeconZ
func (d *Deconz) GetNewAPIKey(devicetype string) (string, error) {
	log.WithFields(log.Fields{
		"host": d.host,
		"port": d.port,
	}).Info("Get a new API Key from DeCONZ")

	url := "http://" + d.host + ":" + fmt.Sprint(d.port) + "/api"

	jsonBody := []byte(`{"devicetype": "` + devicetype + `"}`)
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		log.WithError(err).Fatal("impossible to build request")
		return "", err
	}

	// send the request
	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.WithError(err).Fatal("Failed to send the request")
		return "", err
	}

	defer res.Body.Close()

	statusCode := res.StatusCode
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("impossible to read all body of response: %s", err)
		return "", err
	}

	log.WithFields(log.Fields{
		"Status Code": statusCode,
		"Response":    string(resBody)}).Debug("Get DeconZ API Key")

	switch statusCode {
	case http.StatusForbidden:
		return "", fmt.Errorf("Make sure your Gateway is unlocked by pressing the link button")

	case http.StatusOK:
		var response []DeconzAPIKeyReqResponse
		if err := json.Unmarshal(resBody, &response); err != nil {
			return "", err
		}

		for _, item := range response {
			d.apikey = item.Success.Username
		}

	}

	return d.apikey, nil
}
