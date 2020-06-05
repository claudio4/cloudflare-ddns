package cloudflare

import (
	"errors"
	"fmt"
	"time"

	"github.com/miekg/dns"
)

var (
	ErrNoDNSResults           = errors.New("the dns return no results for the given record")
	ErrUnexpectedResponseType = errors.New("the dns return an unexpected response type")
)

type zoneCacheEntry struct {
	zoneID string
	// expires sets the expire date, is formatted as an Unix timestamp
	expires int64
}

func (api *API) getRecordZoneID(name string) (string, error) {
	api.zoneCacheMutex.RLock()
	entry, ok := api.zoneCache[name]
	api.zoneCacheMutex.RUnlock()

	if ok && entry.expires <= time.Now().Unix() {
		return entry.zoneID, nil
	}
	soa, err := api.findSOA(name)
	if err != nil {
		return "", err
	}
	zoneID, err := api.api.ZoneIDByName(soa.Hdr.Name[:len(soa.Hdr.Name)-1])
	if err != nil {
		return "", fmt.Errorf("error retriving the zone id: %w", err)
	}
	entry = zoneCacheEntry{
		zoneID:  zoneID,
		expires: time.Now().Add(time.Duration(soa.Refresh) * time.Second).Unix(),
	}

	api.zoneCacheMutex.Lock()
	api.zoneCache[name] = entry
	api.zoneCacheMutex.Unlock()

	return zoneID, nil
}

func (api *API) findSOA(name string) (*dns.SOA, error) {
	var res *dns.Msg
	var err error
	c := dns.Client{}
	m := dns.Msg{}
	if name[len(name)-1] != '.' {
		name += "."
	}
	domainIndexes := dns.Split(name)
	for _, domain := range domainIndexes {
		m.SetQuestion(name[domain:], dns.TypeSOA)
		for _, server := range api.DNSServers {
			res, _, err = c.Exchange(&m, server)
			if err == nil {
				break
			}
		}
		if res == nil {
			continue
		}
		switch res.Rcode {
		case dns.RcodeSuccess:
			if len(res.Answer) == 0 {
				continue
			}
			for _, ans := range res.Answer {
				if soa, ok := ans.(*dns.SOA); ok {
					return soa, nil
				}
			}
		case dns.RcodeNameError:
			// NXDOMAIN
		default:
			return nil, ErrUnexpectedResponseType
		}
	}
	return nil, ErrNoDNSResults
}
