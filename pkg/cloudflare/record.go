package cloudflare

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const baseuri = "https://api.cloudflare.com/client/v4/"

type listRecordSResponse struct {
	Success  bool     `json:"success"`
	Errors   []string `json:"errors"`
	Messages []string `json:"messages"`
	Result   []struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Content string `json:"content"`
		TTL     int    `json:"ttl"`
	} `json:"result"`
}

//GetRecordDetails Gets record's id and ttl
func GetRecordDetails(email, apikey, zonedID, name string) (string, string, int, error) {
	client := &http.Client{Timeout: time.Second * 2}
	req, err := http.NewRequest("GET", fmt.Sprintf("%szones/%s/dns_records", baseuri, zonedID), nil)
	if err != nil {
		return "", "", 0, errors.Wrap(err, "Error creating GET request")
	}
	req.Header.Set("X-Auth-Email", email)
	req.Header.Set("X-Auth-Key", apikey)

	res, err := client.Do(req)
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
		return "", "", 0, fmt.Errorf("Cloudflare Errors: %+v - Cloudflare Messages: %+v", r.Errors, r.Messages)
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
	client := &http.Client{Timeout: time.Second * 2}

	data := fmt.Sprintf(`{"type":"%s","name":"%s","content":"%s","ttl":%d}`, recordType, name, content, ttl)
	req, err := http.NewRequest("PUT", fmt.Sprintf("%szones/%s/dns_records/%s", baseuri, zonedID, recordID), strings.NewReader(data))
	if err != nil {
		return errors.Wrap(err, "Error creating PUT request")
	}

	req.ContentLength = int64(len(data))
	req.Header.Set("X-Auth-Email", email)
	req.Header.Set("X-Auth-Key", apikey)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "Error doing request to Cloudflare")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "Error reading response body")
    }

    var resMap map[string]interface{}
	json.Unmarshal(body, &resMap)
	if !resMap["success"].(bool) {
		return fmt.Errorf("Cloudflare Errors: %+v - Cloudflare Messages: %+v", resMap["errors"], resMap["messages"])
	}
	return nil
}
