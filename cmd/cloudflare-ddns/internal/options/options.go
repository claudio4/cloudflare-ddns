package options

import (
	"errors"

	"github.com/caarlos0/env/v6"
	"github.com/jessevdk/go-flags"
)

// Options program options
type Options struct {
	APIKey       string `short:"k" long:"api-key" description:"Cloudflare account Global API Key (requires email to be set) (incompatible with token auth) [$CF_API_KEY]" env:"CF_API_KEY"`
	Email        string `short:"e" long:"email" description:"Cloudflare account email (only necessary for api key auth) [$CF_EMAIL]" env:"CF_EMAIL"`
	PrintVersion bool   `short:"v" long:"version" description:"Print program version and exit"`
	APIToken     string `short:"t" long:"token" description:"API token [$CF_TOKEN]" env:"CF_TOKEN"`
}

// Populate the Options struct with data from environment variables and arguments
func (opts *Options) Populate() error {
	err := env.Parse(opts)
	if err != nil {
		return err
	}
	_, err = flags.Parse(opts)
	if err != nil {
		return err
	}

	return opts.validate()
}

func (opts *Options) validate() error {
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

	return nil
}
