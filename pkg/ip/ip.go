package ip

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

//Get Returns your ip
func Get() (string, error) {
	res, err := http.Get("https://checkip.amazonaws.com/")
	if err != nil {
		return "", errors.Wrap(err, "Error requesting the ip")
	}
	defer res.Body.Close()

	ip, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "Error reading the request")
	}

	return strings.TrimSpace(string(ip)), nil
}
