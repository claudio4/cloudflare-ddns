package cloudflare

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

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
	ip6 := ip.To16()
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
	records, err := api.api.DNSRecords(context.Background(), zoneID, rr)
	if err != nil {
		return fmt.Errorf("error retriving zone records: %w", err)
	}
	if len(records) == 0 {
		rr.Content = content
		resp, err := api.api.CreateDNSRecord(context.Background(), zoneID, rr)
		if err != nil {
			return fmt.Errorf("error creating the record: %w", err)
		}
		if len(resp.Errors) > 0 {
			var b strings.Builder
			b.WriteString("error creating the record: ")
			lastELement := len(resp.Errors) - 1
			for i, err := range resp.Errors {
				if i != 0 && i < lastELement {
					b.WriteString("; ")
				}
				b.WriteString("(code ")
				b.WriteString(strconv.Itoa(err.Code))
				b.WriteRune(')')
				b.WriteString(err.Message)
			}
			return errors.New(b.String())
		}
		return nil
	}

	rr = records[0]
	if rr.Content == content {
		return nil
	}
	rr.Content = content
	err = api.api.UpdateDNSRecord(context.Background(), zoneID, rr.ID, rr)
	if err != nil {
		return fmt.Errorf("error updating the record: %w", err)
	}

	return nil
}
