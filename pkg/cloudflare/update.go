package cloudflare

import (
	"errors"
	"fmt"
	"net"

	cf "github.com/cloudflare/cloudflare-go"
)

var (
	ErrInvalidIPType = errors.New("provided IP type doesn't match record type")
)

func (api *API) UpdateARecord(name string, ip net.IP) error {
	ip4 := ip.To4()
	if ip4 == nil {
		return ErrInvalidIPType
	}
	return api.updateRecordContent(name, ip4.String(), "A")
}

func (api *API) UpdateAAAARecord(name string, ip net.IP) error {
	ip6 := ip.To4()
	if ip6 == nil {
		return ErrInvalidIPType
	}
	return api.updateRecordContent(name, ip6.String(), "AAAA")
}

func (api *API) updateRecordContent(name, content string, recordType string) error {
	zoneID, err := api.getRecordZoneID(name)
	if err != nil {
		return fmt.Errorf("error guesssing the zone ID: %w", err)
	}
	rr := cf.DNSRecord{
		Name: name,
		Type: recordType,
	}
	records, err := api.api.DNSRecords(zoneID, rr)
	if err != nil {
		return fmt.Errorf("error retriving zone records: %w", err)
	}

	rr = records[0]
	rr.Content = content
	err = api.api.UpdateDNSRecord(zoneID, rr.ID, rr)
	if err != nil {
		return fmt.Errorf("error updating the record: %w", err)
	}

	return nil
}
