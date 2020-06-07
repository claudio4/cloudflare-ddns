package options

import (
	"errors"
	"time"

	"github.com/jessevdk/go-flags"
)

var (
	ErrGratefulStop   = errors.New("no error ocurred but the app should be stopped")
	ErrWithoutMessage = errors.New("an error ocurred but no message is required")
)

// Options program options
type Options struct {
	APIKey          string        `short:"k" long:"api-key" value-name:"<key>" description:"Cloudflare account Global API Key (requires email to be set) (incompatible with token auth)" env:"CF_API_KEY"`
	APIToken        string        `short:"t" long:"token" value-name:"<token>" description:"API token" env:"CF_TOKEN"`
	Domains         []string      `short:"d" long:"domain" value-name:"<domain.tld>" description:"Domains (or subdomains) to be updated" env:"CF_DOMAINS"`
	Email           string        `short:"e" long:"email" value-name:"<mail@example.cf>" description:"Cloudflare account email (only necessary for api key auth)" env:"CF_EMAIL"`
	ForceUpdate     bool          `short:"f" long:"force" description:"Force update even if the IP didn't change" env:"CF_FORCE_UPDATE"`
	RefreshTime     time.Duration `short:"r" long:"refresh-every" description:"Time between refreshing the IP on the domains (enables daemon mode)" env:"CF_REFRESH_EVERY"`
	Resolvers       []string      `long:"resolver" value-name:"<server>:<port>" description:"DNS resolvers to be used" env:"CF_RESOLVER" default-mask:"Cloudflare DNS"`
	LogInJSONFormat bool          `long:"json-log" description:"Format log as JSON" env:"CF_JSON_LOG"`
	OnlyIPv4        bool          `short:"4" long:"only-ipv4" description:"Only use IPv4 (A records)" env:"CF_ONLY_IPV4"`
	OnlyIPv6        bool          `short:"6" long:"only-ipv6" description:"Only use IPv6" env:"CF_ONLY_IPV6"`
	PrintVersion    bool          `short:"v" long:"version" description:"Print program version and exit"`
}

// Populate the Options struct with data from environment variables and arguments
func (opts *Options) Populate() error {
	_, err := flags.Parse(opts)
	if err != nil {
		if ferr, ok := err.(*flags.Error); ok {
			if ferr.Type == flags.ErrHelp {
				return ErrGratefulStop
			}
		}
		// the flag package automatically prints the error, so no more printing is required
		return ErrWithoutMessage
	}
	return opts.validate()
}

func (opts *Options) validate() error {
	if opts.PrintVersion {
		return nil
	}
	if opts.APIToken != "" {
		if opts.APIKey != "" || opts.Email != "" {
			return errors.New("token auth and api key auth can not be used at the same time")
		}
	} else {
		if opts.APIKey == "" && opts.Email == "" {
			return errors.New("no auth credentials provided")
		}
		if opts.APIKey == "" || opts.Email == "" {
			return errors.New("api key auth requires both, api key and email")
		}
	}
	if opts.OnlyIPv4 && opts.OnlyIPv6 {
		return errors.New("--only-ipv4 and --only-ipv6 can not be present at the same time")
	}
	if len(opts.Domains) == 0 {
		return errors.New("no domains specified")
	}

	return nil
}
