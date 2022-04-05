package ip

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type version int

const (
	apiIPv4Endpoint = "https://ipv4.icanhazip.com/"
	apiIPv6Endpoint = "https://ipv6.icanhazip.com/"
)

const (
	v4 version = iota
	v6
)

var (
	ErrInvalidIPReceived = errors.New("the received IP is invalid")
)

// HTTPClient is the client used to perform the http requests against the API
var HTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

func getIP(v version) (net.IP, error) {
	var endpoint string
	if v == v4 {
		endpoint = apiIPv4Endpoint
	} else {
		endpoint = apiIPv6Endpoint
	}
	res, err := HTTPClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error fetching the IP: %w", err)
	}

	ipStr, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response %w", err)
	}

	ip := net.ParseIP(string(bytes.TrimSpace(ipStr)))
	if ip == nil {
		return nil, ErrInvalidIPReceived
	}

	return ip, nil
}

func GetV4() (net.IP, error) {
	return getIP(v4)
}

func GetV6() (net.IP, error) {
	return getIP(v6)
}
