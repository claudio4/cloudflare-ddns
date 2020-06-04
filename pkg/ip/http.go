package ip

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

var dialer = &net.Dialer{
	Timeout:   30 * time.Second,
	KeepAlive: 30 * time.Second,
	DualStack: true,
}

func httpDoIPVersionRestricted(req *http.Request, ver version) (resp *http.Response, err error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
		var address string
		var ips []net.IP
		useIPv6 := ver == v6
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		ips4, ips6, err := lookupIPs(host)
		if err != nil {
			return nil, err
		}
		if useIPv6 {
			ips = ips6
		} else {
			ips = ips4
		}

		// retry with all IPs
		for _, ip := range ips {
			if useIPv6 {
				address = fmt.Sprintf("[%s]:%s", ip.String(), port)
			} else {
				address = fmt.Sprintf("%s:%s", ip.String(), port)
			}
			conn, err = dialer.DialContext(ctx, network, address)
			if err == nil {
				break
			}
		}
		return
	}
	client := http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}
	return client.Do(req)
}
