package main

import (
	"flag"
	"fmt"
	"github.com/claudio4/cloudflare-ddns/pkg/cloudflare"
	"github.com/claudio4/cloudflare-ddns/pkg/ip"
	"os"
)

func main() {
	var email, apikey, zoneID, domain string
	flag.StringVar(&email, "email", "", "Your Cloudflare email")
	flag.StringVar(&apikey, "key", "", "Your Cloudflare API Key")
	flag.StringVar(&zoneID, "zoneid", "", "Domain's zone ID")
	flag.StringVar(&domain, "domain", "", "Domain to change")
	flag.Parse()
	if email == "" || apikey == "" || zoneID == "" || domain == "" {
		fmt.Println("Missing parameters: All the following parameters are required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	IP, err := ip.Get()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(10)
	}

	id, content, ttl, err := cloudflare.GetRecordDetails(email, apikey, zoneID, domain)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(20)
	}

	if content == IP {
		os.Exit(0)
	}

	err = cloudflare.SetRecord(email, apikey, zoneID, id, "A", domain, IP, ttl)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(30)
	}
}
