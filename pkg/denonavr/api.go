package denonavr

import (
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

type DenonAVRApi struct {
	Host    string
	Port    int
	timeout int
}

func (api DenonAVRApi) get(request_url string, port int) (*http.Response, error) {

	resp, err := http.Get("http://" + api.Host + url.QueryEscape(request_url))
	if err != nil {
		log.WithError(err).Error("Cannot make Get call")
		return nil, err
	}

	return resp, nil

}
