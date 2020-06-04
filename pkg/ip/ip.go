package ip

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

type version int

const (
	v4 version = iota
	v6
)

var (
	ErrInvalidIPReceived = errors.New("the received IP is invalid")
)

const apiEndpoint = "https://api6.ipify.org/"

func getIP(v version) (net.IP, error) {
	req, err := http.NewRequest("GET", "https://api6.ipify.org/", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating the request %w", err)
	}

	res, err := httpDoIPVersionRestricted(req, v)
	if err != nil {
		return nil, fmt.Errorf("error fetching the IP %w", err)
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
