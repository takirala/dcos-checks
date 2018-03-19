package common

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dcos/dcos-checks/client"
	"github.com/dcos/dcos-checks/constants"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// URLFields is used to construct the url
type URLFields struct {
	Host string
	Port int
	Path string
}

// HTTPRequest verifies the results of the request
func HTTPRequest(cfg *CLIConfigFlags, urlOptions URLFields) (int, []byte, error) {
	httpClient, err := client.NewClient(cfg.IAMConfig, cfg.CACert)
	if err != nil {
		return 0, nil, errors.Wrap(err, "unable to create HTTP client")
	}

	url, err := GetURL(httpClient, cfg, urlOptions)
	if err != nil {
		return 0, nil, err
	}

	logrus.Debugf("GET %s", url)
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return 0, nil, errors.Wrap(err, "unable to create a new HTTP request")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, nil, errors.Wrapf(err, "unable to execute GET %s", url)
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, errors.Wrapf(err, "unable to read response body")
	}
	defer resp.Body.Close()

	return resp.StatusCode, responseData, nil
}

// GetURL returns a URL appropriate for the supplied config flags and url fields
func GetURL(httpClient *http.Client, cfg *CLIConfigFlags, urlOptions URLFields) (*url.URL, error) {
	scheme := constants.HTTPScheme
	if cfg.ForceTLS {
		scheme = constants.HTTPSScheme
	}
	host := urlOptions.Host
	if host == "" {
		ip, err := cfg.IP(httpClient)
		if err != nil {
			return nil, err
		}
		host = ip.String()
	}
	if urlOptions.Port == 0 {
		return &url.URL{
			Scheme: scheme,
			Host:   host,
			Path:   urlOptions.Path,
		}, nil
	}
	return &url.URL{
		Scheme: scheme,
		Host:   net.JoinHostPort(host, strconv.Itoa(urlOptions.Port)),
		Path:   urlOptions.Path,
	}, nil
}
