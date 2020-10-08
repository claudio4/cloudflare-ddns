package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/claudio4/cloudflare-ddns/cmd/cloudflare-ddns/internal/options"
	"github.com/claudio4/cloudflare-ddns/cmd/cloudflare-ddns/internal/unixsignals"
	"github.com/claudio4/cloudflare-ddns/pkg/cloudflare"
	"github.com/claudio4/cloudflare-ddns/pkg/ip"
	"github.com/rs/zerolog"
)

var version = "internal"
var opts options.Options
var api *cloudflare.API
var logger zerolog.Logger

func main() {
	var err error
	if err = opts.Populate(); err != nil {
		if err == options.ErrGratefulStop {
			os.Exit(0)
		}
		if err != options.ErrWithoutMessage {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		os.Exit(1)
	}
	if opts.PrintVersion {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}
	if opts.APIToken != "" {
		api, err = cloudflare.NewWithToken(opts.APIToken)
	} else {
		api, err = cloudflare.NewWithKey(opts.APIKey, opts.Email)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(5)
	}
	if len(opts.Resolvers) > 0 {
		api.DNSServers = opts.Resolvers
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if !opts.LogInJSONFormat {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "03:04", NoColor: opts.NoColors}).With().Timestamp().Logger()
	} else {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
	if opts.RefreshTime > 0 {
		daemonMode()
	} else {
		updateDomains()
	}
}

var (
	oldIPv4 net.IP
	oldIPv6 net.IP
)

func daemonMode() {
	sigc := signalProxy()
	ticker := time.Tick(opts.RefreshTime)
	logger.Info().Msgf("daemon mode activated, updating every %s", opts.RefreshTime)
	updateDomains()
	for {
		select {
		case <-sigc:
			return
		case <-ticker:
			updateDomains()
		}
	}
}

// This proxy listens to close signal and sends them trough sigproxyc
// if the same signal is received twice the program closes with exit code 55
func signalProxy() (sigproxyc chan os.Signal) {
	sigproxyc = make(chan os.Signal, 1)
	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, os.Interrupt)
		unixsignals.ListenUnixCloseSignals(sigc)

		sig := <-sigc
		logger.Info().
			Str("signal", sig.String()).
			Msg("stop signal received, if an update is ongoing the program will close once finished, resend the signal to force stop")
		sigproxyc <- sig

		sig = <-sigc
		logger.Info().
			Str("signal", sig.String()).
			Msg("second close signal received, forcing stop...")
		os.Exit(55)
	}()
	return sigproxyc
}

func updateDomains() {
	var wg sync.WaitGroup
	if !opts.OnlyIPv6 {
		wg.Add(1)
		go updateIPv4(&wg)
	}
	if !opts.OnlyIPv4 {
		wg.Add(1)
		go updateIPv6(&wg)
	}
	wg.Wait()
}

func updateIPv4(wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Info().Msg("fetching Ipv4")
	ip4, err := ip.GetV4()
	if err != nil {
		logger.Warn().
			Err(err).
			Msg("Could not get IPv4, skipping A records")
		return
	}
	logger.Info().
		Str("ip4", ip4.String()).
		Msg("")

	if ip4.Equal(oldIPv4) && !opts.ForceUpdate {
		logger.Info().Msg("IPv4 didn't change, skipping update")
		return
	}
	oldIPv4 = ip4

	var dwg sync.WaitGroup
	dwg.Add(len(opts.Domains))
	for _, domain := range opts.Domains {
		go updateARecord(&dwg, domain, ip4)
	}

	dwg.Wait()
}

func updateARecord(wg *sync.WaitGroup, domain string, ip net.IP) {
	defer wg.Done()
	logger.Info().
		Str("domain", domain).
		Str("recordType", "A").
		Msg("updating record")
	if err := api.UpdateARecord(domain, ip); err != nil {
		logger.Error().
			Err(err).
			Str("domain", domain).
			Msg("error updating domain")
	}
}

func updateIPv6(wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Info().Msg("fetching Ipv6")
	ip6, err := ip.GetV6()
	if err != nil {
		logger.Warn().
			Err(err).
			Msg("Could not get IPv6, skipping AAAA records")
		return
	}
	logger.Info().
		Str("ip6", ip6.String()).
		Msg("")
	if ip6.Equal(oldIPv6) && !opts.ForceUpdate {
		logger.Info().Msg("IPv6 didn't change, skipping update")
		return
	}
	oldIPv6 = ip6

	var dwg sync.WaitGroup
	dwg.Add(len(opts.Domains))
	for _, domain := range opts.Domains {
		go updateAAAARecord(&dwg, domain, ip6)
	}

	dwg.Wait()
}

func updateAAAARecord(wg *sync.WaitGroup, domain string, ip net.IP) {
	defer wg.Done()
	logger.Info().
		Str("domain", domain).
		Str("recordType", "AAAA").
		Msg("updating record")
	if err := api.UpdateAAAARecord(domain, ip); err != nil {
		logger.Error().
			Err(err).
			Str("domain", domain).
			Msg("error updating domain")
	}
}
