package cloudflare

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const baseuri = "https://api.cloudflare.com/client/v4/"

type basicCloudflareResponse struct {
	Success  bool     `json:"success"`
	Errors   []string `json:"errors"`
	Messages []string `json:"messages"`
}

type listRecordSResponse struct {
	basicCloudflareResponse
	Result []struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Content string `json:"content"`
		TTL     int    `json:"ttl"`
	} `json:"result"`
}

//GetRecordDetails Gets record's id and ttl
func GetRecordDetails(email, apikey, zonedID, name string) (string, string, int, error) {
	res, err := doCloudFlareRequest("GET", fmt.Sprintf("zones/%s/dns_records", zonedID), email, apikey, nil)
	if err != nil {
		return "", "", 0, errors.Wrap(err, "Error doing request to Cloudflare")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", "", 0, errors.Wrap(err, "Error reading response body")
	}

	var r listRecordSResponse
	json.Unmarshal(body, &r)
	if !r.Success {
		return "", "", 0, fmt.Errorf("The request failed with the following errors: %+v", r.Errors)
	}
	var id, content string
	var ttl int
	for _, record := range r.Result {
		if record.Name == name {
			id = record.ID
			content = record.Content
			ttl = record.TTL
			break
		}
	}

	return id, content, ttl, nil
}

//SetRecord sets record
func SetRecord(email, apikey, zonedID, recordID, recordType, name, content string, ttl int) error {
	res, err := doCloudFlareRequest(
		"PUT",
		fmt.Sprintf("zones/%s/dns_records/%s", zonedID, recordID),
		email,
		apikey,
		strings.NewReader(fmt.Sprintf(`{"type":"%s","name":"%s","content":"%s","ttl":%d}`, recordType, name, content, ttl)),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "Error reading response body")
	}

	var resp basicCloudflareResponse
	json.Unmarshal(body, &resp)
	if !resp.Success {
		return fmt.Errorf("The request failed with the following errors: %+v", resp.Errors)
	}
	return nil
}

func doCloudFlareRequest(method, uri, email, apikey string, requestBody io.Reader) (*http.Response, error) {
	client := &http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequest(method, baseuri+uri, requestBody)
	if err != nil {
		return nil, errors.Wrapf(err, "Error creating %v request", method)
	}

	req.Header.Set("X-Auth-Email", email)
	req.Header.Set("X-Auth-Key", apikey)
	if method != "GET" {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error doing request to Cloudflare")
	}

	return res, nil
}
