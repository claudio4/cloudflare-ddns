package cloudflare

import (
	"sync"

	cf "github.com/cloudflare/cloudflare-go"
)

type API struct {
	api            *cf.API
	DNSServers     []string
	zoneCache      map[string]zoneCacheEntry
	zoneCacheMutex sync.RWMutex
}

func newDefaultValuesAPI(cfapi *cf.API) API {
	return API{
		api:        cfapi,
		DNSServers: []string{"1.1.1.1:53", "1.0.0.1:53"},
		zoneCache:  map[string]zoneCacheEntry{},
	}
}

// NewWithKey creates a new Cloudflare API client using API Key
func NewWithKey(apiKey, email string) (*API, error) {
	cfapi, err := cf.New(apiKey, email)
	if err != nil {
		return nil, err
	}
	api := newDefaultValuesAPI(cfapi)

	return &api, nil
}

// NewWithToken creates a new Cloudflare API client using API Tokens
func NewWithToken(token string) (*API, error) {
	cfapi, err := cf.NewWithAPIToken(token)
	if err != nil {
		return nil, err
	}
	api := newDefaultValuesAPI(cfapi)

	return &api, nil
}
