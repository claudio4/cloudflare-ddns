# 🌩️ Cloudflare DDNS
Dynamic DNS (DDNS) allows you to automatically update a domain's record with the current IP address of your device. This is particularly useful for devices with dynamic IP addresses, such as those connected to a home router.

This application provides an unofficial DDNS solution for Cloudflare users.

## 📋 Features
- IPv4 and IPv6 support
- Daemon mode for continuous updates
- Support for both API Token and API KEY authentication methods
- Pretty logger for easy debugging and troubleshooting
- JSON logger for structured logging

## 🚀 Getting started
1. Sign up for a Cloudflare account if you don't already have one.
[Generate an API Token](https://developers.cloudflare.com/fundamentals/api/get-started/create-token/) or [retrieve your API Key](https://developers.cloudflare.com/fundamentals/api/get-started/keys/) from the "My Settings" page in your Cloudflare account.
2. Download `cloudflare-ddns` from the [releases page](https://github.com/claudio4/cloudflare-ddns/releases).
3. Run the `cloudflare-ddns` command with the appropriate options. For example:
```sh
# Set the API Token
export CF_TOKEN=<api_token>
# Set the domain(s) to be updated
export CF_DOMAINS=example.com,test.com
# Execute the program
cloudflare-ddns

# or set the options via arguments
cloudflare-ddns -t <api_token> -d example.com -d test.com
```
This will update the DNS records for `example.com` and `test.com` with your current IP address.

### 🔄 Daemon mode
To continuously update your DNS records at a regular interval, you can use the `-r` / `--refresh-every` option, or set the `CF_REFRESH_EVERY` environment variable, to enable daemon mode. For example:
```sh
# Using environment variable
export CF_REFRESH_EVERY=1h
cloudflare-ddns

# Using arguments
cloudflare-ddns -t <api_token> -d example.com -r 1h
```
This will update the DNS records for `example.com` every hour. You can specify the interval with a unit of `s` for seconds, `m` for minutes, `h` for hours, or `d` for days.

## 🙏 Support
If you encounter any issues or have any questions, please open an [issue](https://github.com/claudio4/cloudflare-ddns/issues) on the repository. We will do our best to assist you.

### 🔎 Troubleshooting
If you encounter any issues while using this application, you can enable the `--json-log` option, or set the `CF_JSON_LOG` environment variable, to get structured logging output that may be helpful in debugging. You can also check the documentation or open an issue on the repository for additional help.

## 🎥 Demo
[![Video demo](https://asciinema.org/a/ek8WqVSKqGWuUaFUxNMX7dSv6.svg)](https://asciinema.org/a/ek8WqVSKqGWuUaFUxNMX7dSv6)
